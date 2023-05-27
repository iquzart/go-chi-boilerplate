package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "chi_http_request_duration_seconds",
			Help:    "Histogram of HTTP request duration.",
			Buckets: []float64{0.1, 0.3, 1, 3, 5, 10},
		},
		[]string{"service", "route", "method", "status"},
	)

	httpRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chi_http_request_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"service", "route", "method", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestCounter)
}

// PrometheusMetrics is a middleware that instruments the router with Prometheus metrics.
func PrometheusMetrics(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the handler with instrumentation code.
			recorder := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(recorder, r)

			duration := time.Since(start)
			route := chi.RouteContext(r.Context()).RoutePattern()
			method := r.Method
			status := strconv.Itoa(recorder.Status())

			// Record the duration and counter metrics.
			httpRequestDuration.WithLabelValues(serviceName, route, method, status).Observe(duration.Seconds())
			httpRequestCounter.WithLabelValues(serviceName, route, method, status).Inc()
		})
	}
}
