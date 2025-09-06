package cache

import (
	"context"
	"time"
)

// WithNamespace prefixes all keys with the provided namespace followed by ':'
// if namespace is empty, it returns the cache as is.
func WithNamespace(c Cache, ns string) Cache {
	if ns == "" {
		return c
	}
	return &namespaced{next: c, ns: ns}
}

type namespaced struct {
	next Cache
	ns   string
}

func (n *namespaced) prefix(key string) string {
	return n.ns + ":" + key
}

func (n *namespaced) Set(ctx context.Context, key, value string) error {
	return n.next.Set(ctx, n.prefix(key), value)
}

func (n *namespaced) Get(ctx context.Context, key string) (string, error) {
	return n.next.Get(ctx, n.prefix(key))
}

func (n *namespaced) Del(ctx context.Context, keys ...string) error {
	prefixed := make([]string, len(keys))
	for i, k := range keys {
		prefixed[i] = n.prefix(k)
	}
	return n.next.Del(ctx, prefixed...)
}

func (n *namespaced) Exists(ctx context.Context, key string) (bool, error) {
	return n.next.Exists(ctx, n.prefix(key))
}

func (n *namespaced) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	return n.next.SetWithTTL(ctx, n.prefix(key), value, ttl)
}

func (n *namespaced) SetBytes(ctx context.Context, key string, value []byte) error {
	return n.next.SetBytes(ctx, n.prefix(key), value)
}

func (n *namespaced) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return n.next.GetBytes(ctx, n.prefix(key))
}
