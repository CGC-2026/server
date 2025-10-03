// internal/data/store.go
package data

import (
	"context"
	"os"
	"time"

	sqlc "server/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewStore(ctx context.Context) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil { return nil, err }

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil { return nil, err }

	return &Store{Pool: pool, Queries: sqlc.New(pool)}, nil
}

func (s *Store) Close() { s.Pool.Close() }
