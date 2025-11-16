package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

// DB holds the database connection
var DB *sql.DB

// Connect initializes the database connection
func Connect() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var err error
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		return err
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Ping DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := DB.PingContext(ctx); err != nil {
		return err
	}

	log.Println("Database connection established")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateTables creates the necessary database tables
func CreateTables() error {
	users := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`
	books := `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL
		);
	`

	if _, err := DB.Exec(users); err != nil {
		return err
	}
	if _, err := DB.Exec(books); err != nil {
		return err
	}

	log.Println("Database tables created/verified")
	return nil
}
