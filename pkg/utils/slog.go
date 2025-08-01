package utils

import (
	"context"
	"go_finance/internal/api/middleware"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func WithContext(ctx context.Context) *slog.Logger {
	return Logger.With("user_id", ctx.Value(middleware.UserIDKey))
}
