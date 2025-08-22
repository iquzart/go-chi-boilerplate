package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// responseWriter captures status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// ensure statusCode is set if WriteHeader wasn't called
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware returns a chi-compatible middleware that logs requests using the provided slog.Logger
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip /system/health if desired
			if r.URL.Path == "/system/health" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			ww := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			// capture caller file:line
			_, file, line, ok := runtime.Caller(2)
			caller := "unknown"
			if ok {
				caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
			}

			// Extract trace ID from context
			span := trace.SpanFromContext(r.Context())
			traceID := ""
			if span != nil && span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			logger.Info("http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.statusCode,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"caller", caller,
				"trace_id", traceID, // <-- added trace ID
			)
		})
	}
}
