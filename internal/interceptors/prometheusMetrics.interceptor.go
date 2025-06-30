package interceptors

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_request",
			Help: "gRPC запросы",
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter)
}

func PrometheusInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		status := "OK"
		resp, err = handler(ctx, req)
		if err != nil {
			status = "ERROR"
		}
		RequestCounter.WithLabelValues(info.FullMethod, status).Inc()
		return resp, err
	}
}
