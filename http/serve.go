package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/marlaone/shepard"
)

// bufioReaderPool reuses bufio.Reader instances to reduce allocations
var bufioReaderPool = sync.Pool{
	New: func() any { return bufio.NewReaderSize(nil, 4096) },
}

// responseBufPool reuses byte buffers for writing HTTP responses
var responseBufPool = sync.Pool{
	New: func() any { return make([]byte, 0, 512) },
}

// maxPooledBufSize is the maximum buffer size to return to the pool
// Larger buffers are left for GC to prevent memory bloat
const maxPooledBufSize = 64 * 1024 // 64KB

// ParseHTTPRequest parses an HTTP/1.1 request from a reader
func ParseHTTPRequest(reader io.Reader) shepard.Result[Request, error] {
	// Note: bufioReaderPool is available but disabled for now
	// due to issues with POST requests where body isn't read
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
		contentLength, err := strconv.ParseInt(contentLengthOpt.Unwrap(), 10, 64)
		if err == nil && contentLength > 0 {
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
	// Get pooled buffer
	buf := responseBufPool.Get().([]byte)
	buf = buf[:0]

	// Return to pool when done (unless oversized)
	defer func() {
		if cap(buf) <= maxPooledBufSize {
			responseBufPool.Put(buf[:0])
		}
	}()

	// Build status line
	buf = append(buf, "HTTP/1.1 "...)
	buf = strconv.AppendInt(buf, int64(resp.StatusCode.Int()), 10)
	buf = append(buf, ' ')
	buf = append(buf, resp.StatusCode.Text()...)
	buf = append(buf, '\r', '\n')

	// Build headers
	resp.Headers.Foreach(func(key string, values []string) {
		for _, value := range values {
			buf = append(buf, key...)
			buf = append(buf, ':', ' ')
			buf = append(buf, value...)
			buf = append(buf, '\r', '\n')
		}
	})

	// Add Connection: close since we don't support keep-alive yet
	buf = append(buf, "Connection: close\r\n"...)

	// Add Content-Length if body exists and not already set
	bodyBytes := resp.Body.Bytes()
	if len(bodyBytes) > 0 {
		contentLengthOpt := resp.Headers.Get("Content-Length")
		if contentLengthOpt.IsNone() {
			buf = append(buf, "Content-Length: "...)
			buf = strconv.AppendInt(buf, int64(len(bodyBytes)), 10)
			buf = append(buf, '\r', '\n')
		}
	}

	// End headers
	buf = append(buf, '\r', '\n')

	// Write headers
	if _, err := w.Write(buf); err != nil {
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

	// Drain any unread body data before responding
	// This prevents EOF errors when client is still sending
	if req.Body != nil {
		io.Copy(io.Discard, req.Body.reader)
	}

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
