package cache

import (
	"context"
	"sync"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// CacheEntry represents a single cache entry with expiration
type CacheEntry struct {
	Value      interface{}
	Expiration time.Time
}

// InMemoryCache implements Cache interface using in-memory storage
type InMemoryCache struct {
	store sync.Map
}

// NewInMemoryCache creates a new in-memory cache instance
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	if value, ok := c.store.Load(key); ok {
		entry := value.(CacheEntry)
		if entry.Expiration.IsZero() || time.Now().Before(entry.Expiration) {
			return entry.Value, true
		}
		// Entry has expired, remove it
		c.store.Delete(key)
	}
	return nil, false
}

// Set stores a value in the cache with optional TTL
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	entry := CacheEntry{
		Value: value,
	}
	if ttl > 0 {
		entry.Expiration = time.Now().Add(ttl)
	}
	c.store.Store(key, entry)
	return nil
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.store.Delete(key)
	return nil
}

// Clear removes all values from the cache
func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.store = sync.Map{}
	return nil
}

// CacheEvict is equivalent to Spring's @CacheEvict
type CacheEvict struct {
	Key              string
	AllEntries       bool
	BeforeInvocation bool
}

// CachePut is equivalent to Spring's @CachePut
type CachePut struct {
	Key       string
	TTL       time.Duration
	Condition string
}
