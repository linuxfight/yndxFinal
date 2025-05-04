package tasks

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/controllers/tasks/gen"
	"orchestrator/internal/db"
	"orchestrator/internal/db/expressions"
	"strconv"
)

type TasksServer struct {
	gen.UnimplementedOrchestratorServer

	server *grpc.Server
	cache  *db.Cache
	expr   *expressions.Queries
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

func (s *TasksServer) GetTask(_ *emptypb.Empty, stream grpc.ServerStreamingServer[gen.TaskResponse]) error {
	// Get tasks from cache
	tasks, err := s.cache.GetTasks(context.Background())
	if err != nil {
		return err
	}

	// Stream tasks one by one
	for _, item := range *tasks {
		resp := &gen.TaskResponse{}

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

		resp.Operator = getOperator(item)
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

func (s *TasksServer) UpdateTask(ctx context.Context, req *gen.UpdateTaskRequest) (*emptypb.Empty, error) {
	task, err := s.cache.GetTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	task.Result = req.Result

	if err := s.cache.SetTask(ctx, task); err != nil {
		return nil, err
	}

	if item, err := s.expr.GetById(ctx, req.Id); err == nil {
		val, err := strconv.ParseFloat(req.Result, 64)
		if err != nil {
			item.Error = true
		} else {
			item.Res = val
		}
		item.Finished = true
		if err := s.expr.Update(ctx, expressions.UpdateParams{
			Res:      item.Res,
			Finished: item.Finished,
			Error:    item.Error,
			ID:       item.ID,
		}); err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func NewGrpc(storage *db.Cache, expr *expressions.Queries) *TasksServer {
	server := grpc.NewServer()
	ctl := TasksServer{
		cache:  storage,
		server: server,
		expr:   expr,
	}
	gen.RegisterOrchestratorServer(server, &ctl)
	reflection.Register(server)

	return &ctl
}
