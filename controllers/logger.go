package controllers

import (
	"log/slog"
	"os"
)

// Add an os exec command to log stdout to jq and write a log-TIMESTAMP.json file.
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

// TODO: implement structured logging and error handling.
