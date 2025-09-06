package infrakit

import (
	"fmt"

	"github.com/carlosealves2/go-infrakit/cache"
	"github.com/carlosealves2/go-infrakit/cache/memory"
	redisdrv "github.com/carlosealves2/go-infrakit/cache/redis"
)

// NewCache initializes a cache according to the provided options.
// It validates the driver, configures the provider, applies namespace and observability.
func NewCache(opts cache.Options) (cache.Cache, error) {
	var c cache.Cache
	var err error
	switch opts.Driver {
	case cache.MemoryDriver:
		c = memory.New()
	case cache.RedisDriver:
		c, err = redisdrv.New(opts)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown cache driver: %s", opts.Driver)
	}
	// apply namespace
	c = cache.WithNamespace(c, opts.Namespace)
	// observability wrapper
	c = cache.WithObservability(c, string(opts.Driver), opts)
	return c, nil
}
