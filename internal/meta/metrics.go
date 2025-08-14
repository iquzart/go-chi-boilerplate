package meta

import (
	"go-chi-boilerplate/internal/adapters/secondary/database/postgresql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Database metrics
	DBOpenConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_open_connections",
			Help: "Number of open connections to the database",
		},
	)
	DBIdleConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_idle_connections",
			Help: "Number of idle connections in the database pool",
		},
	)
	DBMaxOpenConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_max_open_connections",
			Help: "Maximum allowed open connections to the database",
		},
	)
)

// InitMetrics registers all metrics with Prometheus
func InitMetrics() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration)
}

// InitDBMetrics registers database metrics and updates them periodically
func InitDBMetrics(db *postgresql.PostgresDB) {
	prometheus.MustRegister(DBOpenConns, DBIdleConns, DBMaxOpenConns)

	updateMetrics := func() {
		stats := db.DB.Stats()
		DBOpenConns.Set(float64(stats.OpenConnections))
		DBIdleConns.Set(float64(stats.Idle))
		DBMaxOpenConns.Set(float64(stats.MaxOpenConnections))
	}

	// Update metrics immediately
	updateMetrics()

	// Update periodically
	go func() {
		ticker := time.NewTicker(10 * time.Second) // adjust interval as needed
		defer ticker.Stop()
		for range ticker.C {
			updateMetrics()
		}
	}()
}
