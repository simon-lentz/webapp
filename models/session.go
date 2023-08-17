package models

import "database/sql"

type Session struct {
	ID     int
	UserID int
	// The token field is only set when creating a new session.
	// For a session lookup the field will be empty, only the
	// TokenHash field persists. This heads off a raw token leak.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// create session token
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
