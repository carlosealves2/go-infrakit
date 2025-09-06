package cache

import "github.com/phuslu/log"

// Logger is a minimal interface for structured logging. Implementations may
// forward the fields to any logging backend.
type Logger interface {
	Log(fields map[string]any)
}

// PhusluLogger adapts github.com/phuslu/log to the Logger interface.
type PhusluLogger struct{ l *log.Logger }

// NewPhusluLogger creates a Logger backed by phuslu/log.
func NewPhusluLogger(l *log.Logger) Logger { return &PhusluLogger{l: l} }

// Log forwards fields to the underlying phuslu logger. Only a subset of types
// used by the cache are handled.
func (p *PhusluLogger) Log(fields map[string]any) {
	if p == nil || p.l == nil {
		return
	}
	e := p.l.Log()
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			e = e.Str(k, val)
		case int:
			e = e.Int(k, val)
		case int64:
			e = e.Int64(k, val)
		case error:
			e = e.Err(val)
		}
	}
	e.Msg("")
}
