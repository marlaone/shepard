package http

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
)

func Serve(addr string, r *Router) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("[http.Serve] listen failed: %w", err)
	}
	defer l.Close()

	host, port, _ := strings.Cut(addr, ":")

	if host == "" {
		host = "localhost"
	}

	log.Println("Listening on http://" + host + ":" + port)

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

			potentialRes := HandleRequest(conn, req.Unwrap(), r)
			if potentialRes.IsErr() {
				conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
				// TODO remove this log
				log.Println(potentialRes.UnwrapErr())
				return
			}
			res := potentialRes.Unwrap()

			conn.Write([]byte("HTTP/" + res.Version().String() + " " + res.StatusCode().String() + "\r\n"))
			res.Headers().Iter().Foreach(func(_ int, value hashmap.Pair[*string, *slice.Slice[string]]) {
				headerValues := ""
				value.Value.Iter().Foreach(func(i int, value string) {
					if i > 0 {
						headerValues += ", "
					}
					headerValues += value
				})
				conn.Write([]byte(*value.Key + ": " + headerValues + "\r\n"))
			})
			conn.Write([]byte("\r\n"))
			buf := slice.New[byte]()
			for {
				res.Body().Read(&buf)

				iter := buf.Iter()
				next := iter.Next()
				data := make([]byte, 0, buf.Len())
				for next.IsSome() {
					data = append(data, next.Unwrap())
					next = iter.Next()
				}

				conn.Write(data)

				if res.Body().Closed() {
					break
				}
			}
			conn.Close()
		}(conn)
	}
}
