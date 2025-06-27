package config

import (
	"flag"
	"os"
)

type Config struct {
	GRPCPort string
	HTTPPort string
}

func InitConfig() *Config {
	grpcPort := flag.String("grpcPort", ":50051", "grpcPort to listen on")
	httpPort := flag.String("httpPort", ":2113", "httpPort for metrics")

	flag.Parse()

	cfg := &Config{
		GRPCPort: *grpcPort,
		HTTPPort: *httpPort,
	}

	if *grpcPort == ":50051" {
		if envPort := os.Getenv("GRPC_PORT"); envPort != "" {
			cfg.GRPCPort = envPort
		}
	}

	if *httpPort == ":2113" {
		if envPort := os.Getenv("HTTP_PORT"); envPort != "" {
			cfg.HTTPPort = envPort
		}
	}
	return cfg
}
