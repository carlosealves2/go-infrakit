package tracing

import (
	"github.com/carlosealves2/go-infrakit/observability/logger"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig configures tracing providers.
type TracingConfig struct {
	Provider string
	Endpoint string
	Service  string
	Logger   logger.Logger
}

// Tracing holds a tracer instance.
type Tracing struct {
	tracer trace.Tracer
}

// NewTracing creates a tracing provider. It is a stubbed no-op implementation.
func NewTracing(cfg TracingConfig) (*Tracing, error) {
	return &Tracing{tracer: trace.NoopTracer}, nil
}

// Tracer returns the underlying tracer.
func (t *Tracing) Tracer() trace.Tracer { return t.tracer }
