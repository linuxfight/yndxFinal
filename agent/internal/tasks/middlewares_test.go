package tasks

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPanicUnaryClientInterceptor(t *testing.T) {
	tests := []struct {
		name        string
		invoker     grpc.UnaryInvoker
		expectError bool
		errorCode   codes.Code
	}{
		{
			name: "normal execution",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "panic with string",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				panic("test panic")
			},
			expectError: true,
			errorCode:   codes.Internal,
		},
		{
			name: "panic with error",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				panic(errors.New("test error"))
			},
			expectError: true,
			errorCode:   codes.Internal,
		},
		{
			name: "returned error without panic",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return status.Error(codes.NotFound, "not found")
			},
			expectError: true,
			errorCode:   codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := PanicUnaryClientInterceptor

			err := interceptor(
				context.Background(),
				"test/method",
				nil, nil,
				&grpc.ClientConn{},
				tt.invoker,
			)

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, tt.errorCode, status.Code(err))

				if tt.errorCode == codes.Internal {
					assert.Contains(t, err.Error(), "client panic:")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPanicStreamClientInterceptor(t *testing.T) {
	mockStream := &mockClientStream{}

	tests := []struct {
		name        string
		streamer    grpc.Streamer
		expectError bool
		errorCode   codes.Code
	}{
		{
			name: "normal execution",
			streamer: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return mockStream, nil
			},
			expectError: false,
		},
		{
			name: "panic with string",
			streamer: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				panic("test panic")
			},
			expectError: true,
			errorCode:   codes.Internal,
		},
		{
			name: "panic with error",
			streamer: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				panic(errors.New("test error"))
			},
			expectError: true,
			errorCode:   codes.Internal,
		},
		{
			name: "returned error without panic",
			streamer: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return nil, status.Error(codes.Unavailable, "service unavailable")
			},
			expectError: true,
			errorCode:   codes.Unavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := PanicStreamClientInterceptor

			stream, err := interceptor(
				context.Background(),
				&grpc.StreamDesc{},
				&grpc.ClientConn{},
				"test/method",
				tt.streamer,
			)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, stream)
				assert.Equal(t, tt.errorCode, status.Code(err))

				if tt.errorCode == codes.Internal {
					assert.Contains(t, err.Error(), "client panic:")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, mockStream, stream)
			}
		})
	}
}

func TestPanicWithCustomType(t *testing.T) {
	type CustomPanic struct {
		Message string
		Code    int
	}

	t.Run("unary with custom panic type", func(t *testing.T) {
		interceptor := PanicUnaryClientInterceptor

		err := interceptor(
			context.Background(),
			"test/method",
			nil, nil,
			&grpc.ClientConn{},
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				panic(CustomPanic{Message: "custom", Code: 42})
			},
		)

		require.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Contains(t, err.Error(), "client panic: {custom 42}")
	})

	t.Run("stream with custom panic type", func(t *testing.T) {
		interceptor := PanicStreamClientInterceptor

		_, err := interceptor(
			context.Background(),
			&grpc.StreamDesc{},
			&grpc.ClientConn{},
			"test/method",
			func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				panic(CustomPanic{Message: "custom", Code: 42})
			},
		)

		require.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Contains(t, err.Error(), "client panic: {custom 42}")
	})
}
