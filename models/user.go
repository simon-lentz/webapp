package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %v", err)
	}

	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id;`, email, passwordHash)

	if err = row.Scan(&user.ID); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := us.DB.QueryRow(`
	SELECT id, password_hash
	FROM users WHERE email=$1`, email)

	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		fmt.Printf("Invalid Password: %v\n", password)
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &user, nil
}

func (us *UserService) UpdatePassword(userID uint, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %v", err)
	}

	passwordHash := string(hashedBytes)

	if _, err = us.DB.Exec(`
	UPDATE users
	SET password_hash = $2
	WHERE id = $1;`, userID, passwordHash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
