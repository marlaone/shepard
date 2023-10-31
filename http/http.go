package http

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/marlaone/shepard/collections/hashmap"
)

func Serve(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("[http.Serve] listen failed: %w", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("[http.Serve] accept failed: %w", err)
		}

		go func(conn net.Conn) {
			req := RequestFromConnection(conn)

			if req.IsErr() {
				conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
				// TODO remove this log
				log.Println(req.UnwrapErr())
				return
			}

			potentialRes := HandleRequest(conn, req.Unwrap())
			if potentialRes.IsErr() {
				conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
				// TODO remove this log
				log.Println(potentialRes.UnwrapErr())
				return
			}
			res := potentialRes.Unwrap()

			conn.Write([]byte("HTTP/" + res.Version().String() + " " + res.StatusCode().String() + "\r\n"))
			res.Headers().Iter().Foreach(func(_ int, value hashmap.Pair[*string, *[]string]) {
				values := strings.Join(*value.Value, ", ")
				conn.Write([]byte(*value.Key + ": " + values + "\r\n"))
			})
			conn.Write([]byte("\r\n"))
			for data := range res.Body().Read() {
				conn.Write(data)
			}
			conn.Close()
		}(conn)
	}
}
