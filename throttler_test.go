package tb

import (
	"io"
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func BenchmarkThrottle_allocs(b *testing.B) {
	keys := make([]string, 10000)
	for i := 0; i < len(keys); i++ {
		keys[i] = strconv.Itoa(i)
	}

	th := NewThrottler(1 * time.Millisecond)
	defer th.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		th.Throttle(keys[i%(len(keys)-1)], 1, 1000)
	}
}

func BenchmarkThrottle_sequential(b *testing.B) {
	th := NewThrottler(1 * time.Millisecond)
	defer th.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		th.Throttle("1", 1, 1000)
	}
}

func ExampleThrottle(t *testing.T) {
	ln, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatal(err)
	}
	th := NewThrottler(100 * time.Millisecond)
	defer th.Close()

	echo := func(conn net.Conn) {
		defer conn.Close()

		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}
		// Throttle to 10 connection per second from the same host
		// Handle non-conformity by dropping the connection
		if out := th.Throttle(host, 1, 10); out < 1 {
			log.Printf("Throttled %s", host)
			return
		}
		log.Printf("Echoing payload from %s:%s", host, port)
		io.Copy(conn, conn)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go echo(conn)
	}
}

func TestThrottler_Close(t *testing.T) {
	th := NewThrottler(1 * time.Millisecond)
	for i := 0; i < 5; i++ {
		th.Throttle(strconv.Itoa(i), 1000, 1000)
	}
	th.Close()
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < 5; i++ {
		if w, g := int64(0), th.Throttle(strconv.Itoa(i), 1, 1000); w != g {
			t.Errorf("Want: %d Got: %d", w, g)
		}
	}
}
