package main

import (
	"orchestrator/internal/app"
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
// @description "Введите 'Bearer TOKEN' чтобы правильно использовать JWT API Token"
func main() {
	a := app.New()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a.Start()

	<-sigChan

	a.Stop()
}
