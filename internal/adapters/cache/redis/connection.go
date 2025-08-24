package redis

import (
	"context"
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
	opts := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := redis.NewClient(opts)

	// Ping with timeout to validate connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("failed to ping Redis",
			"addr", cfg.Addr, "error", err,
		)
		return nil, err
	}

	logger.Info("connected to Redis",
		"addr", cfg.Addr, "db", cfg.DB,
	)

	return &RedisDB{
		Client: client,
		Logger: logger,
	}, nil
}

// Close closes the Redis client
func (r *RedisDB) Close() {
	if err := r.Client.Close(); err != nil {
		r.Logger.Error("failed to close Redis connection", "error", err)
	} else {
		r.Logger.Info("Redis connection closed")
	}
}
