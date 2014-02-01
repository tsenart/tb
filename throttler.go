package tb

import (
	"sync"
)

// DefaultThrottler is an utility instance of Throttler used by Throttle.
var DefaultThrottler = &Throttler{buckets: map[string]*Bucket{}}

// Throttle throttles a quantity 'in' to the specified 'rate' per second,
// with a Bucket keyed by key, returning the permitted quantity.
func Throttle(key string, in, rate int64) (out int64) {
	return DefaultThrottler.Throttle(key, in, rate)
}

// Throttler is a thread-safe wrapper around a map of buckets and an easy to
// use API for generic throttling.
type Throttler struct {
	mu      sync.RWMutex
	buckets map[string]*Bucket
}

// NewThrottler returns an instance of Throttler
func NewThrottler() *Throttler {
	return &Throttler{buckets: map[string]*Bucket{}}
}

// Throttle throttles a quantity 'in' to the specified 'rate' per second,
// with a Bucket keyed by key, returning the permitted quantity.
// This method is thread-safe, locks are used only to synchronize access to
// the bucket map.
func (t *Throttler) Throttle(key string, in, rate int64) (out int64) {
	t.mu.RLock()
	b := t.buckets[key]
	t.mu.RUnlock()

	if b == nil {
		b = NewBucket(rate)
		t.mu.Lock()
		t.buckets[key] = b
		t.mu.Unlock()
	}

	return b.Take(in)
}

// Stop halts the addition of new tokens to all buckets. It is not thread-safe.
func (t *Throttler) Stop() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, b := range t.buckets {
		b.Stop()
	}
}
