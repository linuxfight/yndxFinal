package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"orchestrator/cmd/app"
	"orchestrator/internal/db"
	"orchestrator/internal/db/tasksStorage"
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
	// connect to db
	cache, err := tasksStorage.NewCache()
	if err != nil {
		panic(err)
	}
	log.Info("connected to cache")

	userRepo, exprRepo, dbConn, err := db.NewSql()
	if err != nil {
		if dbConn != nil {
			if dbErr := dbConn.Close(context.Background()); dbErr != nil {
				panic(dbErr)
			}
		}
		panic(err)
	}
	log.Info("connected to database")

	web := app.NewFiber(userRepo, exprRepo, cache)
	stub := app.NewGrpc(cache)

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

func shutdown(app *fiber.App, dbConn *pgx.Conn, cache *tasksStorage.Cache) {
	err := app.Shutdown()
	if err != nil {
		panic(err)
	}

	db.CloseConnections(dbConn, cache.Storage)

	log.Info("shutdown complete")
	os.Exit(0)
}
