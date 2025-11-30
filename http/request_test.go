package http

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_New(t *testing.T) {
	body := NewBodyReader(strings.NewReader("test body"))
	req := NewRequest(MethodGet, "/users", body)

	assert.Equal(t, MethodGet, req.Method)
	assert.Equal(t, "/users", req.URL.Path)
	assert.NotNil(t, req.Body)
	assert.NotNil(t, req.Headers)
	assert.NotNil(t, req.Context)
}

func TestRequest_WithContext(t *testing.T) {
	req := NewRequest(MethodGet, "/", EmptyBodyReader())

	type ctxKey string
	ctx := context.WithValue(context.Background(), ctxKey("user"), "john")

	req = req.WithContext(ctx)

	assert.Equal(t, "john", req.Context.Value(ctxKey("user")))
}

func TestRequest_WithHeader(t *testing.T) {
	req := NewRequest(MethodGet, "/", EmptyBodyReader())

	req = req.WithHeader("Content-Type", "application/json")

	result := req.Headers.Get("Content-Type")
	assert.True(t, result.IsSome())
	assert.Equal(t, "application/json", *result.Unwrap())
}

func TestRequest_Header(t *testing.T) {
	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	req.Headers.Set("Authorization", "Bearer token123")

	result := req.Header("Authorization")
	assert.True(t, result.IsSome())
	assert.Equal(t, "Bearer token123", result.Unwrap())

	result = req.Header("X-Not-Found")
	assert.True(t, result.IsNone())
}

func TestRequest_Query(t *testing.T) {
	body := EmptyBodyReader()
	req := NewRequest(MethodGet, "/search?q=hello&page=1", body)

	q := req.Query("q")
	assert.True(t, q.IsSome())
	assert.Equal(t, "hello", q.Unwrap())

	page := req.Query("page")
	assert.True(t, page.IsSome())
	assert.Equal(t, "1", page.Unwrap())

	notFound := req.Query("missing")
	assert.True(t, notFound.IsNone())
}

func TestRequest_ContentType(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())
	req.Headers.Set("Content-Type", "application/json; charset=utf-8")

	ct := req.ContentType()
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/json; charset=utf-8", ct.Unwrap())
}

func TestRequest_ContentLength(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())
	req.Headers.Set("Content-Length", "123")

	cl := req.ContentLength()
	assert.True(t, cl.IsSome())
	assert.Equal(t, int64(123), cl.Unwrap())
}

func TestRequest_ContentLength_Invalid(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())
	req.Headers.Set("Content-Length", "not-a-number")

	cl := req.ContentLength()
	assert.True(t, cl.IsNone())
}

func TestRequest_IsJSON(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())

	assert.False(t, req.IsJSON())

	req.Headers.Set("Content-Type", "application/json")
	assert.True(t, req.IsJSON())

	req.Headers.Set("Content-Type", "application/json; charset=utf-8")
	assert.True(t, req.IsJSON())
}

func TestRequest_IsForm(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())

	assert.False(t, req.IsForm())

	req.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.True(t, req.IsForm())
}

func TestRequestBuilder(t *testing.T) {
	req := NewRequestBuilder().
		Method(MethodPost).
		Path("/api/users").
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer token").
		Body(NewBodyReaderFromBytes([]byte(`{"name":"John"}`))).
		Build()

	assert.Equal(t, MethodPost, req.Method)
	assert.Equal(t, "/api/users", req.URL.Path)
	assert.True(t, req.Headers.Contains("Content-Type"))
	assert.True(t, req.Headers.Contains("Authorization"))
}
