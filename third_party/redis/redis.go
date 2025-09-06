package redis

import (
    "context"
    "crypto/tls"
    "errors"
    "fmt"
    "sync"
    "time"
)

var (
    // Nil is returned when a key does not exist.
    Nil = errors.New("redis: nil")
)

type Options struct {
    Addr     string
    DB       int
    Username string
    Password string
    TLSConfig *tls.Config
}

type Client struct {
    mu    sync.RWMutex
    store map[string]item
}

type item struct {
    val string
    exp time.Time
}

func NewClient(opts *Options) *Client {
    return &Client{store: make(map[string]item)}
}

type StatusCmd struct{ err error }
func (c *Client) Ping(ctx context.Context) *StatusCmd { return &StatusCmd{} }
func (s *StatusCmd) Err() error { return s.err }

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *StatusCmd {
    c.mu.Lock()
    defer c.mu.Unlock()
    var v string
    switch t := value.(type) {
    case string:
        v = t
    case []byte:
        v = string(t)
    default:
        v = fmt.Sprint(t)
    }
    it := item{val: v}
    if ttl > 0 {
        it.exp = time.Now().Add(ttl)
    }
    c.store[key] = it
    return &StatusCmd{}
}

type StringCmd struct {
    val string
    err error
}

func (c *Client) Get(ctx context.Context, key string) *StringCmd {
    c.mu.RLock()
    it, ok := c.store[key]
    c.mu.RUnlock()
    if !ok || ( !it.exp.IsZero() && time.Now().After(it.exp) ) {
        if ok {
            c.mu.Lock(); delete(c.store, key); c.mu.Unlock()
        }
        return &StringCmd{err: Nil}
    }
    return &StringCmd{val: it.val}
}

func (s *StringCmd) Result() (string, error) { return s.val, s.err }
func (s *StringCmd) Bytes() ([]byte, error) {
    if s.err != nil { return nil, s.err }
    return []byte(s.val), nil
}

type IntCmd struct {
    val int64
    err error
}

func (c *Client) Del(ctx context.Context, keys ...string) *IntCmd {
    c.mu.Lock()
    defer c.mu.Unlock()
    for _, k := range keys {
        delete(c.store, k)
    }
    return &IntCmd{val: 0}
}

func (i *IntCmd) Err() error { return i.err }
func (i *IntCmd) Result() (int64, error) { return i.val, i.err }

func (c *Client) Exists(ctx context.Context, keys ...string) *IntCmd {
    c.mu.RLock()
    defer c.mu.RUnlock()
    var count int64
    for _, k := range keys {
        it, ok := c.store[k]
        if ok && (it.exp.IsZero() || time.Now().Before(it.exp)) {
            count++
        }
    }
    return &IntCmd{val: count}
}
