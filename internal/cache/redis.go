package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go-spring.com/internal/observability"
)

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client  *redis.Client
	metrics *observability.CacheMetrics
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(host string, port int, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisCache{
		client:  client,
		metrics: observability.NewCacheMetrics(),
	}, nil
}

// Get retrieves a value from Redis
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, bool) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		c.metrics.RecordMiss()
		return nil, false
	}
	c.metrics.RecordHit()

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		c.metrics.RecordError()
		return nil, false
	}
	return result, true
}

// Set stores a value in Redis
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var val string

	// Try to marshal as JSON first
	if bytes, err := json.Marshal(value); err == nil {
		val = string(bytes)
	} else {
		// If not JSON-serializable, convert to string
		val = fmt.Sprintf("%v", value)
	}

	return c.client.Set(ctx, key, val, ttl).Err()
}

// Delete removes a value from Redis
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Clear removes all values from Redis
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}
