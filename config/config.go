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
	JaegerPort     string
}

func InitConfig() *Config {
	grpcPort := flag.String("grpcPort", ":50051", "grpcPort to listen on")
	httpPort := flag.String("httpPort", ":2113", "httpPort for metrics")
	spotInstrumentAddr := flag.String("spotInstrument", "spot-instrument-service:50052", "Spot instrument address")
	redisPort := flag.String("redisPort", "redis:6379", "redisPort to redis client")
	JaegerPort := flag.String("jaegerPort", "jaeger:4318", "port for tracing")
	flag.Parse()

	cfg := &Config{
		GRPCPort:       *grpcPort,
		HTTPPort:       *httpPort,
		SpotInstrument: *spotInstrumentAddr,
		RedisPort:      *redisPort,
		JaegerPort:     *JaegerPort,
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

	if *spotInstrumentAddr == "spot-instrument-service:50052" {
		if envPort := os.Getenv("SPOT_INSTRUMENT_ADDR"); envPort != "" {
			cfg.SpotInstrument = envPort
		}
	}

	if *redisPort == "localhost:6379" {
		if envPort := os.Getenv("REDIS_PORT"); envPort != "" {
			cfg.RedisPort = envPort
		}
	}

	if *JaegerPort == "jaeger:4318" {
		if envPort := os.Getenv("JAEGER_AGENT_PORT"); envPort != "" {
			cfg.JaegerPort = envPort
		}
	}

	return cfg
}
