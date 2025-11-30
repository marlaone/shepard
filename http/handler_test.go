package http

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/stretchr/testify/assert"
)

func TestHandlerFunc_Handle(t *testing.T) {
	handler := HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
		return shepard.Ok[Response, HttpError](Text("Hello from handler"))
	})

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	result := handler.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello from handler", string(resp.Body.Bytes()))
}

func TestHandlerFunc_HandleError(t *testing.T) {
	handler := HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
		return shepard.Err[Response, HttpError](NotFound("resource not found"))
	})

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	result := handler.Handle(ctx)

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusNotFound, err.StatusCode)
}

func TestSimpleHandlerFunc_Handle(t *testing.T) {
	handler := SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("Simple response")
	})

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	result := handler.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, "Simple response", string(resp.Body.Bytes()))
}

func TestWriterHandlerFunc_Handle(t *testing.T) {
	handler := WriterHandlerFunc(func(ctx *Context, w *ResponseWriter) {
		w.WriteHeader(StatusCreated)
		w.Header().Set("X-Custom", "value")
		w.Write([]byte("Written directly"))
	})

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	result := handler.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusCreated, resp.StatusCode)
	assert.Equal(t, "Written directly", string(resp.Body.Bytes()))

	custom := resp.Headers.Get("X-Custom")
	assert.True(t, custom.IsSome())
	assert.Equal(t, "value", *custom.Unwrap())
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := NewResponseWriter()

	w.WriteHeader(StatusNotFound)

	result := w.IntoResponse()
	assert.True(t, result.IsOk())
	assert.Equal(t, StatusNotFound, result.Unwrap().StatusCode)
}

func TestResponseWriter_Write(t *testing.T) {
	w := NewResponseWriter()

	result := w.Write([]byte("Hello"))
	assert.True(t, result.IsOk())
	assert.Equal(t, 5, result.Unwrap())

	result = w.Write([]byte(", World!"))
	assert.True(t, result.IsOk())

	resp := w.IntoResponse()
	assert.True(t, resp.IsOk())
	assert.Equal(t, "Hello, World!", string(resp.Unwrap().Body.Bytes()))
}

func TestResponseWriter_DefaultStatus(t *testing.T) {
	w := NewResponseWriter()
	w.Write([]byte("content"))

	resp := w.IntoResponse()
	assert.True(t, resp.IsOk())
	// Default status should be 200 OK
	assert.Equal(t, StatusOK, resp.Unwrap().StatusCode)
}

func TestResponseWriter_Header(t *testing.T) {
	w := NewResponseWriter()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", "123")

	resp := w.IntoResponse()
	assert.True(t, resp.IsOk())

	ct := resp.Unwrap().Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/json", *ct.Unwrap())
}

func TestHandler_Interface(t *testing.T) {
	// All handler types should implement Handler interface
	var h Handler

	h = HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
		return shepard.Ok[Response, HttpError](Text("1"))
	})
	assert.NotNil(t, h)

	h = SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("2")
	})
	assert.NotNil(t, h)

	h = WriterHandlerFunc(func(ctx *Context, w *ResponseWriter) {
		w.Write([]byte("3"))
	})
	assert.NotNil(t, h)
}
