package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "temp",
		Password: "pmet",
		Database: "webapp",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close() // Do not close connection to DB until application closes.

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Successful DB Connection!")

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		description TEXT
	);`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created.")

	name := "John Smith"
	email := "john@smith.com"
	_, err = db.Exec(`
	INSERT INTO users (name, email)
	VALUES ($1, $2);`, name, email)
	if err != nil {
		panic(err)
	}
	fmt.Println("User created.")

}
