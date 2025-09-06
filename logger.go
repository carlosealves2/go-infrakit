package infrakit

import "github.com/carlosealves2/go-infrakit/observability/logger"

// NewLogger constructs a logger using the provided adapter.
func NewLogger(cfg logger.Config) logger.Logger {
	if cfg.Adapter == nil {
		return nil
	}
	return cfg.Adapter.New(cfg)
}
