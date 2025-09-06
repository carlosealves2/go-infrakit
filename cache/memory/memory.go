package memory

import (
	"context"
	"sync"
	"time"

	"github.com/carlosealves2/go-infrakit/cache"
)

// Cache is an in-memory implementation of cache.Cache.
// It is safe for concurrent use.
type Cache struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// New creates a new in-memory cache.
func New() *Cache {
	return &Cache{store: make(map[string][]byte)}
}

func (c *Cache) checkCtx(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return cache.ErrTimeout
	}
	return nil
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	return c.SetBytes(ctx, key, []byte(value))
}

func (c *Cache) SetBytes(ctx context.Context, key string, value []byte) error {
	if err := c.checkCtx(ctx); err != nil {
		return err
	}
	c.mu.Lock()
	c.store[key] = append([]byte(nil), value...)
	c.mu.Unlock()
	return nil
}

func (c *Cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	if err := c.Set(ctx, key, value); err != nil {
		return err
	}
	if ttl > 0 {
		time.AfterFunc(ttl, func() {
			c.mu.Lock()
			delete(c.store, key)
			c.mu.Unlock()
		})
	}
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	b, err := c.GetBytes(ctx, key)
	return string(b), err
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	if err := c.checkCtx(ctx); err != nil {
		return nil, err
	}
	c.mu.RLock()
	v, ok := c.store[key]
	c.mu.RUnlock()
	if !ok {
		return nil, cache.ErrNotFound
	}
	return append([]byte(nil), v...), nil
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	if err := c.checkCtx(ctx); err != nil {
		return err
	}
	c.mu.Lock()
	for _, k := range keys {
		delete(c.store, k)
	}
	c.mu.Unlock()
	return nil
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	if err := c.checkCtx(ctx); err != nil {
		return false, err
	}
	c.mu.RLock()
	_, ok := c.store[key]
	c.mu.RUnlock()
	return ok, nil
}

var _ cache.Cache = (*Cache)(nil)
