package meta

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// NewLogger creates and returns a new structured JSON logger with the specified log level.
// logLevel can be "debug", "info", "warn", "error". Default is "info".
func NewLogger(logLevel string) *slog.Logger {
	level := parseLogLevel(logLevel)
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

// parseLogLevel converts a string log level into slog.Level.
// Supports "debug", "info", "warn"/"warning", "error". Defaults to Info level.
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

// AuthEvent logs an authentication-related event for audit purposes.
// Parameters:
// - logger: the slog.Logger instance to use
// - event: event type (e.g., "login_success", "login_failed")
// - userID: the user ID if available
// - email: the email attempted or associated with the event
// - role: the role of the user
// - ip: client IP address
// - reason: reason for failure, if any
func AuthEvent(logger *slog.Logger, event, userID, email, role, ip, reason string) {
	logger.Info("auth_event",
		slog.String("event", event),
		slog.String("user_id", userID),
		slog.String("email", email),
		slog.String("role", role),
		slog.String("ip", ip),
		slog.String("reason", reason),
		slog.String("time", time.Now().Format(time.RFC3339)),
	)
}

// Fatal logs an error message and terminates the application immediately.
// Parameters:
// - logger: the slog.Logger instance to use
// - msg: the error message to log
// - args: optional key-value pairs for additional context
func Fatal(logger *slog.Logger, msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}
