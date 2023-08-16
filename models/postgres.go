package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "temp",
		Password: "pmet",
		Database: "webapp",
		SSLMode:  "disable",
	}
}

// Open does not handle closing the SQL connection
// that it opens, so the db.Close() method should
// be used in conjunction with the Open(cfg) func.
func Open(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	return db, nil
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
