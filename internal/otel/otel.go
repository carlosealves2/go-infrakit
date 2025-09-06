package otel

import (
    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/trace"
)

// Tracer returns a no-op tracer.
func Tracer(name string) trace.Tracer { return trace.NoopTracer }

// Meter returns a no-op meter.
func Meter(name string) metric.Meter { return metric.Meter{} }
