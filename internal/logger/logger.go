package logger

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

func SetupLogger(env string) *slog.Logger {
	var level slog.Level
	switch env {
	case "production":
		level = slog.LevelInfo
	case "development":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		devopts := &devslog.Options{
			HandlerOptions:    opts,
			MaxSlicePrintSize: 10,
			SortKeys:          true,
			NewLineAfterLog:   true,
			StringerFormatter: true,
		}
		handler = devslog.NewHandler(os.Stdout, devopts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
