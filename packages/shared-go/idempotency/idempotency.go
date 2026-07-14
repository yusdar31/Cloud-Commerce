package idempotency

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store defines the interface for idempotency key storage.
type Store interface {
	// Exists checks if an idempotency key exists.
	Exists(ctx context.Context, key string) (bool, error)
	// Set stores an idempotency key with expiration.
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	// Get retrieves the value for an idempotency key.
	Get(ctx context.Context, key string) (string, error)
}

// RedisStore is a Redis implementation of idempotency.Store.
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new Redis-based idempotency store.
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

// Exists checks if an idempotency key exists in Redis.
func (s *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check idempotency key: %w", err)
	}
	return result > 0, nil
}

// Set stores an idempotency key in Redis with expiration.
func (s *RedisStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if err := s.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set idempotency key: %w", err)
	}
	return nil
}

// Get retrieves the value for an idempotency key from Redis.
func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	result, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get idempotency key: %w", err)
	}
	return result, nil
}

// GenerateKey creates a deterministic idempotency key from method, path, and body hash.
func GenerateKey(method, path string, body []byte) string {
	hash := sha256.Sum256(body)
	return fmt.Sprintf("idem:%s:%s:%x", method, path, hash[:8])
}

// GenerateKeyWithPrefix creates a prefixed idempotency key.
func GenerateKeyWithPrefix(prefix, method, path string, body []byte) string {
	hash := sha256.Sum256(body)
	return fmt.Sprintf("idem:%s:%s:%s:%x", prefix, method, path, hash[:8])
}
