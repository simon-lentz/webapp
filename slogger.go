package main

import (
	"log/slog"
	"os"
)

func init() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Info msg goes here", "key", "value")
}
