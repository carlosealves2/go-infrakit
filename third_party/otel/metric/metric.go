package metric

import (
    "context"
    "go.opentelemetry.io/otel/attribute"
)

type Meter struct{}

type Int64Counter struct{}

type Float64Histogram struct{}

func (m Meter) Int64Counter(name string, opts ...interface{}) (Int64Counter, error) {
    return Int64Counter{}, nil
}

func (m Meter) Float64Histogram(name string, opts ...interface{}) (Float64Histogram, error) {
    return Float64Histogram{}, nil
}

func (c Int64Counter) Add(ctx context.Context, value int64, opts ...interface{}) {}

func (h Float64Histogram) Record(ctx context.Context, value float64, opts ...interface{}) {}

func WithAttributes(attrs ...attribute.KeyValue) interface{} { return nil }
