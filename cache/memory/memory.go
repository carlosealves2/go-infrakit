package memory

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/carlosealves2/go-infrakit/cache"
)

// Cache is an in-memory implementation of cache.Cache.
// It is safe for concurrent use.
type Cache struct {
	mu      sync.RWMutex
	store   map[string][]byte
	ns      string
	logger  cache.Logger
	tracer  trace.Tracer
	counter metric.Int64Counter
	latency metric.Float64Histogram
}

// New creates a new in-memory cache configured by opts.
func New(opts cache.Options) *Cache {
	c := &Cache{
		store:  make(map[string][]byte),
		ns:     opts.Namespace,
		logger: opts.Logger,
		tracer: opts.Tracer,
	}
	if opts.Meter != (metric.Meter{}) {
		c.counter, _ = opts.Meter.Int64Counter("cache_ops_total")
		c.latency, _ = opts.Meter.Float64Histogram("cache_latency_ms")
	}
	return c
}

func (c *Cache) formatKey(key string) (string, int) {
	keyLen := len(key)
	if c.ns != "" {
		key = c.ns + ":" + key
	}
	return key, keyLen
}

func (c *Cache) checkCtx(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return cache.ErrTimeout
	}
	return nil
}

func (c *Cache) observe(ctx context.Context, op string, keyLen int, hit bool, start time.Time, err error) {
	dur := time.Since(start)
	if c.logger != nil {
		fields := map[string]any{
			"mod":      "cache",
			"provider": "memory",
			"op":       op,
			"ns":       c.ns,
			"key_len":  keyLen,
			"dur_ms":   dur.Milliseconds(),
		}
		if err != nil {
			fields["err"] = err
		}
		c.logger.Log(fields)
	}
	if c.tracer != nil {
		_, span := c.tracer.Start(ctx, "cache."+op)
		span.SetAttributes(
			attribute.String("cache.provider", "memory"),
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
			attribute.String("provider", "memory"),
			attribute.String("op", op),
		}
		if op == "get" {
			attrs = append(attrs, attribute.Bool("hit", hit))
		}
		c.counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	if c.latency != (metric.Float64Histogram{}) {
		c.latency.Record(ctx, float64(dur.Milliseconds()), metric.WithAttributes(
			attribute.String("provider", "memory"),
			attribute.String("op", op),
		))
	}
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	return c.SetBytes(ctx, key, []byte(value))
}

func (c *Cache) SetBytes(ctx context.Context, key string, value []byte) error {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	if err := c.checkCtx(ctx); err != nil {
		c.observe(ctx, "set", keyLen, false, start, err)
		return err
	}
	c.mu.Lock()
	c.store[key] = append([]byte(nil), value...)
	c.mu.Unlock()
	c.observe(ctx, "set", keyLen, false, start, nil)
	return nil
}

func (c *Cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	if err := c.checkCtx(ctx); err != nil {
		c.observe(ctx, "set", keyLen, false, start, err)
		return err
	}
	c.mu.Lock()
	c.store[key] = []byte(value)
	c.mu.Unlock()
	if ttl > 0 {
		time.AfterFunc(ttl, func() {
			c.mu.Lock()
			delete(c.store, key)
			c.mu.Unlock()
		})
	}
	c.observe(ctx, "set", keyLen, false, start, nil)
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	b, err := c.GetBytes(ctx, key)
	return string(b), err
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	if err := c.checkCtx(ctx); err != nil {
		c.observe(ctx, "get", keyLen, false, start, err)
		return nil, err
	}
	c.mu.RLock()
	v, ok := c.store[key]
	c.mu.RUnlock()
	if !ok {
		err := cache.ErrNotFound
		c.observe(ctx, "get", keyLen, false, start, err)
		return nil, err
	}
	val := append([]byte(nil), v...)
	c.observe(ctx, "get", keyLen, true, start, nil)
	return val, nil
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
	if err := c.checkCtx(ctx); err != nil {
		c.observe(ctx, "del", keyLen, false, start, err)
		return err
	}
	c.mu.Lock()
	for _, k := range formatted {
		delete(c.store, k)
	}
	c.mu.Unlock()
	c.observe(ctx, "del", keyLen, false, start, nil)
	return nil
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	key, keyLen := c.formatKey(key)
	start := time.Now()
	if err := c.checkCtx(ctx); err != nil {
		c.observe(ctx, "exists", keyLen, false, start, err)
		return false, err
	}
	c.mu.RLock()
	_, ok := c.store[key]
	c.mu.RUnlock()
	c.observe(ctx, "exists", keyLen, false, start, nil)
	return ok, nil
}

var _ cache.Cache = (*Cache)(nil)
