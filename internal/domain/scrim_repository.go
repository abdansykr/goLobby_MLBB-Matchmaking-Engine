package domain

import (
	"context"

	"github.com/google/uuid"
)

// ScrimRequestRepository defines methods for scrim request data access
type ScrimRequestRepository interface {
	// Create a new scrim request
	Create(ctx context.Context, request *ScrimRequest) error

	// Get scrim request by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ScrimRequest, error)

	// Get all searching requests by category
	GetSearchingByCategory(ctx context.Context, category ScrimCategory) ([]*ScrimRequest, error)

	// Get all searching requests (any category)
	GetAllSearching(ctx context.Context) ([]*ScrimRequest, error)

	// Find potential matches for a request
	FindPotentialMatches(ctx context.Context, request *ScrimRequest) ([]*ScrimRequest, error)

	// Update scrim request status
	UpdateStatus(ctx context.Context, id uuid.UUID, status ScrimStatus) error

	// Update match ID
	UpdateMatchID(ctx context.Context, id uuid.UUID, matchID uuid.UUID) error

	// Cancel request
	Cancel(ctx context.Context, id uuid.UUID) error

	// Expire old requests (30+ minutes)
	ExpireOldRequests(ctx context.Context) (int64, error)

	// Get active request by IP (for rate limiting)
	GetActiveByIP(ctx context.Context, ipAddress string) (*ScrimRequest, error)

	// Delete request
	Delete(ctx context.Context, id uuid.UUID) error
}

// ScrimMatchRepository defines methods for scrim match data access
type ScrimMatchRepository interface {
	// Create a new match
	Create(ctx context.Context, match *ScrimMatch) error

	// Get match by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ScrimMatch, error)

	// Get match with team details
	GetWithTeams(ctx context.Context, id uuid.UUID) (*ScrimMatch, error)

	// Update match status
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// Confirm match
	Confirm(ctx context.Context, id uuid.UUID) error

	// Cancel expired matches
	CancelExpiredMatches(ctx context.Context) (int64, error)

	// Get pending matches
	GetPendingMatches(ctx context.Context) ([]*ScrimMatch, error)
}

// RateLimiter defines methods for rate limiting
type RateLimiter interface {
	// Check if IP can make a request
	CanRequest(ctx context.Context, ipAddress string) (bool, error)

	// Set active request for IP
	SetActiveRequest(ctx context.Context, ipAddress string, requestID uuid.UUID) error

	// Remove active request for IP
	RemoveActiveRequest(ctx context.Context, ipAddress string) error

	// Get active request ID for IP
	GetActiveRequestID(ctx context.Context, ipAddress string) (*uuid.UUID, error)
}
