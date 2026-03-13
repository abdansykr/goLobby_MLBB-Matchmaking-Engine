package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache repository
func NewRedisCache(client *redis.Client) domain.CacheRepository {
	return &redisCache{client: client}
}

const (
	queueKey           = "matchmaking:queue"
	lockKeyPrefix      = "matchmaking:lock:"
	matchKeyPrefix     = "matchmaking:match:"
	matchTeamsPostfix  = ":teams"
	consensusPostfix   = ":consensus" // Hash: {reqID1: "pending"|"accepted", reqID2: ...}
)

// matchKey returns the Redis key for a match's status
func matchKey(matchID uuid.UUID) string {
	return matchKeyPrefix + matchID.String()
}

// matchTeamsKey returns the Redis key for storing a match's participants
func matchTeamsKey(matchID uuid.UUID) string {
	return matchKeyPrefix + matchID.String() + matchTeamsPostfix
}

// ──────────────────────────────────────────────────────────────────────────────
// Queue operations
// ──────────────────────────────────────────────────────────────────────────────

// EnqueueTeam adds a team to the matchmaking queue
func (r *redisCache) EnqueueTeam(ctx context.Context, team *domain.Team) error {
	data, err := json.Marshal(team)
	if err != nil {
		return err
	}

	// Use sorted set with timestamp as score for FIFO ordering
	score := float64(time.Now().Unix())
	return r.client.ZAdd(ctx, queueKey, redis.Z{
		Score:  score,
		Member: data,
	}).Err()
}

// DequeueTeam removes and returns the oldest team from the queue
func (r *redisCache) DequeueTeam(ctx context.Context) (*domain.Team, error) {
	result, err := r.client.ZPopMin(ctx, queueKey, 1).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	var team domain.Team
	if err := json.Unmarshal([]byte(result[0].Member.(string)), &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// GetQueueLength returns the current queue size
func (r *redisCache) GetQueueLength(ctx context.Context) (int64, error) {
	return r.client.ZCard(ctx, queueKey).Result()
}

// RemoveFromQueue removes a specific team from the queue
func (r *redisCache) RemoveFromQueue(ctx context.Context, teamID uuid.UUID) error {
	members, err := r.client.ZRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, member := range members {
		var team domain.Team
		if err := json.Unmarshal([]byte(member), &team); err != nil {
			continue
		}

		if team.ID == teamID {
			return r.client.ZRem(ctx, queueKey, member).Err()
		}
	}

	return fmt.Errorf("team not found in queue")
}

// ──────────────────────────────────────────────────────────────────────────────
// Lock operations (anti-ghosting)
// ──────────────────────────────────────────────────────────────────────────────

// LockTeam locks a team for matching
func (r *redisCache) LockTeam(ctx context.Context, teamID uuid.UUID, matchID uuid.UUID, ttl int) error {
	key := lockKeyPrefix + teamID.String()
	return r.client.Set(ctx, key, matchID.String(), time.Duration(ttl)*time.Second).Err()
}

// UnlockTeam removes the lock on a team
func (r *redisCache) UnlockTeam(ctx context.Context, teamID uuid.UUID) error {
	key := lockKeyPrefix + teamID.String()
	return r.client.Del(ctx, key).Err()
}

// IsTeamLocked checks if a team is currently locked
func (r *redisCache) IsTeamLocked(ctx context.Context, teamID uuid.UUID) (bool, error) {
	key := lockKeyPrefix + teamID.String()
	exists, err := r.client.Exists(ctx, key).Result()
	return exists > 0, err
}

// GetTeamLock gets the match ID that locked the team
func (r *redisCache) GetTeamLock(ctx context.Context, teamID uuid.UUID) (*uuid.UUID, error) {
	key := lockKeyPrefix + teamID.String()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	matchID, err := uuid.Parse(val)
	if err != nil {
		return nil, err
	}

	return &matchID, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Match tracking
// ──────────────────────────────────────────────────────────────────────────────

// SetMatchPending marks a match as PENDING and stores its two participant request IDs.
// Two separate keys are written atomically via a pipeline:
//   matchmaking:match:<id>         → "PENDING"  (with TTL)
//   matchmaking:match:<id>:teams   → "team1ID|team2ID"  (with TTL)
func (r *redisCache) SetMatchPending(
	ctx context.Context,
	matchID uuid.UUID,
	team1RequestID uuid.UUID,
	team2RequestID uuid.UUID,
	ttl int,
) error {
	dur := time.Duration(ttl) * time.Second
	pipe := r.client.Pipeline()
	pipe.Set(ctx, matchKey(matchID), "PENDING", dur)
	pipe.Set(ctx, matchTeamsKey(matchID), team1RequestID.String()+"|"+team2RequestID.String(), dur)
	_, err := pipe.Exec(ctx)
	return err
}

// GetMatchStatus gets the status of a match from cache
func (r *redisCache) GetMatchStatus(ctx context.Context, matchID uuid.UUID) (string, error) {
	return r.client.Get(ctx, matchKey(matchID)).Result()
}

// GetMatchParticipants returns the two request IDs that belong to a match.
func (r *redisCache) GetMatchParticipants(ctx context.Context, matchID uuid.UUID) (string, string, error) {
	val, err := r.client.Get(ctx, matchTeamsKey(matchID)).Result()
	if err != nil {
		return "", "", fmt.Errorf("participants not found for match %s: %w", matchID, err)
	}

	var t1, t2 string
	if _, err := fmt.Sscanf(val, "%36s|%36s", &t1, &t2); err != nil {
		// Fallback: split by pipe manually
		parts := splitTwo(val)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("malformed participants value: %q", val)
		}
		t1, t2 = parts[0], parts[1]
	}

	return t1, t2, nil
}

// CancelMatchAtomically uses a Lua script to set match status to CANCELLED only if
// it is currently PENDING. This protects against the race condition where both
// players click "reject" at exactly the same moment — only one goroutine will
// execute the actual cancellation; the other will receive (false, nil).
//
// Lua guarantees atomicity because Redis executes scripts single-threadedly.
var cancelMatchScript = redis.NewScript(`
local key   = KEYS[1]
local cur   = redis.call("GET", key)
if cur == false then
  -- key expired/doesn't exist — treat as already cancelled
  return 0
end
if cur == "CANCELLED" then
  return 0
end
if cur ~= "PENDING" then
  return -1
end
redis.call("SET", key, "CANCELLED", "KEEPTTL")
return 1
`)

// CancelMatchAtomically returns (true, nil) if this call cancelled the match,
// (false, nil) if it was already cancelled by a concurrent caller.
func (r *redisCache) CancelMatchAtomically(ctx context.Context, matchID uuid.UUID) (bool, error) {
	result, err := cancelMatchScript.Run(ctx, r.client, []string{matchKey(matchID)}).Int()
	if err != nil {
		return false, fmt.Errorf("lua cancel script error: %w", err)
	}
	switch result {
	case 1:
		return true, nil  // this caller did the cancellation
	case 0:
		return false, nil // already cancelled, idempotent
	default:
		return false, fmt.Errorf("unexpected match status, cannot cancel")
	}
}

// DeleteMatch removes a match and its participants key from cache
func (r *redisCache) DeleteMatch(ctx context.Context, matchID uuid.UUID) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, matchKey(matchID))
	pipe.Del(ctx, matchTeamsKey(matchID))
	pipe.Del(ctx, consensusKey(matchID))
	_, err := pipe.Exec(ctx)
	return err
}

