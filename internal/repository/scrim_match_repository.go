package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
)

type ScrimMatchRepository struct {
	db *sql.DB
}

func NewScrimMatchRepository(db *sql.DB) *ScrimMatchRepository {
	return &ScrimMatchRepository{db: db}
}

// Create a new scrim match
func (r *ScrimMatchRepository) Create(ctx context.Context, match *domain.ScrimMatch) error {
	query := `
		INSERT INTO scrim_matches (
			id, team1_id, team2_id, category, rank_diff, 
			status, created_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, expires_at
	`

	match.ID = uuid.New()
	match.Status = "pending"
	match.CreatedAt = time.Now()
	match.ExpiresAt = time.Now().Add(60 * time.Second) // 60 seconds to confirm

	err := r.db.QueryRowContext(
		ctx,
		query,
		match.ID,
		match.Team1ID,
		match.Team2ID,
		match.Category,
		match.RankDiff,
		match.Status,
		match.CreatedAt,
		match.ExpiresAt,
	).Scan(&match.ID, &match.CreatedAt, &match.ExpiresAt)

	if err != nil {
		return fmt.Errorf("failed to create scrim match: %w", err)
	}

	return nil
}

// GetByID retrieves a scrim match by ID
func (r *ScrimMatchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ScrimMatch, error) {
	query := `
		SELECT id, team1_id, team2_id, category, rank_diff, 
			   status, created_at, confirmed_at, expires_at
		FROM scrim_matches
		WHERE id = $1
	`

	match := &domain.ScrimMatch{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&match.ID,
		&match.Team1ID,
		&match.Team2ID,
		&match.Category,
		&match.RankDiff,
		&match.Status,
		&match.CreatedAt,
		&match.ConfirmedAt,
		&match.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scrim match not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scrim match: %w", err)
	}

	return match, nil
}

// GetWithTeams retrieves match with team details populated
func (r *ScrimMatchRepository) GetWithTeams(ctx context.Context, id uuid.UUID) (*domain.ScrimMatch, error) {
	query := `
		SELECT 
			sm.id, sm.team1_id, sm.team2_id, sm.category, sm.rank_diff,
			sm.status, sm.created_at, sm.confirmed_at, sm.expires_at,
			t1.id, t1.team_name, t1.whatsapp_number, t1.category, t1.rank_weight,
			t1.status, t1.match_id, t1.ip_address, t1.created_at, t1.updated_at,
			t1.expires_at, t1.matched_at,
			t2.id, t2.team_name, t2.whatsapp_number, t2.category, t2.rank_weight,
			t2.status, t2.match_id, t2.ip_address, t2.created_at, t2.updated_at,
			t2.expires_at, t2.matched_at
		FROM scrim_matches sm
		JOIN scrim_requests t1 ON sm.team1_id = t1.id
		JOIN scrim_requests t2 ON sm.team2_id = t2.id
		WHERE sm.id = $1
	`

	match := &domain.ScrimMatch{}
	team1 := &domain.ScrimRequest{}
	team2 := &domain.ScrimRequest{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&match.ID,
		&match.Team1ID,
		&match.Team2ID,
		&match.Category,
		&match.RankDiff,
		&match.Status,
		&match.CreatedAt,
		&match.ConfirmedAt,
		&match.ExpiresAt,
		// Team 1
		&team1.ID,
		&team1.TeamName,
		&team1.WhatsAppNumber,
		&team1.Category,
		&team1.RankWeight,
		&team1.Status,
		&team1.MatchID,
		&team1.IPAddress,
		&team1.CreatedAt,
		&team1.UpdatedAt,
		&team1.ExpiresAt,
		&team1.MatchedAt,
		// Team 2
		&team2.ID,
		&team2.TeamName,
		&team2.WhatsAppNumber,
		&team2.Category,
		&team2.RankWeight,
		&team2.Status,
		&team2.MatchID,
		&team2.IPAddress,
		&team2.CreatedAt,
		&team2.UpdatedAt,
		&team2.ExpiresAt,
		&team2.MatchedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scrim match not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scrim match with teams: %w", err)
	}

	match.Team1 = team1
	match.Team2 = team2

	return match, nil
}

// UpdateStatus updates the match status
func (r *ScrimMatchRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE scrim_matches
		SET status = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update match status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("scrim match not found")
	}

	return nil
}

// Confirm confirms a match
func (r *ScrimMatchRepository) Confirm(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE scrim_matches
		SET status = $1, confirmed_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, "confirmed", time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to confirm match: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("scrim match not found")
	}

	return nil
}

// CancelExpiredMatches marks expired matches as cancelled
func (r *ScrimMatchRepository) CancelExpiredMatches(ctx context.Context) (int64, error) {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update matches
	query := `
		UPDATE scrim_matches
		SET status = 'cancelled'
		WHERE status = 'pending' AND expires_at < $1
	`

	result, err := tx.ExecContext(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to cancel expired matches: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Update scrim requests to remove match_id and set back to searching
	updateQuery := `
		UPDATE scrim_requests sr
		SET status = 'searching', match_id = NULL
		FROM scrim_matches sm
		WHERE (sr.id = sm.team1_id OR sr.id = sm.team2_id)
		  AND sm.status = 'cancelled'
		  AND sr.status = 'matched'
	`

	_, err = tx.ExecContext(ctx, updateQuery)
	if err != nil {
		return 0, fmt.Errorf("failed to update scrim requests: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return rows, nil
}

// GetPendingMatches retrieves all pending matches
func (r *ScrimMatchRepository) GetPendingMatches(ctx context.Context) ([]*domain.ScrimMatch, error) {
	query := `
		SELECT id, team1_id, team2_id, category, rank_diff,
			   status, created_at, confirmed_at, expires_at
		FROM scrim_matches
		WHERE status = 'pending'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending matches: %w", err)
	}
	defer rows.Close()

	var matches []*domain.ScrimMatch
	for rows.Next() {
		match := &domain.ScrimMatch{}
		err := rows.Scan(
			&match.ID,
			&match.Team1ID,
			&match.Team2ID,
			&match.Category,
			&match.RankDiff,
			&match.Status,
			&match.CreatedAt,
			&match.ConfirmedAt,
			&match.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scrim match: %w", err)
		}
		matches = append(matches, match)
	}

	return matches, nil
}
