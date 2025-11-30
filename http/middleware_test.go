package http

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_BeforeHandler(t *testing.T) {
	var called []string

	mw := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		called = append(called, "middleware:before")
		result := next()
		called = append(called, "middleware:after")
		return result
	}

	router := NewRouter().
		Use(mw).
		Get("/", SimpleHandlerFunc(func(ctx *Context) Response {
			called = append(called, "handler")
			return Text("ok")
		}))

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	router.Handle(ctx)

	assert.Equal(t, []string{"middleware:before", "handler", "middleware:after"}, called)
}

func TestMiddleware_ModifyRequest(t *testing.T) {
	mw := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		// Add a custom header before handling
		ctx.Request.Headers.Set("X-Request-ID", "123")
		return next()
	}

	router := NewRouter().
		Use(mw).
		Get("/", SimpleHandlerFunc(func(ctx *Context) Response {
			id := ctx.Header("X-Request-ID")
			if id.IsSome() {
				return Text("ID: " + id.Unwrap())
			}
			return Text("no id")
		}))

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	assert.Equal(t, "ID: 123", string(result.Unwrap().Body.Bytes()))
}

func TestMiddleware_ModifyResponse(t *testing.T) {
	mw := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		result := next()
		if result.IsOk() {
			resp := result.Unwrap()
			resp.Headers.Set("X-Powered-By", "Shepard")
			return shepard.Ok[Response, HttpError](resp)
		}
		return result
	}

	router := NewRouter().
		Use(mw).
		Get("/", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("ok")
		}))

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	powered := resp.Headers.Get("X-Powered-By")
	assert.True(t, powered.IsSome())
	assert.Equal(t, "Shepard", *powered.Unwrap())
}

func TestMiddleware_ShortCircuit(t *testing.T) {
	var handlerCalled bool

	authMiddleware := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		auth := ctx.Header("Authorization")
		if auth.IsNone() {
			// Short-circuit: don't call next()
			return shepard.Ok[Response, HttpError](TextWithStatus(StatusUnauthorized, "missing auth"))
		}
		return next()
	}

	router := NewRouter().
		Use(authMiddleware).
		Get("/protected", SimpleHandlerFunc(func(ctx *Context) Response {
			handlerCalled = true
			return Text("secret data")
		}))

	// Request without auth header
	req := NewRequest(MethodGet, "/protected", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusUnauthorized, resp.StatusCode)
	assert.False(t, handlerCalled)
}

func TestMiddleware_MultipleMiddlewares(t *testing.T) {
	var order []string

	mw1 := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		order = append(order, "mw1:before")
		result := next()
		order = append(order, "mw1:after")
		return result
	}

	mw2 := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		order = append(order, "mw2:before")
		result := next()
		order = append(order, "mw2:after")
		return result
	}

	router := NewRouter().
		Use(mw1).
		Use(mw2).
		Get("/", SimpleHandlerFunc(func(ctx *Context) Response {
			order = append(order, "handler")
			return Text("ok")
		}))

	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)

	router.Handle(ctx)

	// mw1 wraps mw2 which wraps handler
	assert.Equal(t, []string{"mw1:before", "mw2:before", "handler", "mw2:after", "mw1:after"}, order)
}

func TestLogger_Middleware(t *testing.T) {
	var logs []string

	logger := Logger(func(msg string) {
		logs = append(logs, msg)
	})

	router := NewRouter().
		Use(logger).
		Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("users")
		}))

	req := NewRequest(MethodGet, "/users", EmptyBodyReader())
	ctx := NewContext(req)

	router.Handle(ctx)

	assert.Len(t, logs, 1)
	assert.Contains(t, logs[0], "GET")
	assert.Contains(t, logs[0], "/users")
	assert.Contains(t, logs[0], "200")
}

func TestRecover_Middleware(t *testing.T) {
	router := NewRouter().
		Use(Recover()).
		Get("/panic", SimpleHandlerFunc(func(ctx *Context) Response {
			panic("something went wrong")
		}))

	req := NewRequest(MethodGet, "/panic", EmptyBodyReader())
	ctx := NewContext(req)

	// Should not panic, should return 500
	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusInternalServerError, resp.StatusCode)
}
