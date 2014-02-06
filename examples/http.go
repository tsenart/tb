package examples

import (
	"github.com/tsenart/tb"
	"net/http"
	"time"
)

type RoundTripperFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func ThrottledRoundTripper(rt http.RoundTripper, rate int64) http.RoundTripper {
	hz := time.Duration(1 * time.Millisecond)
	bucket := tb.NewBucket(rate, hz)

	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		var got int64
		for got < r.ContentLength {
			got += bucket.Take(r.ContentLength - got)
			time.Sleep(hz)
		}
		return rt.RoundTrip(r)
	})
}

func NewThrottledClient(rate int64) *http.Client {
	return &http.Client{
		Transport: ThrottledRoundTripper(http.DefaultTransport, rate),
	}
}
