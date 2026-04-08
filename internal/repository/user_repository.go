package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/golobby/matchmaking/internal/domain"
)

// UserRepository handles all database operations for the users table.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, whatsapp_number, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	user.ID = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.WhatsappNumber, user.AvatarURL,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByEmail retrieves a user by email. Returns sql.ErrNoRows if not found.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, whatsapp_number, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.WhatsappNumber, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err // caller checks sql.ErrNoRows
	}
	return user, nil
}

// FindByUsername retrieves a user by username. Returns sql.ErrNoRows if not found.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, whatsapp_number, avatar_url, created_at, updated_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.WhatsappNumber, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, whatsapp_number, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.WhatsappNumber, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateProfile updates user profile information
func (r *UserRepository) UpdateProfile(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, whatsapp_number = $3, avatar_url = $4, updated_at = $5
		WHERE id = $6
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.WhatsappNumber, user.AvatarURL, user.UpdatedAt, user.ID)
	return err
}

// UpdatePasswordHash updates only the password hash for a user
func (r *UserRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, hash, time.Now(), id)
	return err
}
