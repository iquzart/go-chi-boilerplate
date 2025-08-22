package postgresql

import (
	"go-chi-boilerplate/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB     *sql.DB
	Logger *slog.Logger
}

// New creates a new PostgreSQL connection
func New(cfg *config.DatabaseConfigs, logger *slog.Logger) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to open database connection",
			"host", cfg.Host, "db", cfg.DBName, "error", err,
		)
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Use a short timeout for ping to avoid blocking startup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("failed to ping database",
			"host", cfg.Host, "db", cfg.DBName, "error", err,
		)
		return nil, err
	}

	logger.Info("connected to PostgreSQL",
		"host", cfg.Host, "port", cfg.Port, "db", cfg.DBName,
	)

	return &PostgresDB{
		DB:     db,
		Logger: logger,
	}, nil
}

// Close closes the DB connection
func (p *PostgresDB) Close() {
	if err := p.DB.Close(); err != nil {
		p.Logger.Error("failed to close database connection", "error", err)
	} else {
		p.Logger.Info("database connection closed")
	}
}
