package redis

import (
	"context"
	"crypto/tls"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/carlosealves2/go-infrakit/cache"
)

// Cache is a Redis-backed implementation of cache.Cache.
type Cache struct {
	client *goredis.Client
}

// New creates a new Redis cache.
func New(opts cache.Options) (*Cache, error) {
	rOpts := &goredis.Options{
		Addr:     opts.Addr,
		DB:       opts.DB,
		Username: opts.Username,
		Password: opts.Password,
	}
	if opts.TLS {
		rOpts.TLSConfig = &tls.Config{}
	}
	client := goredis.NewClient(rOpts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Cache{client: client}, nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return cache.ErrTimeout
	}
	if err == goredis.Nil {
		return cache.ErrNotFound
	}
	return err
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	return mapError(c.client.Set(ctx, key, value, 0).Err())
}

func (c *Cache) SetBytes(ctx context.Context, key string, value []byte) error {
	return mapError(c.client.Set(ctx, key, value, 0).Err())
}

func (c *Cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	return mapError(c.client.Set(ctx, key, value, ttl).Err())
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	return val, mapError(err)
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	return val, mapError(err)
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return mapError(c.client.Del(ctx, keys...).Err())
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n == 1, mapError(err)
}

var _ cache.Cache = (*Cache)(nil)
