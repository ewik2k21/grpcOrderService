package interceptors

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"time"
)

func LoggerRequestInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		requestID := ""
		if id, ok := md["x-request-id"]; ok && len(id) > 0 {
			requestID = id[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
			md = md.Copy()
			md.Set("x-request-id", requestID)
		}

		start := time.Now()
		logger.Info("Received grpc request",
			slog.String("method", info.FullMethod),
			slog.String("x-request-id", requestID),
			slog.Time("start_time", start),
		)

		resp, err = handler(ctx, req)

		duration := time.Since(start)
		logger.Info("Complete grpc request",
			slog.String("method", info.FullMethod),
			slog.String("x-request-id", requestID),
			slog.Any("duration", duration),
		)

		return resp, err
	}
}
