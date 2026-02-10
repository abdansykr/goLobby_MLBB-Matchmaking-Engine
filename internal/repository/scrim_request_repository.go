package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
)

type ScrimRequestRepository struct {
	db *sql.DB
}

func NewScrimRequestRepository(db *sql.DB) *ScrimRequestRepository {
	return &ScrimRequestRepository{db: db}
}

// Create a new scrim request
func (r *ScrimRequestRepository) Create(ctx context.Context, request *domain.ScrimRequest) error {
	query := `
		INSERT INTO scrim_requests (
			id, team_name, whatsapp_number, category, rank_weight, 
			status, ip_address, created_at, updated_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at, expires_at
	`

	request.ID = uuid.New()
	request.Status = domain.StatusSearching
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.ExpiresAt = time.Now().Add(30 * time.Minute)

	err := r.db.QueryRowContext(
		ctx,
		query,
		request.ID,
		request.TeamName,
		request.WhatsAppNumber,
		request.Category,
		request.RankWeight,
		request.Status,
		request.IPAddress,
		request.CreatedAt,
		request.UpdatedAt,
		request.ExpiresAt,
	).Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt, &request.ExpiresAt)

	if err != nil {
		return fmt.Errorf("failed to create scrim request: %w", err)
	}

	return nil
}

// GetByID retrieves a scrim request by ID
func (r *ScrimRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ScrimRequest, error) {
	query := `
		SELECT id, team_name, whatsapp_number, category, rank_weight, 
			   status, match_id, ip_address, created_at, updated_at, 
			   expires_at, matched_at
		FROM scrim_requests
		WHERE id = $1
	`

	request := &domain.ScrimRequest{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.TeamName,
		&request.WhatsAppNumber,
		&request.Category,
		&request.RankWeight,
		&request.Status,
		&request.MatchID,
		&request.IPAddress,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.ExpiresAt,
		&request.MatchedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scrim request not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scrim request: %w", err)
	}

	return request, nil
}

