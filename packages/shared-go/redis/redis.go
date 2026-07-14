package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrNotFound is returned when a cache key is not found.
var ErrNotFound = errors.New("cache: key not found")

// Cache provides a simple Redis-based caching interface.
type Cache struct {
	client *redis.Client
}

// New creates a new Redis cache client from a connection URL.
func New(redisURL string) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis: failed to parse URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to connect: %w", err)
	}

	return &Cache{client: client}, nil
}

// Client returns the underlying Redis client for advanced operations.
func (c *Cache) Client() *redis.Client {
	return c.client
}

// Set stores a value in the cache with a TTL.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis: failed to marshal value: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from the cache and unmarshals it into the target.
func (c *Cache) Get(ctx context.Context, key string, target interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return fmt.Errorf("redis: failed to get key: %w", err)
	}
	return json.Unmarshal(data, target)
}

// Delete removes a key from the cache.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern (e.g., "products:*").
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis: scan error: %w", err)
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

// Exists checks if a key exists in the cache.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis: exists error: %w", err)
	}
	return n > 0, nil
}

// SetNX sets a value only if the key does not exist (useful for distributed locks).
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("redis: failed to marshal value: %w", err)
	}
	acquired, err := c.client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis: SetNX error: %w", err)
	}
	return acquired, nil
}

// Incr increments a counter key by 1.
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	n, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: incr error: %w", err)
	}
	return n, nil
}

// IncrWithTTL increments a counter and sets TTL if the key is new.
func (c *Cache) IncrWithTTL(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	n, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: incr error: %w", err)
	}
	if n == 1 {
		if err := c.client.Expire(ctx, key, ttl).Err(); err != nil {
			return 0, fmt.Errorf("redis: expire error: %w", err)
		}
	}
	return n, nil
}

// Close closes the Redis connection.
func (c *Cache) Close() error {
	return c.client.Close()
}

// Health checks if the Redis connection is healthy.
func (c *Cache) Health(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
