package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/simon-lentz/webapp/rand"
)

const (
	// Minimum bytes used to create each session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID     uint
	UserID uint
	// The token field is only set when creating a new session.
	// For a session lookup the field will be empty, only the
	// TokenHash field persists. This heads off a raw token leak.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to
	// use when generating each session token. If this value is not
	// set or set to less than the MinBytesPerToken const it will
	// be replaced with MinBytesPerToken.
	BytesPerToken int
}

func (ss *SessionService) Create(userID uint) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;`, session.UserID, session.TokenHash)
	if err = row.Scan(&session.ID); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id,
			users.email,
			users.password_hash 
		FROM sessions
			JOIN users on users.id = sessions.user_id 
		WHERE sessions.token_hash = $1;`, tokenHash)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	if _, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1`, tokenHash); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
