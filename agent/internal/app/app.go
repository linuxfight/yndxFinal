package app

import (
	"agent/internal/config"
	"agent/internal/tasks"
	"agent/internal/tasks/gen"
	"agent/internal/worker"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
)

type App struct {
	cfg    *config.Config
	client gen.OrchestratorClient
	conn   *grpc.ClientConn

	taskCh     chan *gen.TaskResponse
	wg         *sync.WaitGroup
	ShutdownCh chan os.Signal
}

func (a *App) Start() {
	for i := 0; i < a.cfg.ComputingPower; i++ {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			worker.Work(a.taskCh, a.client)
		}()
	}
	go a.getTasks()
}

func (a *App) Stop() {
	err := a.conn.Close()
	if err != nil {
		log.Panicf("failed to close connection: %v", err)
	}

	close(a.taskCh)
	a.wg.Wait()
	log.Println("shutdown complete")
}

func New() *App {
	cfg := config.New()

	log.Printf("Worker started with URL: %s", cfg.ApiAddr)
	log.Printf("Computing power: %d workers", cfg.ComputingPower)

	conn, client := tasks.NewClient(cfg.ApiAddr)

	taskCh := make(chan *gen.TaskResponse, cfg.ComputingPower*2)
	shutdownCh := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}

	return &App{
		cfg:        cfg,
		client:     client,
		conn:       conn,
		taskCh:     taskCh,
		wg:         wg,
		ShutdownCh: shutdownCh,
	}
}
