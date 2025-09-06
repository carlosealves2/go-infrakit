package infrakit

import (
	"fmt"

	"github.com/carlosealves2/go-infrakit/cache"
	"github.com/carlosealves2/go-infrakit/cache/memory"
	redisdrv "github.com/carlosealves2/go-infrakit/cache/redis"
)

// NewCache initializes a cache according to the provided options.
// It validates the driver and configures the selected provider.
func NewCache(opts cache.Options) (cache.Cache, error) {
	switch opts.Driver {
	case cache.MemoryDriver:
		return memory.New(opts), nil
	case cache.RedisDriver:
		return redisdrv.New(opts)
	default:
		return nil, fmt.Errorf("unknown cache driver: %s", opts.Driver)
	}
}
