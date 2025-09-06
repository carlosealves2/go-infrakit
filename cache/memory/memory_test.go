package memory

import (
	"context"
	"testing"
	"time"

	"github.com/carlosealves2/go-infrakit/cache"
)

func TestMemoryCacheBasic(t *testing.T) {
	ctx := context.Background()
	m := New()
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
	m := New()
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
	base := New()
	ns := cache.WithNamespace(base, "ns")
	if err := ns.Set(ctx, "foo", "bar"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if _, err := base.Get(ctx, "foo"); err != cache.ErrNotFound {
		t.Fatalf("prefix not applied")
	}
	if val, err := base.Get(ctx, "ns:foo"); err != nil || val != "bar" {
		t.Fatalf("prefixed get failed: %v %s", err, val)
	}
}
