package domain

import (
	"time"

	"github.com/google/uuid"
)

// TeamStatus represents the current status of a team
type TeamStatus string

const (
	TeamStatusWaiting   TeamStatus = "WAITING"
	TeamStatusMatched   TeamStatus = "MATCHED"
	TeamStatusLocked    TeamStatus = "LOCKED"
	TeamStatusReady     TeamStatus = "READY"
	TeamStatusCancelled TeamStatus = "CANCELLED"
)

// Team represents a team entity in the matchmaking system
type Team struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	CaptainID      uuid.UUID  `json:"captain_id" db:"captain_id"`
	CaptainName    string     `json:"captain_name" db:"captain_name"`
	TeamName       string     `json:"team_name" db:"team_name"`
	AverageRank    int        `json:"average_rank" db:"average_rank"`
	Status         TeamStatus `json:"status" db:"status"`
	MatchID        *uuid.UUID `json:"match_id,omitempty" db:"match_id"`
	ReputationScore int       `json:"reputation_score" db:"reputation_score"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// NewTeam creates a new team with default values
func NewTeam(captainID uuid.UUID, captainName, teamName string, averageRank int) *Team {
	return &Team{
		ID:              uuid.New(),
		CaptainID:       captainID,
		CaptainName:     captainName,
		TeamName:        teamName,
		AverageRank:     averageRank,
		Status:          TeamStatusWaiting,
		ReputationScore: 100, // Default reputation score
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// IsWaiting checks if team is waiting for a match
func (t *Team) IsWaiting() bool {
	return t.Status == TeamStatusWaiting
}

// CanBeMatched checks if team can be matched
func (t *Team) CanBeMatched() bool {
	return t.Status == TeamStatusWaiting || t.Status == TeamStatusCancelled
}

// Lock locks the team for matching
func (t *Team) Lock(matchID uuid.UUID) {
	t.Status = TeamStatusLocked
	t.MatchID = &matchID
	t.UpdatedAt = time.Now()
}

// MarkAsReady marks the team as ready
func (t *Team) MarkAsReady() {
	t.Status = TeamStatusReady
	t.UpdatedAt = time.Now()
}

// Cancel cancels the team's matchmaking
func (t *Team) Cancel() {
	t.Status = TeamStatusCancelled
	t.MatchID = nil
	t.UpdatedAt = time.Now()
}

// ApplyGhostingPenalty applies penalty for not confirming ready
func (t *Team) ApplyGhostingPenalty(penalty int) {
	t.ReputationScore += penalty // penalty is negative
	if t.ReputationScore < 0 {
		t.ReputationScore = 0
	}
	t.UpdatedAt = time.Now()
}
