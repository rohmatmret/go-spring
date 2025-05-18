package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache implements Cache interface using in-memory map
type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]item
}

type item struct {
	value      interface{}
	expiration int64
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]item),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	itm, ok := c.store[key]
	if !ok || (itm.expiration > 0 && itm.expiration < time.Now().UnixNano()) {
		return nil, false
	}
	return itm.value, true
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	c.store[key] = item{value: value, expiration: exp}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}

func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]item)
	return nil
}
