package redis

import (
	"context"
	"crypto/tls"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/carlosealves2/go-infrakit/cache"
	"github.com/carlosealves2/go-infrakit/observability/logger"
)

// Cache is a Redis-backed implementation of cache.Cache.
type Cache struct {
	client  *goredis.Client
	ns      string
	logger  logger.Logger
	tracer  trace.Tracer
	counter metric.Int64Counter
	latency metric.Float64Histogram
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
	c := &Cache{
		client: client,
		ns:     opts.Namespace,
		logger: opts.Logger,
		tracer: opts.Tracer,
	}
	if opts.Meter != (metric.Meter{}) {
		c.counter, _ = opts.Meter.Int64Counter("cache_ops_total")
		c.latency, _ = opts.Meter.Float64Histogram("cache_latency_ms")
	}
	return c, nil
}

func (c *Cache) formatKey(key string) (string, int) {
	keyLen := len(key)
	if c.ns != "" {
		key = c.ns + ":" + key
	}
	return key, keyLen
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

func (c *Cache) observe(ctx context.Context, op string, keyLen int, hit bool, start time.Time, err error) {
	dur := time.Since(start)
	if c.logger != nil {
		entry := c.logger.Info()
		if err != nil {
			entry = c.logger.Error().With("err", err)
		}
		entry.With("mod", "cache").With("provider", "redis").With("op", op).
			With("ns", c.ns).With("key_len", keyLen).With("dur_ms", dur.Milliseconds()).
			Log("msg", "")
	}
	if c.tracer != nil {
		_, span := c.tracer.Start(ctx, "cache."+op)
		span.SetAttributes(
			attribute.String("cache.provider", "redis"),
			attribute.String("cache.namespace", c.ns),
			attribute.Int("cache.key_len", keyLen),
		)
		if op == "get" {
			span.SetAttributes(attribute.Bool("cache.hit", hit))
		}
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}
	if c.counter != (metric.Int64Counter{}) {
		attrs := []attribute.KeyValue{
			attribute.String("provider", "redis"),
			attribute.String("op", op),
		}
		if op == "get" {
			attrs = append(attrs, attribute.Bool("hit", hit))
		}
		c.counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	if c.latency != (metric.Float64Histogram{}) {
		c.latency.Record(ctx, float64(dur.Milliseconds()), metric.WithAttributes(
			attribute.String("provider", "redis"),
			attribute.String("op", op),
		))
	}
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	err := mapError(c.client.Set(ctx, key, value, 0).Err())
	c.observe(ctx, "set", keyLen, false, start, err)
	return err
}

func (c *Cache) SetBytes(ctx context.Context, key string, value []byte) error {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	err := mapError(c.client.Set(ctx, key, value, 0).Err())
	c.observe(ctx, "set", keyLen, false, start, err)
	return err
}

func (c *Cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	err := mapError(c.client.Set(ctx, key, value, ttl).Err())
	c.observe(ctx, "set", keyLen, false, start, err)
	return err
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	val, err := c.client.Get(ctx, key).Result()
	err = mapError(err)
	c.observe(ctx, "get", keyLen, err == nil, start, err)
	return val, err
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	val, err := c.client.Get(ctx, key).Bytes()
	err = mapError(err)
	c.observe(ctx, "get", keyLen, err == nil, start, err)
	return val, err
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	formatted := make([]string, len(keys))
	key, keyLen := c.formatKey(keys[0])
	formatted[0] = key
	for i := 1; i < len(keys); i++ {
		formatted[i], _ = c.formatKey(keys[i])
	}
	start := time.Now()
	err := mapError(c.client.Del(ctx, formatted...).Err())
	c.observe(ctx, "del", keyLen, false, start, err)
	return err
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	n, err := c.client.Exists(ctx, key).Result()
	err = mapError(err)
	c.observe(ctx, "exists", keyLen, false, start, err)
	return n == 1, err
}

var _ cache.Cache = (*Cache)(nil)
