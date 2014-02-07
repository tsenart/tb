package http

import (
	"github.com/tsenart/tb"
	"net/http"
	"time"
)

type roundTripperFunc func(r *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// ByteThrottledRoundTripper wraps another RoundTripper rt,
// throttling all requests to the specified byte rate.
func ByteThrottledRoundTripper(rt http.RoundTripper, rate int64) http.RoundTripper {
	hz := time.Duration(1 * time.Millisecond)
	bucket := tb.NewBucket(rate, hz)

	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		got := bucket.Take(r.ContentLength)
		for got < r.ContentLength {
			got += bucket.Take(r.ContentLength - got)
			time.Sleep(hz)
		}
		return rt.RoundTrip(r)
	})
}

// ReqThrottledRoundTripper wraps another RoundTripper rt,
// throttling all requests to the specified request rate.
func ReqThrottledRoundTripper(rt http.RoundTripper, rate int64) http.RoundTripper {
	hz := time.Duration(1e9 / rate)
	bucket := tb.NewBucket(rate, hz)

	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		got := bucket.Take(1)
		for got != 1 {
			got = bucket.Take(1)
			time.Sleep(hz)
		}
		return rt.RoundTrip(r)
	})
}
