package main

import (
	"agent/internal/app"
	"os/signal"
	"syscall"
)

func main() {
	a := app.New()

	signal.Notify(a.ShutdownCh, syscall.SIGINT, syscall.SIGTERM)

	a.Start()

	<-a.ShutdownCh

	a.Stop()
}
