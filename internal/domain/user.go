package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered player account.
type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"` // never serialise the hash
	WhatsappNumber string    `json:"whatsapp_number"`
	AvatarURL      string    `json:"avatar_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
