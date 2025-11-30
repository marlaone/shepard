package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_New(t *testing.T) {
	req := NewRequest(MethodGet, "/users/123", EmptyBodyReader())
	ctx := NewContext(req)

	assert.Equal(t, MethodGet, ctx.Request.Method)
	assert.Equal(t, "/users/123", ctx.Request.URL.Path)
}

func TestContext_Param(t *testing.T) {
	req := NewRequest(MethodGet, "/users/123", EmptyBodyReader())
	ctx := NewContext(req)

	// Set path params (normally done by router)
	ctx.SetPathParam("id", "123")
	ctx.SetPathParam("name", "john")

	id := ctx.Param("id")
	assert.True(t, id.IsSome())
	assert.Equal(t, "123", id.Unwrap())

	name := ctx.Param("name")
	assert.True(t, name.IsSome())
	assert.Equal(t, "john", name.Unwrap())

	missing := ctx.Param("missing")
	assert.True(t, missing.IsNone())
}

func TestContext_Query(t *testing.T) {
	req := NewRequest(MethodGet, "/search?q=hello&page=1", EmptyBodyReader())
	ctx := NewContext(req)

	q := ctx.Query("q")
	assert.True(t, q.IsSome())
	assert.Equal(t, "hello", q.Unwrap())

	page := ctx.Query("page")
	assert.True(t, page.IsSome())
	assert.Equal(t, "1", page.Unwrap())

	missing := ctx.Query("missing")
	assert.True(t, missing.IsNone())
}

func TestContext_Header(t *testing.T) {
	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	req.Headers.Set("Authorization", "Bearer token123")
	ctx := NewContext(req)

	auth := ctx.Header("Authorization")
	assert.True(t, auth.IsSome())
	assert.Equal(t, "Bearer token123", auth.Unwrap())

	missing := ctx.Header("X-Not-Found")
	assert.True(t, missing.IsNone())
}

func TestContext_Context(t *testing.T) {
	type ctxKey string
	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	// Set a value in the request context
	newReqCtx := context.WithValue(ctx.Request.Context, ctxKey("user"), "john")
	ctx.Request.Context = newReqCtx

	val := ctx.Request.Context.Value(ctxKey("user"))
	assert.Equal(t, "john", val)
}

func TestContext_Method(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())
	ctx := NewContext(req)

	assert.Equal(t, MethodPost, ctx.Method())
}

func TestContext_Path(t *testing.T) {
	req := NewRequest(MethodGet, "/api/users/123", EmptyBodyReader())
	ctx := NewContext(req)

	assert.Equal(t, "/api/users/123", ctx.Path())
}

func TestContext_ContentType(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())
	req.Headers.Set("Content-Type", "application/json")
	ctx := NewContext(req)

	ct := ctx.ContentType()
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/json", ct.Unwrap())
}

func TestContext_IsJSON(t *testing.T) {
	req := NewRequest(MethodPost, "/", EmptyBodyReader())
	req.Headers.Set("Content-Type", "application/json")
	ctx := NewContext(req)

	assert.True(t, ctx.IsJSON())
}

func TestContext_Body(t *testing.T) {
	body := NewBodyReaderFromBytes([]byte("test body"))
	req := NewRequest(MethodPost, "/", body)
	ctx := NewContext(req)

	data := ctx.Body().ReadAll()
	assert.True(t, data.IsOk())
	assert.Equal(t, "test body", string(data.Unwrap()))
}
