package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/marlaone/shepard"
)

// ParseHTTPRequest parses an HTTP/1.1 request from a reader
func ParseHTTPRequest(reader io.Reader) shepard.Result[Request, error] {
	bufReader := bufio.NewReader(reader)

	// Read request line
	requestLine, err := bufReader.ReadString('\n')
	if err != nil {
		return shepard.Err[Request, error](fmt.Errorf("failed to read request line: %w", err))
	}
	requestLine = strings.TrimSuffix(requestLine, "\r\n")
	requestLine = strings.TrimSuffix(requestLine, "\n")

	// Parse request line: METHOD PATH HTTP/VERSION
	parts := strings.SplitN(requestLine, " ", 3)
	if len(parts) != 3 {
		return shepard.Err[Request, error](fmt.Errorf("invalid request line: %s", requestLine))
	}

	methodResult := MethodFromString(parts[0])
	if methodResult.IsErr() {
		return shepard.Err[Request, error](fmt.Errorf("invalid method: %s", parts[0]))
	}
	method := methodResult.Unwrap()
	path := parts[1]

	// Parse headers
	headers := NewHeaders()
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return shepard.Err[Request, error](fmt.Errorf("failed to read header: %w", err))
		}
		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")

		// Empty line marks end of headers
		if line == "" {
			break
		}

		// Parse header: Key: Value
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue // Skip malformed headers
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])
		headers.Set(key, value)
	}

	// Check for Content-Length to read body
	var bodyReader *BodyReader
	contentLengthOpt := headers.Get("Content-Length")
	if contentLengthOpt.IsSome() {
		contentLength, err := strconv.ParseInt(*contentLengthOpt.Unwrap(), 10, 64)
		if err == nil && contentLength > 0 {
			// Read body with limit
			bodyReader = NewBodyReaderWithLimit(bufReader, contentLength)
		} else {
			bodyReader = EmptyBodyReader()
		}
	} else {
		bodyReader = EmptyBodyReader()
	}

	req := NewRequest(method, path, bodyReader)
	req.Headers = headers

	return shepard.Ok[Request, error](req)
}

// WriteHTTPResponse writes an HTTP/1.1 response to a writer
func WriteHTTPResponse(w io.Writer, resp Response) error {
	// Write status line
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.StatusCode.Int(), resp.StatusCode.Text())
	if _, err := w.Write([]byte(statusLine)); err != nil {
		return err
	}

	// Write headers
	resp.Headers.Foreach(func(key string, values []string) {
		for _, value := range values {
			header := fmt.Sprintf("%s: %s\r\n", key, value)
			w.Write([]byte(header))
		}
	})

	// Write Content-Length if body exists
	bodyBytes := resp.Body.Bytes()
	if len(bodyBytes) > 0 {
		contentLengthOpt := resp.Headers.Get("Content-Length")
		if contentLengthOpt.IsNone() {
			contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(bodyBytes))
			w.Write([]byte(contentLength))
		}
	}

	// End headers
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}

	// Write body
	if len(bodyBytes) > 0 {
		if _, err := w.Write(bodyBytes); err != nil {
			return err
		}
	}

	return nil
}

// handleConnection handles a single HTTP connection
func handleConnection(conn net.Conn, router *Router) {
	defer conn.Close()

	// Parse the request
	result := ParseHTTPRequest(conn)
	if result.IsErr() {
		// Write error response
		resp := InternalErrorResponse("failed to parse request")
		WriteHTTPResponse(conn, resp)
		return
	}

	req := result.Unwrap()
	ctx := NewContext(req)

	// Handle the request
	handlerResult := router.Handle(ctx)

	var resp Response
	if handlerResult.IsOk() {
		resp = handlerResult.Unwrap()
	} else {
		resp = handlerResult.UnwrapErr().IntoResponse()
	}

	// Write response
	WriteHTTPResponse(conn, resp)
}

// Server represents an HTTP server
type Server struct {
	addr    string
	router  *Router
	running bool
}

// NewServer creates a new HTTP server
func NewServer(addr string, router *Router) *Server {
	return &Server{
		addr:   addr,
		router: router,
	}
}

// Serve starts the HTTP server and blocks until it's stopped
func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	s.running = true
	for s.running {
		conn, err := listener.Accept()
		if err != nil {
			if !s.running {
				return nil
			}
			continue
		}
		go handleConnection(conn, s.router)
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	s.running = false
}
