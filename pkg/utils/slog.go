package utils

import (
	"context"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func WithContext(ctx context.Context) *slog.Logger {
	return Logger.With("request_id", ctx.Value("request_id"))
}
