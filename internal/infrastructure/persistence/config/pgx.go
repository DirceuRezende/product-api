package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL              string // postgres://user:pass@host:5432/db?sslmode=disable
	AppName          string
	MaxConns         int32
	MinConns         int32
	MaxConnLifetime  time.Duration
	MaxConnIdleTime  time.Duration
	ConnectTimeout   time.Duration
	StatementTimeout time.Duration
}

func DefaultConfig(url, appName string) Config {
	return Config{
		URL:              url,
		AppName:          appName,
		MaxConns:         20,
		MinConns:         2,
		MaxConnLifetime:  time.Hour,
		MaxConnIdleTime:  30 * time.Minute,
		ConnectTimeout:   5 * time.Second,
		StatementTimeout: 10 * time.Second,
	}
}

func NewPool(ctx context.Context, c Config) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(c.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	// Pool
	cfg.MaxConns = c.MaxConns
	cfg.MinConns = c.MinConns
	cfg.MaxConnLifetime = c.MaxConnLifetime
	cfg.MaxConnLifetimeJitter = 5 * time.Minute // evita que todas as conexões expirem juntas
	cfg.MaxConnIdleTime = c.MaxConnIdleTime
	cfg.HealthCheckPeriod = time.Minute

	// Conexão
	cfg.ConnConfig.ConnectTimeout = c.ConnectTimeout
	cfg.ConnConfig.RuntimeParams["application_name"] = c.AppName // aparece no pg_stat_activity
	cfg.ConnConfig.RuntimeParams["statement_timeout"] =
		fmt.Sprintf("%d", c.StatementTimeout.Milliseconds()) // proteção no lado do banco

	// Se usar PgBouncer em transaction mode, descomente:
	// cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// pgxpool.New é lazy e não conecta de verdade. O Ping garante fail-fast no startup.
	pingCtx, cancel := context.WithTimeout(ctx, c.ConnectTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
