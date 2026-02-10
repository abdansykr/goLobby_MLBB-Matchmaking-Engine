package domain

import (
	"time"

	"github.com/google/uuid"
)

// MatchStatus represents the status of a match
type MatchStatus string

const (
	MatchStatusPending   MatchStatus = "PENDING"
	MatchStatusConfirmed MatchStatus = "CONFIRMED"
	MatchStatusCancelled MatchStatus = "CANCELLED"
	MatchStatusCompleted MatchStatus = "COMPLETED"
)

// Match represents a matched pair of teams
type Match struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	Team1ID   uuid.UUID   `json:"team1_id" db:"team1_id"`
	Team2ID   uuid.UUID   `json:"team2_id" db:"team2_id"`
	Status    MatchStatus `json:"status" db:"status"`
	RankDiff  int         `json:"rank_diff" db:"rank_diff"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	ExpiresAt time.Time   `json:"expires_at" db:"expires_at"`
}

// NewMatch creates a new match between two teams
func NewMatch(team1ID, team2ID uuid.UUID, rankDiff int, readyTimeout int) *Match {
	now := time.Now()
	return &Match{
		ID:        uuid.New(),
		Team1ID:   team1ID,
		Team2ID:   team2ID,
		Status:    MatchStatusPending,
		RankDiff:  rankDiff,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(time.Duration(readyTimeout) * time.Second),
	}
}

// IsExpired checks if the match confirmation has expired
func (m *Match) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

// Confirm confirms the match
func (m *Match) Confirm() {
	m.Status = MatchStatusConfirmed
	m.UpdatedAt = time.Now()
}

// Cancel cancels the match
func (m *Match) Cancel() {
	m.Status = MatchStatusCancelled
	m.UpdatedAt = time.Now()
}

// Complete marks the match as completed
func (m *Match) Complete() {
	m.Status = MatchStatusCompleted
	m.UpdatedAt = time.Now()
}
