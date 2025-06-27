package main

import (
	"context"
	"github.com/ewik2k21/grpcOrderService/cmd/server"
	"github.com/ewik2k21/grpcOrderService/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.InitConfig()

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	server.Execute(ctx, cfg, logger)
}
