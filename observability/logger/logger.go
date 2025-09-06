package logger

import (
	"fmt"
	"io"

	phuslog "github.com/phuslu/log"
)

// Level represents log verbosity.
type Level int

const (
	Debug Level = iota
	Info
	Error
)

// Format represents the output encoding.
type Format int

const (
	Text Format = iota
	JSON
)

// Config configures a logger instance.
type Config struct {
	Adapter Adapter
	Level   Level
	Format  Format
	Out     io.Writer
	Err     io.Writer
	Caller  bool
	Stack   bool
}

// Adapter builds a Logger from the provided config.
type Adapter interface {
	New(cfg Config) Logger
}

// Logger is a leveled structured logger.
type Logger interface {
	Debug() Entry
	Info() Entry
	Error() Entry
}

// Entry accumulates fields for a log event.
type Entry interface {
	With(key string, val any) Entry
	Log(key string, val any)
}

// PhusluLogAdapter adapts github.com/phuslu/log.
type PhusluLogAdapter struct{}

// New builds a Logger using phuslu/log.
func (PhusluLogAdapter) New(cfg Config) Logger {
	l := &phuslog.Logger{}
	return &phusluLogger{l: l}
}

type phusluLogger struct{ l *phuslog.Logger }

func (p *phusluLogger) Debug() Entry { return &phusluEntry{e: p.l.Debug()} }
func (p *phusluLogger) Info() Entry  { return &phusluEntry{e: p.l.Info()} }
func (p *phusluLogger) Error() Entry { return &phusluEntry{e: p.l.Error()} }

type phusluEntry struct{ e *phuslog.Entry }

func (e *phusluEntry) With(key string, val any) Entry {
	switch v := val.(type) {
	case string:
		e.e = e.e.Str(key, v)
	case int:
		e.e = e.e.Int(key, v)
	case int64:
		e.e = e.e.Int64(key, v)
	case error:
		e.e = e.e.Err(v)
	default:
		e.e = e.e.Str(key, fmt.Sprint(v))
	}
	return e
}

func (e *phusluEntry) Log(key string, val any) {
	e.With(key, val)
	e.e.Msg("")
}
