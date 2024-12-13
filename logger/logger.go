package logger

import (
	"context"
	"github.com/bool64/ctxd"
	"github.com/bool64/zapctxd"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor wraps the request and response lifecycle to log the handler error.
func UnaryServerInterceptor(logger *zapctxd.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = ctxd.AddFields(ctx, "full_method", info.FullMethod)

		// Call the actual handler
		resp, err := handler(ctx, req)
		if err != nil {
			// Handle or log the error
			logger.Error(ctx, "handling call", "error", err)

			return nil, err
		}

		return resp, nil
	}
}
