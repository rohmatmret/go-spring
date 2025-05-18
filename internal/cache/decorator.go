package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-spring.com/internal/observability"
)

// Cacheable is a decorator function that adds caching to any function
func Cacheable[T any](
	cache Cache,
	keyGenerator func(args ...interface{}) string,
	ttl time.Duration,
) func(fn func(...interface{}) (T, error)) func(...interface{}) (T, error) {
	return func(fn func(...interface{}) (T, error)) func(...interface{}) (T, error) {
		return func(args ...interface{}) (T, error) {
			// Generate cache key
			key := keyGenerator(args...)

			// Try to get from cache
			if cached, found := cache.Get(context.Background(), key); found {
				if result, ok := cached.(T); ok {
					return result, nil
				}
			}

			// If not in cache, call the original function
			result, err := fn(args...)
			if err != nil {
				var zero T
				return zero, err
			}

			// Store in cache
			if err := cache.Set(context.Background(), key, result, ttl); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to cache result: %v\n", err)
			}

			return result, nil
		}
	}
}

// DefaultKeyGenerator generates a cache key from function arguments
func DefaultKeyGenerator(args ...interface{}) string {
	key := ""
	for _, arg := range args {
		// Convert argument to string representation
		if str, ok := arg.(string); ok {
			key += str
		} else {
			// For non-string types, use JSON marshaling
			if bytes, err := json.Marshal(arg); err == nil {
				key += string(bytes)
			}
		}
		key += ":"
	}
	return key
}

// CacheDecorator wraps methods with caching functionality
type CacheDecorator struct {
	cache   Cache
	metrics *observability.CacheMetrics
}

func NewCacheDecorator(cache Cache) *CacheDecorator {
	return &CacheDecorator{
		cache:   cache,
		metrics: observability.NewCacheMetrics(),
	}
}

// Cacheable decorates a function with caching
func (d *CacheDecorator) Cacheable(ctx context.Context, key string, ttl time.Duration, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	// Try to get from cache
	if cached, err := d.cache.Get(ctx, key); err {
		d.metrics.RecordHit()
		return cached, nil
	}
	d.metrics.RecordMiss()

	// Execute function
	result, err := fn(ctx)
	if err != nil {
		return nil, err
	}

	// Cache result
	if err := d.cache.Set(ctx, key, result, ttl); err != nil {
		d.metrics.RecordError()
	}

	return result, nil
}

// CacheEvict decorates a function with cache eviction
func (d *CacheDecorator) CacheEvict(ctx context.Context, key string, allEntries bool, fn func(context.Context) error) error {
	if allEntries {
		return d.cache.Clear(ctx)
	}
	return d.cache.Delete(ctx, key)
}

// CachePut decorates a function with cache update
func (d *CacheDecorator) CachePut(ctx context.Context, key string, ttl time.Duration, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	result, err := fn(ctx)
	if err != nil {
		return nil, err
	}

	if err := d.cache.Set(ctx, key, result, ttl); err != nil {
		d.metrics.RecordError()
	}

	return result, nil
}
