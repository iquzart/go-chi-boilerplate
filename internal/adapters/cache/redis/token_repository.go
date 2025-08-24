package redis

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
)

// Define a tracer for this package
var tracer = otel.Tracer("go-chi-boilerplate/internal/adapters/cache/redis")

type RefreshTokenRepository struct {
	client *RedisDB
	prefix string
	ttl    time.Duration
}

// NewRefreshTokenRepository creates a repository for refresh tokens
func NewRefreshTokenRepository(client *RedisDB, prefix string, ttl time.Duration) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}
}

// Add a helper method to generate the key
func (r *RefreshTokenRepository) createKey(userID string) string {
	return fmt.Sprintf("%s:refresh:%s", r.prefix, userID)
}

// SaveRefreshToken saves a refresh token for a user with the token's specific TTL.
func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, exp time.Time) error {
	ctx, span := tracer.Start(ctx, "RefreshTokenRepository.SaveRefreshToken")
	defer span.End()

	key := r.createKey(userID)
	duration := time.Until(exp)
	if duration <= 0 {
		return nil
	}
	return r.client.Client.Set(ctx, key, token, duration).Err()
}

// GetRefreshToken retrieves a refresh token for a user
func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	ctx, span := tracer.Start(ctx, "RefreshTokenRepository.GetRefreshToken")
	defer span.End()

	key := r.createKey(userID)
	return r.client.Client.Get(ctx, key).Result()
}

// DeleteRefreshToken removes a refresh token (logout)
func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	ctx, span := tracer.Start(ctx, "RefreshTokenRepository.DeleteRefreshToken")
	defer span.End()

	key := r.createKey(userID)
	return r.client.Client.Del(ctx, key).Err()
}
