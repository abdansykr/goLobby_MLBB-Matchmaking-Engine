package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

// keyForIP generates Redis key for IP address
func (r *RedisRateLimiter) keyForIP(ipAddress string) string {
	return fmt.Sprintf("scrim:rate_limit:%s", ipAddress)
}

// CanRequest checks if IP can make a new request
func (r *RedisRateLimiter) CanRequest(ctx context.Context, ipAddress string) (bool, error) {
	key := r.keyForIP(ipAddress)

	// Check if key exists
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// If key doesn't exist, IP can make a request
	return exists == 0, nil
}

// SetActiveRequest sets active request for IP
func (r *RedisRateLimiter) SetActiveRequest(ctx context.Context, ipAddress string, requestID uuid.UUID) error {
	key := r.keyForIP(ipAddress)

	// Set with 30 minute expiration (same as scrim request expiry)
	err := r.client.Set(ctx, key, requestID.String(), 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set active request: %w", err)
	}

	return nil
}

// RemoveActiveRequest removes active request for IP
func (r *RedisRateLimiter) RemoveActiveRequest(ctx context.Context, ipAddress string) error {
	key := r.keyForIP(ipAddress)

	err := r.client.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to remove active request: %w", err)
	}

	return nil
}

// GetActiveRequestID gets the active request ID for IP
func (r *RedisRateLimiter) GetActiveRequestID(ctx context.Context, ipAddress string) (*uuid.UUID, error) {
	key := r.keyForIP(ipAddress)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // No active request
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active request ID: %w", err)
	}

	requestID, err := uuid.Parse(val)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request ID: %w", err)
	}

	return &requestID, nil
}
