package config

import (
	"flag"
	"os"
)

type Config struct {
	GRPCPort       string
	HTTPPort       string
	SpotInstrument string
	RedisPort      string
}

func InitConfig() *Config {
	grpcPort := flag.String("grpcPort", ":50051", "grpcPort to listen on")
	httpPort := flag.String("httpPort", ":2113", "httpPort for metrics")
	spotInstrumentAddr := flag.String("spotInstrument", "localhost:50052", "Spot instrument address")
	redisPort := flag.String("redisPort", "localhost:6379", "redisPort to redis client")
	flag.Parse()

	cfg := &Config{
		GRPCPort:       *grpcPort,
		HTTPPort:       *httpPort,
		SpotInstrument: *spotInstrumentAddr,
		RedisPort:      *redisPort,
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

	if *spotInstrumentAddr == "localhost:50052" {
		if envPort := os.Getenv("SPOT_INSTRUMENT_ADDR"); envPort != "" {
			cfg.SpotInstrument = envPort
		}
	}

	if *redisPort == "localhost:6379" {
		if envPort := os.Getenv("REDIS_PORT"); envPort != "" {
			cfg.RedisPort = envPort
		}
	}

	return cfg
}
