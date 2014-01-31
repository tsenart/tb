package throttle

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestThrottler(t *testing.T) {
	th := &Throttler{buckets: map[string]*Bucket{}}
	var hits int64
	for i := 0; i < 10; i++ {
		go func() {
			for _ = range time.Tick(5 * time.Millisecond) {
				atomic.AddInt64(&hits, th.Throttle("a", 1, 5000))
			}
		}()
	}
	time.Sleep(1 * time.Second)
	t.Log("HITS:", hits)
}
