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
        base, err := New(cache.Options{})
        if err != nil {
                t.Fatalf("new: %v", err)
        }
        c := cache.WithNamespace(base, "ns")
        if err := c.Set(ctx, "foo", "bar"); err != nil {
                t.Fatalf("set: %v", err)
        }
        if _, err := base.Get(ctx, "foo"); err != cache.ErrNotFound {
                t.Fatalf("expected base miss, got %v", err)
        }
        v, err := base.Get(ctx, "ns:foo")
        if err != nil || v != "bar" {
                t.Fatalf("expected namespaced key, got %v %s", err, v)
        }
}
