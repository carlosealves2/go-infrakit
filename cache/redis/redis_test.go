package redis

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"

	"github.com/carlosealves2/go-infrakit/cache"
)

func newTestCache(t *testing.T) cache.Cache {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	c, err := New(cache.Options{Addr: mr.Addr()})
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
	mr, _ := miniredis.Run()
	defer mr.Close()
	c, err := New(cache.Options{Addr: mr.Addr()})
	if err != nil {
		t.Fatalf("new redis: %v", err)
	}
	if err := c.SetWithTTL(ctx, "foo", "bar", time.Second); err != nil {
		t.Fatalf("set ttl: %v", err)
	}
	mr.FastForward(time.Second + time.Millisecond)
	if _, err := c.Get(ctx, "foo"); err != cache.ErrNotFound {
		t.Fatalf("expected expiration, got %v", err)
	}
}

func TestRedisNamespace(t *testing.T) {
	ctx := context.Background()
	mr, _ := miniredis.Run()
	defer mr.Close()
	base, err := New(cache.Options{Addr: mr.Addr()})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	c := cache.WithNamespace(base, "ns")
	if err := c.Set(ctx, "foo", "bar"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if !mr.Exists("ns:foo") {
		t.Fatalf("expected namespaced key in redis")
	}
	if mr.Exists("foo") {
		t.Fatalf("unexpected plain key")
	}
}
