package http

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseRequest_Simple(t *testing.T) {
	raw := "GET /users HTTP/1.1\r\nHost: localhost\r\n\r\n"
	reader := strings.NewReader(raw)

	result := ParseHTTPRequest(reader)

	assert.True(t, result.IsOk())
	req := result.Unwrap()

	assert.Equal(t, MethodGet, req.Method)
	assert.Equal(t, "/users", req.URL.Path)
}

func TestParseRequest_WithHeaders(t *testing.T) {
	raw := "POST /api/users HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Type: application/json\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		`{"name":"Jo"}`

	reader := strings.NewReader(raw)

	result := ParseHTTPRequest(reader)

	assert.True(t, result.IsOk())
	req := result.Unwrap()

	assert.Equal(t, MethodPost, req.Method)
	assert.Equal(t, "/api/users", req.URL.Path)

	ct := req.Header("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/json", ct.Unwrap())
}

func TestParseRequest_WithBody(t *testing.T) {
	body := `{"name":"John"}`
	raw := "POST /users HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Length: " + string(rune(len(body))) + "\r\n" +
		"\r\n" +
		body

	reader := strings.NewReader(raw)

	result := ParseHTTPRequest(reader)

	assert.True(t, result.IsOk())
	req := result.Unwrap()

	// Read the body
	bodyResult := req.Body.ReadAll()
	assert.True(t, bodyResult.IsOk())
	// Note: body parsing is simplified for now
}

func TestParseRequest_InvalidMethod(t *testing.T) {
	raw := "INVALID /users HTTP/1.1\r\n\r\n"
	reader := strings.NewReader(raw)

	result := ParseHTTPRequest(reader)

	assert.True(t, result.IsErr())
}

func TestWriteResponse(t *testing.T) {
	resp := Text("Hello, World!")

	var buf strings.Builder
	err := WriteHTTPResponse(&buf, resp)

	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "HTTP/1.1 200 OK")
	assert.Contains(t, output, "Hello, World!")
}

func TestWriteResponse_WithHeaders(t *testing.T) {
	resp := NewResponse(StatusCreated).
		WithHeader("X-Custom", "value").
		WithBody([]byte("created"))

	var buf strings.Builder
	err := WriteHTTPResponse(&buf, resp)

	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "HTTP/1.1 201 Created")
	assert.Contains(t, output, "X-Custom: value")
	assert.Contains(t, output, "created")
}

func TestServer_Integration(t *testing.T) {
	// Create a router
	router := NewRouter().
		Get("/", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("Welcome!")
		})).
		Get("/users/:id", SimpleHandlerFunc(func(ctx *Context) Response {
			id := ctx.Param("id").UnwrapOr("unknown")
			return Text("User: " + id)
		}))

	// Start server on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	// Run server in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn, router)
		}
	}()

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Make a test request
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send HTTP request
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	assert.NoError(t, err)

	// Read response - may need multiple reads to get complete response
	var response strings.Builder
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second))
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			response.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	responseStr := response.String()
	assert.Contains(t, responseStr, "HTTP/1.1 200 OK")
	assert.Contains(t, responseStr, "Welcome!")
}

func TestServer_PathParams_Integration(t *testing.T) {
	router := NewRouter().
		Get("/users/:id", SimpleHandlerFunc(func(ctx *Context) Response {
			id := ctx.Param("id").UnwrapOr("unknown")
			return Text("User: " + id)
		}))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn, router)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("GET /users/123 HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	assert.NoError(t, err)

	// Read response - may need multiple reads to get complete response
	var response strings.Builder
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second))
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			response.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	responseStr := response.String()
	assert.Contains(t, responseStr, "HTTP/1.1 200 OK")
	assert.Contains(t, responseStr, "User: 123")
}
