package tasks

import (
	"agent/internal/tasks/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient(addr string) (*grpc.ClientConn, gen.OrchestratorClient) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}

	client := gen.NewOrchestratorClient(conn)

	return conn, client
}
