package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"orchestrator/internal/config"
	"orchestrator/internal/controllers"
	"orchestrator/internal/controllers/tasks"
	"orchestrator/internal/db"
	"os"
	"os/signal"
	"syscall"

	_ "orchestrator/docs"
)

// @title           Orchestrator API
// @version         1.0
// @description     Yandex Lyceum Calculator API
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"
func main() {
	cfg := config.New()

	// connect to db
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

	web := controllers.NewFiber(userRepo, exprRepo, cache)
	stub := tasks.NewGrpc(cache, exprRepo, cfg.OperationTime)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	listenCfg := fiber.ListenConfig{
		EnablePrefork:         false,
		DisableStartupMessage: true,
	}

	log.Info("started listening on port 8080")

	go func() {
		// start listening
		if err := web.Listen(":8080", listenCfg); err != nil {
			db.CloseConnections(dbConn, cache.Storage)
			panic(err)
		}
	}()

	log.Info("started listening on port 9090")

	go func() {
		err := stub.Listen()
		if err != nil {
			db.CloseConnections(dbConn, cache.Storage)
			panic(err)
		}
	}()

	<-sigChan
	shutdown(web, dbConn, cache)
}

func shutdown(app *fiber.App, dbConn *pgx.Conn, cache *db.Cache) {
	err := app.Shutdown()
	if err != nil {
		panic(err)
	}

	db.CloseConnections(dbConn, cache.Storage)

	log.Info("shutdown complete")
	os.Exit(0)
}
