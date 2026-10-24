package utils

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// Create a new logger that uses both stdout and OpenTelemetry handlers
func NewLogger(ctx context.Context) *slog.Logger {
	logProvider := NewLoggerProvider(ctx)
	otelHandler := otelslog.NewHandler("notion_saver", otelslog.WithLoggerProvider(logProvider))
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slog.NewMultiHandler(otelHandler, stdoutHandler))

	return logger
}
