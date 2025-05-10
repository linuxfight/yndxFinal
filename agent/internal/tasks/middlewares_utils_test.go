package tasks // mockClientStream implements grpc.ClientStream for testing
import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mockClientStream struct {
	grpc.ClientStream
}

func (m *mockClientStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (m *mockClientStream) Trailer() metadata.MD {
	return nil
}

func (m *mockClientStream) CloseSend() error {
	return nil
}

func (m *mockClientStream) Context() context.Context {
	return context.Background()
}

func (m *mockClientStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockClientStream) RecvMsg(msg interface{}) error {
	return nil
}
