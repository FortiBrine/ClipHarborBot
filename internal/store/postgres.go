package store

import (
	"context"
	"fmt"
	"time"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	poolMaxConns          = 10
	poolMinConns          = 2
	poolMaxConnLifetime   = time.Hour
	poolMaxConnIdleTime   = 15 * time.Minute
	poolHealthCheckPeriod = 30 * time.Second
)

func poolConfig(cfg config.Config) (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("parsing pool config: %w", err)
	}

	poolCfg.ConnConfig.Host = cfg.PostgresHost
	poolCfg.ConnConfig.Port = cfg.PostgresPort
	poolCfg.ConnConfig.User = cfg.PostgresUser
	poolCfg.ConnConfig.Password = cfg.PostgresPassword
	poolCfg.ConnConfig.Database = cfg.PostgresDb
	poolCfg.ConnConfig.TLSConfig = nil

	poolCfg.MaxConns = poolMaxConns
	poolCfg.MinConns = poolMinConns
	poolCfg.MaxConnLifetime = poolMaxConnLifetime
	poolCfg.MaxConnIdleTime = poolMaxConnIdleTime
	poolCfg.HealthCheckPeriod = poolHealthCheckPeriod

	return poolCfg, nil
}

func NewPostgres(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := poolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return pool, nil
}
