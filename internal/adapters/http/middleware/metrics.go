package middleware

import (
	"go-chi-boilerplate/internal/meta"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip /system/metrics and any other paths you want
		if r.URL.Path == "/system/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Wrap ResponseWriter to capture status
		ww := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		meta.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(ww.status)).Inc()
		meta.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
