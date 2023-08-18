package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/simon-lentz/webapp/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

// Maps to DB table.
type PasswordReset struct {
	ID     uint
	UserID uint
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
	// Verify existence of valid user email, retrieve user ID
	email = strings.ToLower(email)
	var userID uint
	row := prs.DB.QueryRow(`
	SELECT id FROM users
	WHERE email = $1;`, email)
	if err := row.Scan(&userID); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// Generate new token for PasswordReset.
	bytesPerToken := prs.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := prs.Duration
	if duration < 10*time.Minute {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: prs.hash(token),
		ExpiresAt: time.Now().Add(duration), // Potential refactor with time.Now func parameter for testing edge cases?
	}

	row = prs.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	if err = row.Scan(&pwReset.ID); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	// Verify password reset and user.
	tokenHash := prs.hash(token)
	var user User
	var pwReset PasswordReset
	row := prs.DB.QueryRow(`
	SELECT password_resets.id,
		password_resets.expires_at,
		users.id,
		users.email,
		users.password_hash
	FROM password_resets
		JOIN users ON users.id = password_resets.user_id
	WHERE password_resets.token_hash = $1;`, tokenHash)

	if err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}

	// Password Reset Token is valid, consume.
	if err := prs.delete(pwReset.ID); err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (prs *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (prs *PasswordResetService) delete(id uint) error {
	if _, err := prs.DB.Exec(`
	DELETE FROM password_resets
	WHERE id = $1;`, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
