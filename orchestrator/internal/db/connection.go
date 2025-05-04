package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/valkey-io/valkey-go"
	"orchestrator/internal/db/expressions"
	"orchestrator/internal/db/users"
)

const defaultConn = "postgres://postgres:password@localhost:5432/db"

func NewSql() (*users.Queries, *expressions.Queries, *pgx.Conn, error) {
	// TODO: add config
	dbConn, err := pgx.Connect(context.Background(), defaultConn)
	if err != nil {
		return nil, nil, dbConn, err
	}

	usersRepo := users.New(dbConn)
	exprRepo := expressions.New(dbConn)

	if err := usersRepo.CreateSchema(context.Background()); err != nil {
		return nil, nil, dbConn, err
	}
	if err := exprRepo.CreateSchema(context.Background()); err != nil {
		return nil, nil, dbConn, err
	}

	return usersRepo, exprRepo, dbConn, nil
}

func CloseConnections(pgConn *pgx.Conn, cache valkey.Client) {
	cache.Close()

	if err := pgConn.Close(context.Background()); err != nil {
		panic(err)
	}
}