// GetSearchingByCategory retrieves all searching requests for a category
func (r *ScrimRequestRepository) GetSearchingByCategory(ctx context.Context, category domain.ScrimCategory) ([]*domain.ScrimRequest, error) {
	query := `
		SELECT id, team_name, whatsapp_number, category, rank_weight, 
			   status, match_id, ip_address, created_at, updated_at, 
			   expires_at, matched_at
		FROM scrim_requests
		WHERE status = $1 AND category = $2
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusSearching, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get searching requests: %w", err)
	}
	defer rows.Close()

	var requests []*domain.ScrimRequest
	for rows.Next() {
		request := &domain.ScrimRequest{}
		err := rows.Scan(
			&request.ID,
			&request.TeamName,
			&request.WhatsAppNumber,
			&request.Category,
			&request.RankWeight,
			&request.Status,
			&request.MatchID,
			&request.IPAddress,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.ExpiresAt,
			&request.MatchedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scrim request: %w", err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// GetAllSearching retrieves ALL searching requests (any category)
func (r *ScrimRequestRepository) GetAllSearching(ctx context.Context) ([]*domain.ScrimRequest, error) {
	query := `
		SELECT id, team_name, whatsapp_number, category, rank_weight, 
			   status, match_id, ip_address, created_at, updated_at, 
			   expires_at, matched_at
		FROM scrim_requests
		WHERE status = $1 AND expires_at > NOW()
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusSearching)
	if err != nil {
		return nil, fmt.Errorf("failed to get all searching requests: %w", err)
	}
	defer rows.Close()

	var requests []*domain.ScrimRequest
	for rows.Next() {
		request := &domain.ScrimRequest{}
		err := rows.Scan(
			&request.ID,
			&request.TeamName,
			&request.WhatsAppNumber,
			&request.Category,
			&request.RankWeight,
			&request.Status,
			&request.MatchID,
			&request.IPAddress,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.ExpiresAt,
			&request.MatchedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scrim request: %w", err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// FindPotentialMatches finds potential matches for a request
func (r *ScrimRequestRepository) FindPotentialMatches(ctx context.Context, request *domain.ScrimRequest) ([]*domain.ScrimRequest, error) {
	var query string
	var args []interface{}

	if request.Category == domain.CategoryPoke {
		// Special case: Rank 9 (Classic/Fun) - exact match only
		if request.RankWeight == 9 {
			query = `
				SELECT id, team_name, whatsapp_number, category, rank_weight, 
					   status, match_id, ip_address, created_at, updated_at, 
					   expires_at, matched_at
				FROM scrim_requests
				WHERE status = $1 
				  AND category = $2
				  AND id != $3
				  AND rank_weight = $4
				  AND expires_at > NOW()
				ORDER BY created_at ASC
				LIMIT 10
			`
			args = []interface{}{
				domain.StatusSearching,
				domain.CategoryPoke,
				request.ID,
				9, // Exact match for Classic/Fun
			}
		} else {
			// POKE: rank tolerance ±1 (for ranks 1-8)
			query = `
				SELECT id, team_name, whatsapp_number, category, rank_weight, 
					   status, match_id, ip_address, created_at, updated_at, 
					   expires_at, matched_at
				FROM scrim_requests
				WHERE status = $1 
				  AND category = $2
				  AND id != $3
				  AND rank_weight BETWEEN $4 AND $5
				  AND expires_at > NOW()
				ORDER BY created_at ASC
				LIMIT 10
			`
			args = []interface{}{
				domain.StatusSearching,
				domain.CategoryPoke,
				request.ID,
				request.RankWeight - 1,
				request.RankWeight + 1,
			}
		}
	} else {
		// WARKOP: no rank tolerance, match anyone
		query = `
			SELECT id, team_name, whatsapp_number, category, rank_weight, 
				   status, match_id, ip_address, created_at, updated_at, 
				   expires_at, matched_at
			FROM scrim_requests
			WHERE status = $1 
			  AND category = $2
			  AND id != $3
			  AND expires_at > NOW()
			ORDER BY created_at ASC
			LIMIT 10
		`
		args = []interface{}{
			domain.StatusSearching,
			domain.CategoryWarkop,
			request.ID,
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find potential matches: %w", err)
	}
	defer rows.Close()

	var matches []*domain.ScrimRequest
	for rows.Next() {
		match := &domain.ScrimRequest{}
		err := rows.Scan(
			&match.ID,
			&match.TeamName,
			&match.WhatsAppNumber,
			&match.Category,
			&match.RankWeight,
			&match.Status,
			&match.MatchID,
			&match.IPAddress,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.ExpiresAt,
			&match.MatchedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan potential match: %w", err)
		}
		matches = append(matches, match)
	}

	return matches, nil
}

// UpdateStatus updates the status of a scrim request
func (r *ScrimRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ScrimStatus) error {
	query := `
		UPDATE scrim_requests
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("scrim request not found")
	}

	return nil
}

// UpdateMatchID updates the match_id and status
func (r *ScrimRequestRepository) UpdateMatchID(ctx context.Context, id uuid.UUID, matchID uuid.UUID) error {
	query := `
		UPDATE scrim_requests
		SET match_id = $1, status = $2, matched_at = $3, updated_at = $4
		WHERE id = $5
	`

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		matchID,
		domain.StatusMatched,
		now,
		now,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update match_id: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("scrim request not found")
	}

	return nil
}

// Cancel cancels a scrim request
func (r *ScrimRequestRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	return r.UpdateStatus(ctx, id, domain.StatusCancelled)
}

// ExpireOldRequests marks old searching requests as expired
func (r *ScrimRequestRepository) ExpireOldRequests(ctx context.Context) (int64, error) {
	query := `
		UPDATE scrim_requests
		SET status = $1, updated_at = $2
		WHERE status = $3 AND expires_at < $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		domain.StatusExpired,
		time.Now(),
		domain.StatusSearching,
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to expire old requests: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

// GetActiveByIP gets active request for an IP (for rate limiting)
func (r *ScrimRequestRepository) GetActiveByIP(ctx context.Context, ipAddress string) (*domain.ScrimRequest, error) {
	query := `
		SELECT id, team_name, whatsapp_number, category, rank_weight, 
			   status, match_id, ip_address, created_at, updated_at, 
			   expires_at, matched_at
		FROM scrim_requests
		WHERE ip_address = $1 
		  AND status IN ($2, $3)
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	request := &domain.ScrimRequest{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		ipAddress,
		domain.StatusSearching,
		domain.StatusMatched,
	).Scan(
		&request.ID,
		&request.TeamName,
		&request.WhatsAppNumber,
		&request.Category,
		&request.RankWeight,
		&request.Status,
		&request.MatchID,
		&request.IPAddress,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.ExpiresAt,
		&request.MatchedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No active request
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active request by IP: %w", err)
	}

	return request, nil
}

// Delete deletes a scrim request
func (r *ScrimRequestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM scrim_requests WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete scrim request: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("scrim request not found")
	}

	return nil
}
