package trace

import (
    "context"
    "go.opentelemetry.io/otel/attribute"
)

type Tracer interface {
    Start(ctx context.Context, name string, opts ...interface{}) (context.Context, Span)
}

type Span interface {
    SetAttributes(...attribute.KeyValue)
    RecordError(error)
    End()
}

type noopTracer struct{}

type noopSpan struct{}

func (t noopTracer) Start(ctx context.Context, name string, opts ...interface{}) (context.Context, Span) {
    return ctx, noopSpan{}
}

func (noopSpan) SetAttributes(...attribute.KeyValue) {}
func (noopSpan) RecordError(error)        {}
func (noopSpan) End()                     {}

var NoopTracer Tracer = noopTracer{}
