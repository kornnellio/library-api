package database

import (
	"context"
	"database/sql"
	"time"

	"library-api/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"go.uber.org/zap"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// Connect initializes the database connection
func Connect(cfg *config.DatabaseConfig, logger *zap.Logger) (*DB, error) {
	if cfg.URL == "" {
		logger.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Ping DB with retry
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		if i == maxRetries-1 {
			return nil, err
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}

	logger.Info("Database connection established")
	return &DB{db}, nil
}

// Ping checks database connectivity
func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

// CreateTables creates the necessary database tables with indexes
func (db *DB) CreateTables(ctx context.Context) error {
	users := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
	books := `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
		CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
	`

	if _, err := db.ExecContext(ctx, users); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, books); err != nil {
		return err
	}

	return nil
}
