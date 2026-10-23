package utils

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// Create a new OpenTelemetry trace provider
func NewTracerProvider(ctx context.Context, env map[string]string) *sdktrace.TracerProvider {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		panic(err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(env["NOTION_SAVER_NAME"]),
			semconv.ServiceVersion(env["NOTION_SAVER_VERSION"]),
		),
	)
	if err != nil {
		panic(err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return provider
}

// Create a new OpenTelemetry log provider
func NewLoggerProvider(ctx context.Context) *sdklog.LoggerProvider {
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		panic(err)
	}

	processor := sdklog.NewBatchProcessor(exporter)
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(processor))

	return provider
}
