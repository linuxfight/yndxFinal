package tasks

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PanicUnaryClientInterceptor catches panics in client calls and converts them to errors
func PanicUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {

	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			err = status.Errorf(codes.Internal, "client panic: %v", r)
		}
	}()

	return invoker(ctx, method, req, reply, cc, opts...)
}

// PanicStreamClientInterceptor handles panics in streaming RPCs
func PanicStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Internal, "client panic: %v", r)
		}
	}()

	return streamer(ctx, desc, cc, method, opts...)
}
