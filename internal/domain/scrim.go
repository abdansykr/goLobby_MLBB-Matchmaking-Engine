package domain

import (
	"time"

	"github.com/google/uuid"
)

// ScrimCategory represents the matchmaking category
type ScrimCategory string

const (
	CategoryPoke   ScrimCategory = "POKE"   // Rank weight 1-8, uses ±2 tolerance
	CategoryWarkop ScrimCategory = "WARKOP" // Rank weight 9-10, no tolerance
)

// ScrimStatus represents the current status of a scrim request
type ScrimStatus string

const (
	StatusSearching ScrimStatus = "searching"
	StatusMatched   ScrimStatus = "matched"
	StatusExpired   ScrimStatus = "expired"
	StatusCancelled ScrimStatus = "cancelled"
)

// ScrimRequest represents a team's scrim matchmaking request
type ScrimRequest struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	TeamName       string        `json:"team_name" db:"team_name"`
	WhatsAppNumber string        `json:"whatsapp_number" db:"whatsapp_number"`
	Category       ScrimCategory `json:"category" db:"category"`
	RankWeight     int           `json:"rank_weight" db:"rank_weight"` // 1-10
	Status         ScrimStatus   `json:"status" db:"status"`
	MatchID        *uuid.UUID    `json:"match_id,omitempty" db:"match_id"`
	IPAddress      string        `json:"ip_address,omitempty" db:"ip_address"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	ExpiresAt      time.Time     `json:"expires_at" db:"expires_at"`
	MatchedAt      *time.Time    `json:"matched_at,omitempty" db:"matched_at"`
}

// ScrimMatch represents a matched pair of teams
type ScrimMatch struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	Team1ID     uuid.UUID     `json:"team1_id" db:"team1_id"`
	Team2ID     uuid.UUID     `json:"team2_id" db:"team2_id"`
	Category    ScrimCategory `json:"category" db:"category"`
	RankDiff    *int          `json:"rank_diff,omitempty" db:"rank_diff"`
	Status      string        `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	ConfirmedAt *time.Time    `json:"confirmed_at,omitempty" db:"confirmed_at"`
	ExpiresAt   time.Time     `json:"expires_at" db:"expires_at"`

	// Populated relations
	Team1 *ScrimRequest `json:"team1,omitempty" db:"-"`
	Team2 *ScrimRequest `json:"team2,omitempty" db:"-"`
}

// Validation methods

// IsValidCategory checks if category is valid
func (sr *ScrimRequest) IsValidCategory() bool {
	return sr.Category == CategoryPoke || sr.Category == CategoryWarkop
}

// IsValidRankWeight checks if rank_weight is valid for the category
func (sr *ScrimRequest) IsValidRankWeight() bool {
	if sr.RankWeight < 1 || sr.RankWeight > 10 {
		return false
	}

	switch sr.Category {
	case CategoryPoke:
		return sr.RankWeight >= 1 && sr.RankWeight <= 8
	case CategoryWarkop:
		return sr.RankWeight >= 9 && sr.RankWeight <= 10
	default:
		return false
	}
}

// IsExpired checks if the request has expired
func (sr *ScrimRequest) IsExpired() bool {
	return time.Now().After(sr.ExpiresAt)
}

// CanMatchWith checks if this request can be matched with another
func (sr *ScrimRequest) CanMatchWith(other *ScrimRequest) bool {
	// Must be same category
	if sr.Category != other.Category {
		return false
	}

	// Both must be searching
	if sr.Status != StatusSearching || other.Status != StatusSearching {
		return false
	}

	// Neither should be expired
	if sr.IsExpired() || other.IsExpired() {
		return false
	}

	// POKE category: strict rank tolerance
	if sr.Category == CategoryPoke {
		// Rank 9 (Classic/Fun): exact match only
		if sr.RankWeight == 9 || other.RankWeight == 9 {
			return sr.RankWeight == other.RankWeight
		}
		// Ranks 1-8: ±1 tolerance
		rankDiff := abs(sr.RankWeight - other.RankWeight)
		return rankDiff <= 1
	}

	// WARKOP category: always can match (no rank restriction)
	return true
}

// GetWhatsAppURL generates WhatsApp URL for contacting opponent
func (sr *ScrimRequest) GetWhatsAppURL(opponentName string) string {
	message := "Hi " + opponentName + "! We've been matched for a scrim. Let's coordinate the match details!"
	return "https://wa.me/" + sr.WhatsAppNumber + "?text=" + urlEncode(message)
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Simple URL encoding for WhatsApp message
func urlEncode(s string) string {
	// Simplified version - in production use url.QueryEscape
	result := ""
	for _, c := range s {
		switch c {
		case ' ':
			result += "%20"
		case '!':
			result += "%21"
		case '\'':
			result += "%27"
		case ',':
			result += "%2C"
		default:
			result += string(c)
		}
	}
	return result
}

// MatchResponse represents the response when a match is found
type MatchResponse struct {
	Match          *ScrimMatch `json:"match"`
	OpponentName   string      `json:"opponent_name"`
	OpponentNumber string      `json:"opponent_number"`
	WhatsAppURL    string      `json:"whatsapp_url"`
	ExpiresIn      int         `json:"expires_in"` // seconds
}
