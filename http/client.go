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

// ThrottledRoundTripper wraps another RoundTripper rt, throttling all requests
// to the specified byte rate.
func ThrottledRoundTripper(rt http.RoundTripper, rate int64) http.RoundTripper {
	hz := time.Duration(1 * time.Millisecond)
	bucket := tb.NewBucket(rate, hz)

	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		var got int64
		for got < r.ContentLength {
			got += bucket.Take(r.ContentLength - got)
			time.Sleep(hz)
		}
		return rt.RoundTrip(r)
	})
}
