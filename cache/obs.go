package cache

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// WithObservability wraps a Cache with optional logging, metrics and tracing.
// The logger, tracer and meter are provided via Options and can be nil to
// disable specific signals.
func WithObservability(c Cache, provider string, opts Options) Cache {
	if opts.Logger == nil && opts.Tracer == nil && opts.Meter == (metric.Meter{}) {
		return c
	}
	obs := &observedCache{
		cache:    c,
		provider: provider,
		ns:       opts.Namespace,
		logger:   opts.Logger,
		tracer:   opts.Tracer,
	}
	if opts.Meter != (metric.Meter{}) {
		obs.counter, _ = opts.Meter.Int64Counter("cache_ops_total")
		obs.latency, _ = opts.Meter.Float64Histogram("cache_latency_ms")
	}
	return obs
}

// observedCache adds observability around a Cache implementation.
type observedCache struct {
	cache    Cache
	provider string
	ns       string
	logger   Logger
	tracer   trace.Tracer
	counter  metric.Int64Counter
	latency  metric.Float64Histogram
}

func (o *observedCache) log(op string, keyLen int, dur time.Duration, err error) {
	if o.logger == nil {
		return
	}
	fields := map[string]any{
		"mod":      "cache",
		"provider": o.provider,
		"op":       op,
		"ns":       o.ns,
		"key_len":  keyLen,
		"dur_ms":   dur.Milliseconds(),
	}
	if err != nil {
		fields["err"] = err
	}
	o.logger.Log(fields)
}

func (o *observedCache) trace(ctx context.Context, op string, keyLen int, hit bool, err error, span trace.Span) {
	if o.tracer == nil {
		return
	}
	span.SetAttributes(
		attribute.String("cache.provider", o.provider),
		attribute.String("cache.namespace", o.ns),
		attribute.Int("cache.key_len", keyLen),
	)
	if op == "get" {
		span.SetAttributes(attribute.Bool("cache.hit", hit))
	}
	if err != nil {
		span.RecordError(err)
	}
	span.End()
}

func (o *observedCache) metrics(ctx context.Context, op string, hit bool, dur time.Duration) {
	if o.counter != (metric.Int64Counter{}) {
		attrs := []attribute.KeyValue{
			attribute.String("provider", o.provider),
			attribute.String("op", op),
		}
		if op == "get" {
			attrs = append(attrs, attribute.Bool("hit", hit))
		}
		o.counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	if o.latency != (metric.Float64Histogram{}) {
		o.latency.Record(ctx, float64(dur.Milliseconds()), metric.WithAttributes(
			attribute.String("provider", o.provider),
			attribute.String("op", op),
		))
	}
}

func (o *observedCache) Set(ctx context.Context, key, value string) error {
	keyLen := len(key)
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.set")
	}
	start := time.Now()
	err := o.cache.Set(ctx, key, value)
	dur := time.Since(start)
	o.log("set", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "set", keyLen, false, err, span)
	}
	o.metrics(ctx, "set", false, dur)
	return err
}

func (o *observedCache) Get(ctx context.Context, key string) (string, error) {
	keyLen := len(key)
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.get")
	}
	start := time.Now()
	val, err := o.cache.Get(ctx, key)
	dur := time.Since(start)
	hit := err == nil
	o.log("get", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "get", keyLen, hit, err, span)
	}
	o.metrics(ctx, "get", hit, dur)
	return val, err
}

func (o *observedCache) Del(ctx context.Context, keys ...string) error {
	keyLen := 0
	if len(keys) > 0 {
		keyLen = len(keys[0])
	}
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.del")
	}
	start := time.Now()
	err := o.cache.Del(ctx, keys...)
	dur := time.Since(start)
	o.log("del", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "del", keyLen, false, err, span)
	}
	o.metrics(ctx, "del", false, dur)
	return err
}

func (o *observedCache) Exists(ctx context.Context, key string) (bool, error) {
	keyLen := len(key)
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.exists")
	}
	start := time.Now()
	b, err := o.cache.Exists(ctx, key)
	dur := time.Since(start)
	o.log("exists", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "exists", keyLen, false, err, span)
	}
	o.metrics(ctx, "exists", false, dur)
	return b, err
}

func (o *observedCache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	keyLen := len(key)
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.set")
	}
	start := time.Now()
	err := o.cache.SetWithTTL(ctx, key, value, ttl)
	dur := time.Since(start)
	o.log("set", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "set", keyLen, false, err, span)
	}
	o.metrics(ctx, "set", false, dur)
	return err
}

func (o *observedCache) SetBytes(ctx context.Context, key string, value []byte) error {
	keyLen := len(key)
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.set")
	}
	start := time.Now()
	err := o.cache.SetBytes(ctx, key, value)
	dur := time.Since(start)
	o.log("set", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "set", keyLen, false, err, span)
	}
	o.metrics(ctx, "set", false, dur)
	return err
}

func (o *observedCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	keyLen := len(key)
	var span trace.Span
	if o.tracer != nil {
		ctx, span = o.tracer.Start(ctx, "cache.get")
	}
	start := time.Now()
	val, err := o.cache.GetBytes(ctx, key)
	dur := time.Since(start)
	hit := err == nil
	o.log("get", keyLen, dur, err)
	if span != nil {
		o.trace(ctx, "get", keyLen, hit, err, span)
	}
	o.metrics(ctx, "get", hit, dur)
	return val, err
}
