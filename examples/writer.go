package examples

import (
	"github.com/tsenart/tb"
	"io"
	"time"
)

type ThrottledWriter struct {
	b *tb.Bucket
	w io.Writer
}

func NewThrottledWriter(rate int64, w io.Writer) io.Writer {
	return &ThrottledWriter{tb.NewBucket(rate, -1), w}
}

func (tw *ThrottledWriter) Write(p []byte) (n int, err error) {
	for wr := 0; wr < len(p); {
		in := len(p) - wr
		if out := tw.b.Take(int64(in)); out == 0 {
			time.Sleep(1 * time.Second)
			continue
		} else if n, err = tw.w.Write(p[wr : wr+int(out)]); err != nil {
			return wr, err
		}
		wr += n
	}
	return len(p), nil
}
