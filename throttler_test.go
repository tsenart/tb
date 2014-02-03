package tb

import (
	"io"
	"log"
	"net"
	"strconv"
	"testing"
)

func BenchmarkThrottle_allocs(b *testing.B) {
	keys := make([]string, 10000)
	for i := 0; i < len(keys); i++ {
		keys[i] = strconv.Itoa(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Throttle(keys[i%(len(keys)-1)], 1, 1000)
	}
}

func BenchmarkThrottle_sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Throttle("1", 1, 1000)
	}
}

func ExampleThrottle(t *testing.T) {
	ln, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatal(err)
	}

	echo := func(conn net.Conn) {
		defer conn.Close()

		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}
		// Throttle to 10 connection per second from the same host
		// Handle non-conformity by dropping the connection
		if out := Throttle(host, 1, 10); out < 1 {
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
