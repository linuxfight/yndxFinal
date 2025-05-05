package main

import (
	"agent/internal/config"
	"agent/internal/tasks"
	"agent/internal/tasks/gen"
	"agent/internal/worker"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	c := config.New()
	log.Printf("Worker started with URL: %s", c.ApiAddr)
	log.Printf("Computing power: %d workers", c.ComputingPower)

	// Initialize gRPC connection
	conn, client := tasks.NewClient(c.ApiAddr)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	// Create buffered task channel and shutdown channel
	taskCh := make(chan *gen.TaskResponse, c.ComputingPower*2)
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	// Worker pool
	var wg sync.WaitGroup
	for i := 0; i < c.ComputingPower; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.Work(taskCh, client)
		}()
	}

	// Start task receiver
	go taskReceiver(client, taskCh, shutdownCh)

	// Wait for shutdown signal
	<-shutdownCh
	log.Println("Shutting down agent...")

	// Close task channel and wait for workers to finish
	close(taskCh)
	wg.Wait()
}

func taskReceiver(client gen.OrchestratorClient, taskCh chan<- *gen.TaskResponse, shutdownCh <-chan os.Signal) {

	for {
		select {
		case <-shutdownCh:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
			stream, err := client.GetTask(ctx, &emptypb.Empty{})
			time.Sleep(config.RequestTimeout)
			if err != nil {
				log.Printf("Failed to get task stream: %v (retrying)", err)
				cancel()
				time.Sleep(config.RequestTimeout)
				continue
			}

			receiveTasks(stream, taskCh, shutdownCh)

			cancel()
		}
	}
}

func receiveTasks(stream gen.Orchestrator_GetTaskClient, taskCh chan<- *gen.TaskResponse, shutdownCh <-chan os.Signal) {

	defer func(stream gen.Orchestrator_GetTaskClient) {
		err := stream.CloseSend()
		if err != nil {
			log.Printf("Error closing stream: %v", err)
		}
	}(stream)

	for {
		select {
		case <-shutdownCh:
			return
		default:
			task, err := stream.Recv()
			if err != nil {
				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					log.Printf("Stream error: %v (reconnecting)", err)
				}
				return
			}

			select {
			case taskCh <- task:
				log.Printf("Received new task: %s", task.Id)
			default:
				log.Println("Task channel full - dropping task")
			}
		}
	}
}
