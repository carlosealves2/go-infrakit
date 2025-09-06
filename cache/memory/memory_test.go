package memory

import (
	"context"
	"testing"
	"time"

	"github.com/carlosealves2/go-infrakit/cache"
)

func TestMemoryCacheBasic(t *testing.T) {
	ctx := context.Background()
	m := New(cache.Options{})
	if err := m.Set(ctx, "foo", "bar"); err != nil {
		t.Fatalf("set: %v", err)
	}
	val, err := m.Get(ctx, "foo")
	if err != nil || val != "bar" {
		t.Fatalf("get: %v %s", err, val)
	}
	ok, err := m.Exists(ctx, "foo")
	if err != nil || !ok {
		t.Fatalf("exists: %v %v", err, ok)
	}
	if err := m.Del(ctx, "foo"); err != nil {
		t.Fatalf("del: %v", err)
	}
	if _, err := m.Get(ctx, "foo"); err != cache.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryCacheTTL(t *testing.T) {
	ctx := context.Background()
	m := New(cache.Options{})
	if err := m.SetWithTTL(ctx, "foo", "bar", 50*time.Millisecond); err != nil {
		t.Fatalf("set ttl: %v", err)
	}
	time.Sleep(60 * time.Millisecond)
	if _, err := m.Get(ctx, "foo"); err != cache.ErrNotFound {
		t.Fatalf("expected expiration, got %v", err)
	}
}

func TestMemoryNamespace(t *testing.T) {
	ctx := context.Background()
	c := New(cache.Options{Namespace: "ns"})
	if err := c.Set(ctx, "foo", "bar"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if _, err := c.Get(ctx, "ns:foo"); err != cache.ErrNotFound {
		t.Fatalf("should not require manual prefix")
	}
	if val, err := c.Get(ctx, "foo"); err != nil || val != "bar" {
		t.Fatalf("get failed: %v %s", err, val)
	}
}
