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
	queueKey       = "matchmaking:queue"
	lockKeyPrefix  = "matchmaking:lock:"
	matchKeyPrefix = "matchmaking:match:"
)

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
	// Get the team with the lowest score (oldest)
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
	// We need to find and remove the team by ID
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

// LockTeam locks a team for matching (anti-ghosting mechanism)
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

// SetMatchPending marks a match as pending in cache
func (r *redisCache) SetMatchPending(ctx context.Context, matchID uuid.UUID, ttl int) error {
	key := matchKeyPrefix + matchID.String()
	return r.client.Set(ctx, key, "PENDING", time.Duration(ttl)*time.Second).Err()
}

// GetMatchStatus gets the status of a match from cache
func (r *redisCache) GetMatchStatus(ctx context.Context, matchID uuid.UUID) (string, error) {
	key := matchKeyPrefix + matchID.String()
	return r.client.Get(ctx, key).Result()
}

// DeleteMatch removes a match from cache
func (r *redisCache) DeleteMatch(ctx context.Context, matchID uuid.UUID) error {
	key := matchKeyPrefix + matchID.String()
	return r.client.Del(ctx, key).Err()
}
