package server

import (
	"context"
	"github.com/ewik2k21/grpcOrderService/config"
	"github.com/ewik2k21/grpcOrderService/internal/handlers"
	"github.com/ewik2k21/grpcOrderService/internal/interceptors"
	"github.com/ewik2k21/grpcOrderService/internal/repositories"
	"github.com/ewik2k21/grpcOrderService/internal/services"
	order_service_v1 "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	spot_instrument_service_v1 "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Execute(ctx context.Context, cfg *config.Config, logger *slog.Logger) {
	wg := sync.WaitGroup{}

	//conn for spot instrument client
	conn, err := grpc.Dial(cfg.SpotInstrument, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to spot instrument service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close()

	spotInstrumentClient := spot_instrument_service_v1.NewSpotInstrumentServiceClient(conn)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor(),
			interceptors.LoggerRequestInterceptor(logger),
			interceptors.PrometheusInterceptor(),
			interceptors.UnaryPanicRecoveryInterceptor(logger)))

	//redis client create
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisPort,
		Password: "",
		DB:       0,
	})
	cacheTTL := 1 * time.Minute

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("failed to connect to Redis: %v", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Redis connect on ", slog.String("port", cfg.RedisPort))

	orderRepo := repositories.NewOrderRepository(logger)
	orderService := services.NewOrderService(*orderRepo, spotInstrumentClient, logger, redisClient, cacheTTL)
	orderHandler := handlers.NewOrderHandler(logger, orderService)

	order_service_v1.RegisterOrderServiceServer(grpcServer, orderHandler)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logger.Error("failed listen tcp server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server listening at %s", slog.Any("port", lis.Addr()))

	//start grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("start tcp server")
		if err = grpcServer.Serve(lis); err != nil {
			logger.Error("failed start grpc server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	//cfg for httpserver for metrics
	metricsServer := &http.Server{
		Addr: cfg.HTTPPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				promhttp.Handler().ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	//start server for metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("metrics endpoint start on :", slog.String("port", cfg.HTTPPort))
		http.Handle("/metrics", promhttp.Handler())
		if err = metricsServer.ListenAndServe(); err != nil {
			logger.Error("metrics endpoint failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("received shutdown signal, start graceful shutdown")
	//shutdown grpc
	grpcServer.GracefulStop()
	logger.Info("grpc server stopped")

	//shutdown redis
	if err = redisClient.Close(); err != nil {
		logger.Error("redis client shutdown failed", slog.String("error", err.Error()))
	}

	//shutdown metrics server
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("metrics server shutdown failed", slog.String("error", err.Error()))
	}

	wg.Wait()
	logger.Info("all stopped")
}
