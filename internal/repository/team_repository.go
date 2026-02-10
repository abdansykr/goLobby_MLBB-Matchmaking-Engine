package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type postgresTeamRepository struct {
	db *sql.DB
}

// NewPostgresTeamRepository creates a new PostgreSQL team repository
func NewPostgresTeamRepository(db *sql.DB) domain.TeamRepository {
	return &postgresTeamRepository{db: db}
}

func (r *postgresTeamRepository) Create(ctx context.Context, team *domain.Team) error {
	query := `
		INSERT INTO teams (id, captain_id, captain_name, team_name, average_rank, status, reputation_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		team.ID,
		team.CaptainID,
		team.CaptainName,
		team.TeamName,
		team.AverageRank,
		team.Status,
		team.ReputationScore,
		team.CreatedAt,
		team.UpdatedAt,
	)
	return err
}

func (r *postgresTeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	query := `
		SELECT id, captain_id, captain_name, team_name, average_rank, status, match_id, reputation_score, created_at, updated_at
		FROM teams WHERE id = $1
	`
	team := &domain.Team{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&team.ID,
		&team.CaptainID,
		&team.CaptainName,
		&team.TeamName,
		&team.AverageRank,
		&team.Status,
		&team.MatchID,
		&team.ReputationScore,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	return team, err
}

func (r *postgresTeamRepository) GetByCaptainID(ctx context.Context, captainID uuid.UUID) (*domain.Team, error) {
	query := `
		SELECT id, captain_id, captain_name, team_name, average_rank, status, match_id, reputation_score, created_at, updated_at
		FROM teams WHERE captain_id = $1 ORDER BY created_at DESC LIMIT 1
	`
	team := &domain.Team{}
	err := r.db.QueryRowContext(ctx, query, captainID).Scan(
		&team.ID,
		&team.CaptainID,
		&team.CaptainName,
		&team.TeamName,
		&team.AverageRank,
		&team.Status,
		&team.MatchID,
		&team.ReputationScore,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	return team, err
}

func (r *postgresTeamRepository) Update(ctx context.Context, team *domain.Team) error {
	query := `
		UPDATE teams 
		SET captain_name = $2, team_name = $3, average_rank = $4, status = $5, match_id = $6, reputation_score = $7
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		team.ID,
		team.CaptainName,
		team.TeamName,
		team.AverageRank,
		team.Status,
		team.MatchID,
		team.ReputationScore,
	)
	return err
}

func (r *postgresTeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresTeamRepository) GetWaitingTeams(ctx context.Context) ([]*domain.Team, error) {
	query := `
		SELECT id, captain_id, captain_name, team_name, average_rank, status, match_id, reputation_score, created_at, updated_at
		FROM teams 
		WHERE status = $1 
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, domain.TeamStatusWaiting)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		team := &domain.Team{}
		err := rows.Scan(
			&team.ID,
			&team.CaptainID,
			&team.CaptainName,
			&team.TeamName,
			&team.AverageRank,
			&team.Status,
			&team.MatchID,
			&team.ReputationScore,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

func (r *postgresTeamRepository) UpdateStatus(ctx context.Context, teamID uuid.UUID, status domain.TeamStatus) error {
	query := `UPDATE teams SET status = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, teamID, status)
	return err
}

func (r *postgresTeamRepository) UpdateReputationScore(ctx context.Context, teamID uuid.UUID, score int) error {
	query := `UPDATE teams SET reputation_score = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, teamID, score)
	return err
}
