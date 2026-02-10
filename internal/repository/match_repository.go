package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
)

type postgresMatchRepository struct {
	db *sql.DB
}

// NewPostgresMatchRepository creates a new PostgreSQL match repository
func NewPostgresMatchRepository(db *sql.DB) domain.MatchRepository {
	return &postgresMatchRepository{db: db}
}

func (r *postgresMatchRepository) Create(ctx context.Context, match *domain.Match) error {
	query := `
		INSERT INTO matches (id, team1_id, team2_id, status, rank_diff, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		match.ID,
		match.Team1ID,
		match.Team2ID,
		match.Status,
		match.RankDiff,
		match.CreatedAt,
		match.UpdatedAt,
		match.ExpiresAt,
	)
	return err
}

func (r *postgresMatchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Match, error) {
	query := `
		SELECT id, team1_id, team2_id, status, rank_diff, created_at, updated_at, expires_at
		FROM matches WHERE id = $1
	`
	match := &domain.Match{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&match.ID,
		&match.Team1ID,
		&match.Team2ID,
		&match.Status,
		&match.RankDiff,
		&match.CreatedAt,
		&match.UpdatedAt,
		&match.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("match not found")
	}
	return match, err
}

func (r *postgresMatchRepository) Update(ctx context.Context, match *domain.Match) error {
	query := `
		UPDATE matches 
		SET status = $2, updated_at = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, match.ID, match.Status, time.Now())
	return err
}

func (r *postgresMatchRepository) GetPendingMatches(ctx context.Context) ([]*domain.Match, error) {
	query := `
		SELECT id, team1_id, team2_id, status, rank_diff, created_at, updated_at, expires_at
		FROM matches 
		WHERE status = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, domain.MatchStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*domain.Match
	for rows.Next() {
		match := &domain.Match{}
		err := rows.Scan(
			&match.ID,
			&match.Team1ID,
			&match.Team2ID,
			&match.Status,
			&match.RankDiff,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}

func (r *postgresMatchRepository) GetExpiredMatches(ctx context.Context) ([]*domain.Match, error) {
	query := `
		SELECT id, team1_id, team2_id, status, rank_diff, created_at, updated_at, expires_at
		FROM matches 
		WHERE status = $1 AND expires_at < $2
		ORDER BY expires_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, domain.MatchStatusPending, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*domain.Match
	for rows.Next() {
		match := &domain.Match{}
		err := rows.Scan(
			&match.ID,
			&match.Team1ID,
			&match.Team2ID,
			&match.Status,
			&match.RankDiff,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}
