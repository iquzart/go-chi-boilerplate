package redis

import (
	"context"
	"time"
)

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

// SaveRefreshToken saves a refresh token for a user with the token's specific TTL.
// The ttl parameter should be derived from the JWT's "exp" claim.
func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, exp time.Time) error {
	key := r.prefix + ":" + "refresh:" + userID
	duration := time.Until(exp)
	if duration <= 0 {
		return nil // Token is already expired, no need to save
	}
	return r.client.Client.Set(ctx, key, token, duration).Err()
}

// GetRefreshToken retrieves a refresh token for a user
func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := r.prefix + ":" + "refresh:" + userID
	return r.client.Client.Get(ctx, key).Result()
}

// DeleteRefreshToken removes a refresh token (logout)
func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := r.prefix + ":" + "refresh:" + userID
	return r.client.Client.Del(ctx, key).Err()
}
