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
	SetMatchPending(ctx context.Context, matchID uuid.UUID, ttl int) error
	GetMatchStatus(ctx context.Context, matchID uuid.UUID) (string, error)
	DeleteMatch(ctx context.Context, matchID uuid.UUID) error
}