// ──────────────────────────────────────────────────────────────────────────────
// Double Opt-in Consensus
// ──────────────────────────────────────────────────────────────────────────────

// consensusKey returns the Redis key for a match's Double Opt-in Hash.
func consensusKey(matchID uuid.UUID) string {
	return matchKeyPrefix + matchID.String() + consensusPostfix
}

// InitConsensus creates a Redis Hash with both participants set to "pending".
// Uses HSETNX so that re-initialization is idempotent (never overwrites existing data).
func (r *redisCache) InitConsensus(
	ctx context.Context,
	matchID uuid.UUID,
	requestID1 uuid.UUID,
	requestID2 uuid.UUID,
	ttlSeconds int,
) error {
	key := consensusKey(matchID)
	dur := time.Duration(ttlSeconds) * time.Second
	pipe := r.client.Pipeline()
	pipe.HSetNX(ctx, key, requestID1.String(), "pending")
	pipe.HSetNX(ctx, key, requestID2.String(), "pending")
	pipe.Expire(ctx, key, dur)
	_, err := pipe.Exec(ctx)
	return err
}

// recordAcceptanceLua atomically:
//  1. Sets the caller's field to "accepted".
//  2. Checks whether ALL fields in the hash are now "accepted".
//  3. Returns 1 (all accepted) or 0 (still waiting on others).
//  4. Returns -1 if the key has already been cancelled/deleted.
var recordAcceptanceLua = redis.NewScript(`
local key      = KEYS[1]
local field    = ARGV[1]

-- If key has expired / been deleted, do nothing
if redis.call("EXISTS", key) == 0 then
  return -1
end

-- Mark this participant as accepted
redis.call("HSET", key, field, "accepted")

-- Check all values
local vals = redis.call("HVALS", key)
for _, v in ipairs(vals) do
  if v ~= "accepted" then
    return 0   -- at least one is still pending
  end
end
return 1       -- everyone has accepted!
`)

// RecordAcceptance atomically marks one participant as "accepted".
// Returns (true, nil) when BOTH participants have now accepted.
func (r *redisCache) RecordAcceptance(ctx context.Context, matchID uuid.UUID, requestID uuid.UUID) (bool, error) {
	result, err := recordAcceptanceLua.Run(
		ctx, r.client,
		[]string{consensusKey(matchID)},
		requestID.String(),
	).Int()
	if err != nil {
		return false, fmt.Errorf("RecordAcceptance lua error: %w", err)
	}
	switch result {
	case 1:
		return true, nil  // all accepted!
	case 0:
		return false, nil // still waiting
	case -1:
		return false, fmt.Errorf("consensus session expired or cancelled")
	default:
		return false, fmt.Errorf("unexpected lua result: %d", result)
	}
}

// CancelConsensus removes the consensus hash (call on rejection).
func (r *redisCache) CancelConsensus(ctx context.Context, matchID uuid.UUID) error {
	return r.client.Del(ctx, consensusKey(matchID)).Err()
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

func splitTwo(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == '|' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}
