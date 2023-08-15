package main

import (
	"log/slog"
	"os"
)

func NewSlogger() *slog.Logger {
	logHandler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug, // default level for logHandler
			AddSource: true,            // log origin
		})

	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	return logger
}
