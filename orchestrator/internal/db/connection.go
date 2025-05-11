package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
	"orchestrator/internal/db/expressions"
	"orchestrator/internal/db/users"
	"time"
)

func NewSql(conn string) (*users.Queries, *expressions.Queries, *pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, nil, nil, err
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, nil, nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	usersRepo := users.New(pool)
	exprRepo := expressions.New(pool)

	if err := usersRepo.CreateSchema(context.Background()); err != nil {
		return nil, nil, pool, err
	}
	if err := exprRepo.CreateSchema(context.Background()); err != nil {
		return nil, nil, pool, err
	}

	return usersRepo, exprRepo, pool, nil
}

func CloseConnections(pgConn *pgxpool.Pool, cache valkey.Client) {
	cache.Close()
	pgConn.Close()
}
