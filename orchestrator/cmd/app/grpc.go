package app

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/controllers/tasksServer"
	"orchestrator/internal/controllers/utils"
	"orchestrator/internal/db/tasksStorage"
)

type TasksServer struct {
	tasksServer.UnimplementedOrchestratorServer

	server *grpc.Server
	cache  *tasksStorage.Cache
}

func (s *TasksServer) Listen() error {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	if err := s.server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *TasksServer) GetTask(_ *emptypb.Empty, stream grpc.ServerStreamingServer[tasksServer.TaskResponse]) error {
	// Get tasks from cache
	tasks, err := s.cache.GetTasks(context.Background())
	if err != nil {
		return err
	}

	// Stream tasks one by one
	for _, item := range *tasks {
		resp := &tasksServer.TaskResponse{}

		res, ok := item.Result.(string)
		if !ok || res != dto.Processing {
			continue
		}

		if _, ok := item.Arg1.(float64); !ok {
			continue
		}

		if _, ok := item.Arg2.(float64); !ok {
			continue
		}

		resp.Operator = utils.GetOperator(item)
		resp.Id = item.ID
		resp.Arg1 = fmt.Sprintf("%v", item.Arg1)
		resp.Arg2 = fmt.Sprintf("%v", item.Arg2)
		resp.Time = 1000 // TODO: Update with actual value

		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// Keep existing UpdateTask implementation
func (s *TasksServer) UpdateTask(ctx context.Context, req *tasksServer.UpdateTaskRequest) (*emptypb.Empty, error) {
	fmt.Printf("Received update for task %s: %s\n", req.Id, req.Result)
	return &emptypb.Empty{}, nil
}

func NewGrpc(storage *tasksStorage.Cache) *TasksServer {
	server := grpc.NewServer()
	ctl := TasksServer{
		cache:  storage,
		server: server,
	}
	tasksServer.RegisterOrchestratorServer(server, &ctl)
	reflection.Register(server)

	return &ctl
}
