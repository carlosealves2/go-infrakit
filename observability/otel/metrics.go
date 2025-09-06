package otel

import (
	"github.com/carlosealves2/go-infrakit/observability/logger"
	"go.opentelemetry.io/otel/metric"
)

// MetricsConfig configures Metrics.
type MetricsConfig struct {
	Namespace string
	Server    string
	Logger    logger.Logger
}

// Metrics holds a Meter instance.
type Metrics struct {
	meter metric.Meter
}

// NewMetrics creates a metrics provider. It is a stubbed no-op implementation.
func NewMetrics(cfg MetricsConfig) (*Metrics, error) {
	return &Metrics{meter: metric.Meter{}}, nil
}

// Meter returns the underlying meter.
func (m *Metrics) Meter() metric.Meter { return m.meter }
