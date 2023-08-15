package controllers

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	logHandler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug, // default level for logHandler
			AddSource: true,            // log origin
		})

	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	return logger
}
