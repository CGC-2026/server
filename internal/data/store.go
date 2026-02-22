package data

import (
	"context"
	"fmt"
	"os"
	"time"

	"server/internal/db"
	sqlc "server/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewStore(ctx context.Context) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Store{Pool: pool, Queries: sqlc.New(pool)}, nil
}

func (s *Store) Close() { s.Pool.Close() }

// WithTransaction handles the boilerplate of starting, committing, and rolling back a transaction.
func (s *Store) WithTransaction(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	// Ensure rollback if not committed
	defer tx.Rollback(ctx)

	// Create a new Queries instance with the transaction
	qtx := s.Queries.WithTx(tx)

	// Execute the function
	if err := fn(qtx); err != nil {
		return err
	}

	// If no errror, commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
