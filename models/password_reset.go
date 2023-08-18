package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

// Maps to DB table.
type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	// Given NOT NULL constraint we don't need sql.NullTime.
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// See SessionService in session.go.
	BytesPerToken int
	// Time for which PasswordReset is valid.
	Duration time.Duration
}

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Implement prs.Create")
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implements prs.Consume")
}
