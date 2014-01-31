package throttle

import (
	"sync"
	"sync/atomic"
	"time"
)

type Throttler struct {
	mu      sync.RWMutex
	buckets map[string]*Bucket
}

func (t *Throttler) Throttle(key string, in, rate int64) (out int64) {
	var b *Bucket

	t.mu.Lock()
	if b = t.buckets[key]; b == nil {
		b = NewBucket(rate)
		t.buckets[key] = b
	}
	t.mu.Unlock()

	return b.Put(in)
}

func (t *Throttler) Stop() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, b := range t.buckets {
		b.Stop()
	}
}

type Bucket struct {
	ticker   *time.Ticker
	capacity int64
	tokens   int64
}

func NewBucket(capacity int64) *Bucket {
	b := &Bucket{time.NewTicker(time.Duration(1e9 / capacity)), capacity, 0}
	go b.fill()
	return b
}

func (b *Bucket) Put(n int64) int64 {
	for {
		if tokens := atomic.LoadInt64(&b.tokens); tokens == 0 {
			return 0
		} else if n <= tokens {
			if !atomic.CompareAndSwapInt64(&b.tokens, tokens, tokens-n) {
				continue
			}
			return n
		} else if atomic.CompareAndSwapInt64(&b.tokens, tokens, 0) { // Spill
			return tokens
		}
	}
}

func (b Bucket) Stop() {
	b.ticker.Stop()
}

func (b *Bucket) fill() {
	for _ = range b.ticker.C {
		if tokens := atomic.LoadInt64(&b.tokens); tokens < b.capacity {
			atomic.AddInt64(&b.tokens, 1)
		}
	}
}
