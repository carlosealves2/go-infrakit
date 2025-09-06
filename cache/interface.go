package cache

import (
	"context"
	"time"
)

// Cache is the unified cache interface for memory and Redis providers.
// It is intentionally string-focused with optional byte helpers.
type Cache interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error
	SetBytes(ctx context.Context, key string, value []byte) error
	GetBytes(ctx context.Context, key string) ([]byte, error)
}
