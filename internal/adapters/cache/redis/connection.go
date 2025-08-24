package redis

import (
	"context"
	"fmt"
	"go-chi-boilerplate/internal/config"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client *redis.Client
	Logger *slog.Logger
}

// New creates a new Redis connection
func New(cfg *config.RedisConfigs, logger *slog.Logger) (*RedisDB, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Use a short timeout for ping to avoid blocking startup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("failed to ping Redis",
			"host", cfg.Host, "port", cfg.Port, "db", cfg.DB, "error", err,
		)
		return nil, err
	}

	logger.Info("connected to Redis",
		"host", cfg.Host, "port", cfg.Port, "db", cfg.DB,
	)

	return &RedisDB{
		Client: rdb,
		Logger: logger,
	}, nil
}

// Close closes the Redis connection
func (r *RedisDB) Close() {
	if err := r.Client.Close(); err != nil {
		r.Logger.Error("failed to close Redis connection", "error", err)
	} else {
		r.Logger.Info("Redis connection closed")
	}
}
