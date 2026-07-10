package store

import (
	"context"
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func dsn(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresDb,
	)
}

func NewPostgres(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return pool, nil
}
