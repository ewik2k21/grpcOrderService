package tracing

import (
	"context"
	"github.com/ewik2k21/grpcOrderService/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitJaeger(ctx context.Context, serviceName string, cfg config.Config) (*trace.TracerProvider, error) {

	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(cfg.JaegerPort), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName)),
		))
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil

}
