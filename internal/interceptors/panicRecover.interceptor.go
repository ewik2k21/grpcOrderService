package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"runtime/debug"
)

func UnaryPanicRecoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {

		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered",
					slog.Any("method", info.FullMethod),
					slog.Any("panic", r),
					slog.String("stacktrace", string(debug.Stack())),
				)

				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
