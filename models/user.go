package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Return structs accept interfaces.
type User struct {
	ID           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

/*
type NewUser struct {
	Email    string
	Password string
}

func (us *UserService) Create(nu NewUser) (*User, error) {
	//
	return &User{}, nil
}
*/

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
		return nil, fmt.Errorf("create user: %v", err)
	}

	return &user, nil
}
