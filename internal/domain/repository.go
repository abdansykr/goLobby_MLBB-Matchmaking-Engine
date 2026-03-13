package domain

import (
	"context"

	"github.com/google/uuid"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*Team, error)
	GetByCaptainID(ctx context.Context, captainID uuid.UUID) (*Team, error)
	Update(ctx context.Context, team *Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWaitingTeams(ctx context.Context) ([]*Team, error)
	UpdateStatus(ctx context.Context, teamID uuid.UUID, status TeamStatus) error
	UpdateReputationScore(ctx context.Context, teamID uuid.UUID, score int) error
}

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	Create(ctx context.Context, match *Match) error
	GetByID(ctx context.Context, id uuid.UUID) (*Match, error)
	Update(ctx context.Context, match *Match) error
	GetPendingMatches(ctx context.Context) ([]*Match, error)
	GetExpiredMatches(ctx context.Context) ([]*Match, error)
}

// CacheRepository defines the interface for cache operations (Redis)
type CacheRepository interface {
	// Queue operations
	EnqueueTeam(ctx context.Context, team *Team) error
	DequeueTeam(ctx context.Context) (*Team, error)
	GetQueueLength(ctx context.Context) (int64, error)
	RemoveFromQueue(ctx context.Context, teamID uuid.UUID) error
	
	// Lock operations for anti-ghosting
	LockTeam(ctx context.Context, teamID uuid.UUID, matchID uuid.UUID, ttl int) error
	UnlockTeam(ctx context.Context, teamID uuid.UUID) error
	IsTeamLocked(ctx context.Context, teamID uuid.UUID) (bool, error)
	GetTeamLock(ctx context.Context, teamID uuid.UUID) (*uuid.UUID, error)
	
	// Match tracking
	SetMatchPending(ctx context.Context, matchID uuid.UUID, team1RequestID uuid.UUID, team2RequestID uuid.UUID, ttl int) error
	GetMatchStatus(ctx context.Context, matchID uuid.UUID) (string, error)
	GetMatchParticipants(ctx context.Context, matchID uuid.UUID) (team1RequestID, team2RequestID string, err error)
	// CancelMatchAtomically cancels a match only if it is still PENDING.
	// Returns (true, nil) if this caller was the one who performed the cancellation.
	// Returns (false, nil) if the match was already cancelled by someone else (idempotent).
	CancelMatchAtomically(ctx context.Context, matchID uuid.UUID) (cancelled bool, err error)
	DeleteMatch(ctx context.Context, matchID uuid.UUID) error

	// ── Double Opt-in Consensus ───────────────────────────────────────────
	// InitConsensus creates a Redis Hash match:{matchID}:consensus with both
	// participants set to "pending" and a TTL of ttlSeconds.
	InitConsensus(ctx context.Context, matchID uuid.UUID, requestID1 uuid.UUID, requestID2 uuid.UUID, ttlSeconds int) error
	// RecordAcceptance marks one participant as "accepted" in the hash.
	// Returns (allAccepted bool, err). allAccepted is true only when BOTH
	// participants have accepted — guaranteed atomic via Lua script.
	RecordAcceptance(ctx context.Context, matchID uuid.UUID, requestID uuid.UUID) (allAccepted bool, err error)
	// CancelConsensus deletes the consensus hash (used on rejection).
	CancelConsensus(ctx context.Context, matchID uuid.UUID) error
}
