package log

// Logger is a minimal stub implementing structured logging methods used in the project.
type Logger struct{}

// Entry represents a logging entry.
type Entry struct{}

// Log returns a new entry.
func (l *Logger) Log() *Entry { return &Entry{} }

func (e *Entry) Str(key, val string) *Entry { return e }
func (e *Entry) Int(key string, val int) *Entry { return e }
func (e *Entry) Int64(key string, val int64) *Entry { return e }
func (e *Entry) Err(err error) *Entry { return e }
func (e *Entry) Msg(msg string) {}
