package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	if os.Getenv("ENVIRONMENT") == "prod" {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
