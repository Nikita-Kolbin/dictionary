package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

var (
	op     = &slog.HandlerOptions{Level: slog.LevelInfo}
	jh     = slog.NewJSONHandler(os.Stdout, op)
	logger = slog.New(jh)
)

func Error(ctx context.Context, msg string, args ...any) {
	args = append(args, funcName())
	logger.ErrorContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	args = append(args, funcName())
	logger.InfoContext(ctx, msg, args...)
}

func funcName() slog.Attr {
	pc, _, _, _ := runtime.Caller(2)
	name := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	return slog.String("func", name[len(name)-1])
}
