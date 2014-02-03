package tb

import (
	"sync/atomic"
	"time"
)

// Bucket defines a generic lock-free implementation of a Token Bucket.
type Bucket struct{ tokens int64 }

// NewBucket returns a full empty Bucket with c capacity and asynchronously
// starts filling it c times per second.
func NewBucket(c int64) *Bucket {
	b := &Bucket{tokens: c}
	go b.fill(c)
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

func (b *Bucket) fill(capacity int64) {
	for _ = range time.Tick(time.Duration(1e9 / capacity)) {
		if tokens := atomic.LoadInt64(&b.tokens); tokens < capacity {
			atomic.AddInt64(&b.tokens, 1)
		}
	}
}
