package cache

import (
	"errors"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Driver represents the cache backend driver
type Driver string

const (
	MemoryDriver Driver = "memory"
	RedisDriver  Driver = "redis"
)

// Options defines configuration for cache instances.
// Fields are flat to keep usage simple and align with README examples.
type Options struct {
	Driver    Driver
	Namespace string

	// Redis specific fields
	Addr     string
	DB       int
	Username string
	Password string
	TLS      bool

	// Observability adapters
	Logger Logger
	Tracer trace.Tracer
	Meter  metric.Meter
}

var (
	ErrNotFound = errors.New("cache: not found")
	ErrTimeout  = errors.New("cache: timeout")
	ErrClosed   = errors.New("cache: closed")
)
