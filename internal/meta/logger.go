package meta

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(logLevel string) *slog.Logger {
	level := parseLogLevel(logLevel)
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
