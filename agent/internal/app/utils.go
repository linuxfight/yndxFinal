package app

import (
	"agent/internal/config"
	"agent/internal/tasks/gen"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"time"
)

func (a *App) getTasks() {
	for {
		select {
		case <-a.ShutdownCh:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
			stream, err := a.client.GetTask(ctx, &emptypb.Empty{})
			time.Sleep(config.RequestTimeout)
			if err != nil {
				log.Printf("Failed to get task stream: %v (retrying)", err)
				cancel()
				time.Sleep(config.RequestTimeout)
				continue
			}

			a.processStream(stream)

			cancel()
		}
	}
}

func (a *App) processStream(stream gen.Orchestrator_GetTaskClient) {

	defer func(stream gen.Orchestrator_GetTaskClient) {
		err := stream.CloseSend()
		if err != nil {
			log.Printf("Error closing stream: %v", err)
		}
	}(stream)

	for {
		select {
		case <-a.ShutdownCh:
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
			case a.taskCh <- task:
				log.Printf("Received new task: %s", task.Id)
			default:
				log.Println("Task channel full - dropping task")
			}
		}
	}
}
