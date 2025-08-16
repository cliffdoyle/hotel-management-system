// File: backend/internal/database/db.go

package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect establishes a connection pool to the PostgreSQL database.
// It configures the pool for performance and longevity.
func Connect(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Set pool configuration
	config.MaxConns = 25                  // Max number of connections
	config.MinConns = 5                   // Min number of connections
	config.MaxConnLifetime = time.Hour    // Max lifetime of a connection
	config.MaxConnIdleTime = time.Minute * 30 // Idle time before closing a connection

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}