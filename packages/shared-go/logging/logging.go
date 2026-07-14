package logging

import (
	"log/slog"
	"os"
)

// New creates a configured slog logger for the given service.
func New(serviceName, environment string) *slog.Logger {
	level := slog.LevelInfo
	if environment == "development" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: environment == "development",
	})

	logger := slog.New(handler).With(
		"service", serviceName,
		"env", environment,
	)

	slog.SetDefault(logger)
	return logger
}
