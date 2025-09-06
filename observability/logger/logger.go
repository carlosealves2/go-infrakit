package logger

import (
	"io"
	"time"

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
	Str(key, val string) Entry
	Int(key string, val int) Entry
	Int64(key string, val int64) Entry
	Float64(key string, val float64) Entry
	Bool(key string, val bool) Entry
	Dur(key string, val time.Duration) Entry
	Time(key string, val time.Time) Entry
	Err(err error) Entry
	Msg(msg string)
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

func (e *phusluEntry) Str(key, val string) Entry { e.e = e.e.Str(key, val); return e }
func (e *phusluEntry) Int(key string, val int) Entry {
	e.e = e.e.Int(key, val)
	return e
}
func (e *phusluEntry) Int64(key string, val int64) Entry {
	e.e = e.e.Int64(key, val)
	return e
}
func (e *phusluEntry) Float64(key string, val float64) Entry {
	e.e = e.e.Float64(key, val)
	return e
}
func (e *phusluEntry) Bool(key string, val bool) Entry {
	e.e = e.e.Bool(key, val)
	return e
}
func (e *phusluEntry) Dur(key string, val time.Duration) Entry {
	e.e = e.e.Dur(key, val)
	return e
}
func (e *phusluEntry) Time(key string, val time.Time) Entry {
	e.e = e.e.Time(key, val)
	return e
}
func (e *phusluEntry) Err(err error) Entry { e.e = e.e.Err(err); return e }
func (e *phusluEntry) Msg(msg string)      { e.e.Msg(msg) }
