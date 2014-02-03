package tb

import (
	"sync/atomic"
	"time"
)

// Bucket defines a generic lock-free implementation of a Token Bucket.
type Bucket struct {
	ticker   *time.Ticker
	capacity int64
	tokens   int64
}

// NewBucket returns a new full Bucket with c capacity and asynchronously
// starts filling it c times per second.
func NewBucket(c int64) *Bucket {
	b := &Bucket{time.NewTicker(time.Duration(1e9 / c)), c, c}
	go b.fill()
	return b
}

// Take will take n tokens out of the bucket. If there aren't enough
// tokens, the difference is returned. Otherwise n is returned.
// This method is thread-safe.
func (b *Bucket) Take(n int64) int64 {
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

// Stop halts the addition of new tokens to the bucket. It is not thread-safe.
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
