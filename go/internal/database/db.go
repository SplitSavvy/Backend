package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect establishes a connection pool and returns it.
func Connect(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	// Ping to ensure the connection is actually valid before returning
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
