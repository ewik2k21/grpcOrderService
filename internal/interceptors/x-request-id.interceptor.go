package interceptors

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDInterceptor() grpc.UnaryServerInterceptor {
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

		newCtx := metadata.NewIncomingContext(ctx, md)

		md, ok = metadata.FromIncomingContext(newCtx)
		if ok {
			if ids, exists := md["x-request-id"]; exists && len(ids) > 0 {
				grpc.SendHeader(ctx, metadata.Pairs("x-request-id", ids[0]))
			}
		}

		return handler(newCtx, req)
	}
}
