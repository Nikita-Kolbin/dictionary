package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	op     = &slog.HandlerOptions{Level: slog.LevelInfo}
	jh     = slog.NewJSONHandler(os.Stdout, op)
	logger = slog.New(jh)
)

func Error(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}
