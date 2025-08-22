package meta

import (
	"go-chi-boilerplate/internal/adapters/database/postgresql"
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

	// PostgreSQL database metrics
	PostgresDBOpenConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "postgresql_db_open_connections",
			Help: "Number of open connections to the PostgreSQL database",
		},
	)
	PostgresDBIdleConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "postgresql_db_idle_connections",
			Help: "Number of idle connections in the PostgreSQL pool",
		},
	)
	PostgresDBMaxOpenConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "postgresql_db_max_open_connections",
			Help: "Maximum allowed open connections to the PostgreSQL database",
		},
	)
)

// InitMetrics registers all metrics with Prometheus
func InitMetrics() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration)
}

// InitDBMetrics registers PostgreSQL metrics and updates them periodically
func InitDBMetrics(db *postgresql.PostgresDB) {
	prometheus.MustRegister(PostgresDBOpenConns, PostgresDBIdleConns, PostgresDBMaxOpenConns)

	updateMetrics := func() {
		stats := db.DB.Stats()
		PostgresDBOpenConns.Set(float64(stats.OpenConnections))
		PostgresDBIdleConns.Set(float64(stats.Idle))
		PostgresDBMaxOpenConns.Set(float64(stats.MaxOpenConnections))
	}

	// Initial update
	updateMetrics()

	// Update periodically
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateMetrics()
		}
	}()
}
