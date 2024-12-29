package logger

import (
	"context"
	"log/slog"
	"os"
)

var op = &slog.HandlerOptions{Level: slog.LevelInfo}
var jh = slog.NewJSONHandler(os.Stdout, op)
var logger = slog.New(jh)

func Error(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}
