package redis

import (
	"context"
	"testing"
	"time"

	"github.com/carlosealves2/go-infrakit/cache"
)

func newTestCache(t *testing.T) cache.Cache {
	c, err := New(cache.Options{})
	if err != nil {
		t.Fatalf("new redis: %v", err)
	}
	return c
}

func TestRedisCacheBasic(t *testing.T) {
	ctx := context.Background()
	c := newTestCache(t)
	if err := c.Set(ctx, "foo", "bar"); err != nil {
		t.Fatalf("set: %v", err)
	}
	v, err := c.Get(ctx, "foo")
	if err != nil || v != "bar" {
		t.Fatalf("get: %v %s", err, v)
	}
	ok, err := c.Exists(ctx, "foo")
	if err != nil || !ok {
		t.Fatalf("exists: %v %v", err, ok)
	}
	if err := c.Del(ctx, "foo"); err != nil {
		t.Fatalf("del: %v", err)
	}
	if _, err := c.Get(ctx, "foo"); err != cache.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRedisCacheTTL(t *testing.T) {
	ctx := context.Background()
	c := newTestCache(t)
	if err := c.SetWithTTL(ctx, "foo", "bar", 50*time.Millisecond); err != nil {
		t.Fatalf("set ttl: %v", err)
	}
	time.Sleep(60 * time.Millisecond)
	if _, err := c.Get(ctx, "foo"); err != cache.ErrNotFound {
		t.Fatalf("expected expiration, got %v", err)
	}
}

func TestRedisNamespace(t *testing.T) {
	ctx := context.Background()
	c, err := New(cache.Options{Namespace: "ns"})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := c.Set(ctx, "foo", "bar"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if _, err := c.Get(ctx, "ns:foo"); err != cache.ErrNotFound {
		t.Fatalf("should not require manual prefix")
	}
	v, err := c.Get(ctx, "foo")
	if err != nil || v != "bar" {
		t.Fatalf("expected namespaced key, got %v %s", err, v)
	}
}
