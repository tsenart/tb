package tb

import (
	"sync/atomic"
	"time"
)

// Bucket defines a generic lock-free implementation of a Token Bucket.
type Bucket struct {
	tokens  int64
	closing chan struct{}
}

// NewBucket returns a full Bucket with c capacity and asynchronously
// starts filling it c times per second. You must call Close when you're done
// with the Bucket in order to not leak a go-routine and a system timer.
func NewBucket(c int64) *Bucket {
	b := &Bucket{tokens: c, closing: make(chan struct{})}
	go b.fill(c)
	return b
}

// Take will attempt to take n tokens out of the bucket.
// If available tokens == 0, nothing will be taken.
// If n <= available tokens, n tokens will be taken.
// If n > available tokens, all available tokens will be taken.
//
// This method is thread-safe.
func (b *Bucket) Take(n int64) (taken int64) {
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

// Close halts filling the bucket with tokens and finishes execution of the
// respective go-routine.
func (b *Bucket) Close() error {
	close(b.closing)
	return nil
}

func (b *Bucket) fill(capacity int64) {
	ticker := time.NewTicker(time.Duration(1e9 / capacity))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if atomic.LoadInt64(&b.tokens) < capacity {
				atomic.AddInt64(&b.tokens, 1)
			}
		case <-b.closing:
			return
		default:
		}
	}
}
