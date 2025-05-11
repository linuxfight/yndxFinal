package app

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"orchestrator/internal/config"
	"orchestrator/internal/controllers"
	"orchestrator/internal/controllers/tasks"
	"orchestrator/internal/db"
)

var listenConfig = fiber.ListenConfig{DisableStartupMessage: true}

type App struct {
	web   *fiber.App
	stub  *tasks.Server
	db    *pgx.Conn
	cache *db.Cache
}

func (a App) Start() {
	log.Info("started listening on port 8080")

	go func() {
		if err := a.web.Listen(":8080", listenConfig); err != nil {
			db.CloseConnections(a.db, a.cache.Storage)
			panic(err)
		}
	}()

	log.Info("started listening on port 9090")

	go func() {
		err := a.stub.Listen()
		if err != nil {
			db.CloseConnections(a.db, a.cache.Storage)
			panic(err)
		}
	}()
}

func (a App) Stop() {
	err := a.web.Shutdown()
	if err != nil {
		panic(err)
	}
	a.stub.Close()
	db.CloseConnections(a.db, a.cache.Storage)

	log.Info("shutdown complete")
}

func New() *App {
	cfg := config.New()

	cache, err := db.NewCache(cfg.ValkeyConn)
	if err != nil {
		panic(err)
	}
	log.Info("connected to cache")

	userRepo, exprRepo, dbConn, err := db.NewSql(cfg.PostgresConn)
	if err != nil {
		if dbConn != nil {
			if dbErr := dbConn.Close(context.Background()); dbErr != nil {
				panic(dbErr)
			}
		}
		panic(err)
	}
	log.Info("connected to database")

	web := controllers.NewFiber(userRepo, exprRepo, cache, cfg.JwtSecret)
	stub := tasks.NewGrpc(cache, exprRepo, cfg)

	return &App{
		web:   web,
		stub:  stub,
		db:    dbConn,
		cache: cache,
	}
}
