package cache

import (
	"errors"
	"github.com/phuslu/log"
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

	// Observability
	Logger        *log.Logger
	EnableMetrics bool
	EnableTracing bool
}

var (
	ErrNotFound = errors.New("cache: not found")
	ErrTimeout  = errors.New("cache: timeout")
	ErrClosed   = errors.New("cache: closed")
)
