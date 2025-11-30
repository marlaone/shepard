package http

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/stretchr/testify/assert"
)

func TestRouter_New(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router)
}

func TestRouter_Get(t *testing.T) {
	router := NewRouter().
		Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("users list")
		}))

	req := NewRequest(MethodGet, "/users", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "users list", string(resp.Body.Bytes()))
}

func TestRouter_Post(t *testing.T) {
	router := NewRouter().
		Post("/users", SimpleHandlerFunc(func(ctx *Context) Response {
			return TextWithStatus(StatusCreated, "user created")
		}))

	req := NewRequest(MethodPost, "/users", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusCreated, resp.StatusCode)
}

func TestRouter_PathParams(t *testing.T) {
	router := NewRouter().
		Get("/users/:id", SimpleHandlerFunc(func(ctx *Context) Response {
			id := ctx.Param("id")
			if id.IsSome() {
				return Text("User " + id.Unwrap())
			}
			return TextWithStatus(StatusBadRequest, "missing id")
		}))

	req := NewRequest(MethodGet, "/users/123", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, "User 123", string(resp.Body.Bytes()))
}

func TestRouter_NotFound(t *testing.T) {
	router := NewRouter().
		Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("users")
		}))

	req := NewRequest(MethodGet, "/posts", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusNotFound, resp.StatusCode)
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	router := NewRouter().
		Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("users")
		}))

	req := NewRequest(MethodPost, "/users", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusMethodNotAllowed, resp.StatusCode)
}

func TestRouter_CustomNotFoundHandler(t *testing.T) {
	router := NewRouter().
		NotFound(SimpleHandlerFunc(func(ctx *Context) Response {
			return TextWithStatus(StatusNotFound, "custom 404")
		})).
		Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("users")
		}))

	req := NewRequest(MethodGet, "/posts", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, "custom 404", string(resp.Body.Bytes()))
}

func TestRouter_AllMethods(t *testing.T) {
	handler := SimpleHandlerFunc(func(ctx *Context) Response {
		return Text(ctx.Method().String())
	})

	router := NewRouter().
		Get("/test", handler).
		Post("/test", handler).
		Put("/test", handler).
		Delete("/test", handler).
		Patch("/test", handler).
		Head("/test", handler).
		Options("/test", handler)

	methods := []Method{MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch, MethodHead, MethodOptions}

	for _, method := range methods {
		t.Run(method.String(), func(t *testing.T) {
			req := NewRequest(method, "/test", EmptyBodyReader())
			ctx := NewContext(req)

			result := router.Handle(ctx)

			assert.True(t, result.IsOk())
			resp := result.Unwrap()
			assert.Equal(t, StatusOK, resp.StatusCode)
			// HEAD doesn't return body
			if method != MethodHead {
				assert.Equal(t, method.String(), string(resp.Body.Bytes()))
			}
		})
	}
}

func TestRouter_Route(t *testing.T) {
	router := NewRouter().
		Route(MethodGet, "/custom", SimpleHandlerFunc(func(ctx *Context) Response {
			return Text("custom route")
		}))

	req := NewRequest(MethodGet, "/custom", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	assert.Equal(t, "custom route", string(result.Unwrap().Body.Bytes()))
}

func TestRouter_HandlerReturnsError(t *testing.T) {
	router := NewRouter().
		Get("/fail", HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
			return shepard.Err[Response, HttpError](BadRequest("validation failed"))
		}))

	req := NewRequest(MethodGet, "/fail", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	// Router should convert the error to a response
	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusBadRequest, resp.StatusCode)
}

// --- RouteGroup Tests ---

func TestRouteGroup_BasicPrefix(t *testing.T) {
	router := NewRouter()

	api := router.Group("/api")
	api.Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("users list")
	}))

	req := NewRequest(MethodGet, "/api/users", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "users list", string(resp.Body.Bytes()))
}

func TestRouteGroup_AllMethods(t *testing.T) {
	router := NewRouter()

	api := router.Group("/api")
	api.Get("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("GET") }))
	api.Post("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("POST") }))
	api.Put("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("PUT") }))
	api.Delete("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("DELETE") }))
	api.Patch("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("PATCH") }))
	api.Head("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("HEAD") }))
	api.Options("/test", SimpleHandlerFunc(func(ctx *Context) Response { return Text("OPTIONS") }))

	tests := []struct {
		method   Method
		expected string
	}{
		{MethodGet, "GET"},
		{MethodPost, "POST"},
		{MethodPut, "PUT"},
		{MethodDelete, "DELETE"},
		{MethodPatch, "PATCH"},
		{MethodOptions, "OPTIONS"},
	}

	for _, tc := range tests {
		t.Run(tc.method.String(), func(t *testing.T) {
			req := NewRequest(tc.method, "/api/test", EmptyBodyReader())
			ctx := NewContext(req)

			result := router.Handle(ctx)
			assert.True(t, result.IsOk())
			assert.Equal(t, tc.expected, string(result.Unwrap().Body.Bytes()))
		})
	}
}

func TestRouteGroup_WithPathParams(t *testing.T) {
	router := NewRouter()

	api := router.Group("/api/v1")
	api.Get("/users/:id", SimpleHandlerFunc(func(ctx *Context) Response {
		id := ctx.Param("id")
		if id.IsSome() {
			return Text("User: " + id.Unwrap())
		}
		return TextWithStatus(StatusBadRequest, "missing id")
	}))

	req := NewRequest(MethodGet, "/api/v1/users/42", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	assert.Equal(t, "User: 42", string(result.Unwrap().Body.Bytes()))
}

func TestRouteGroup_WithMiddleware(t *testing.T) {
	router := NewRouter()

	var middlewareCalled bool
	authMiddleware := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		middlewareCalled = true
		return next()
	}

	api := router.Group("/api")
	api.Use(authMiddleware)
	api.Get("/protected", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("protected resource")
	}))

	req := NewRequest(MethodGet, "/api/protected", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	assert.True(t, middlewareCalled, "middleware should have been called")
	assert.Equal(t, "protected resource", string(result.Unwrap().Body.Bytes()))
}

func TestRouteGroup_MiddlewareRejects(t *testing.T) {
	router := NewRouter()

	authMiddleware := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		// Reject - don't call next
		return shepard.Ok[Response, HttpError](TextWithStatus(StatusUnauthorized, "unauthorized"))
	}

	api := router.Group("/api")
	api.Use(authMiddleware)
	api.Get("/protected", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("should not reach here")
	}))

	req := NewRequest(MethodGet, "/api/protected", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusUnauthorized, resp.StatusCode)
}

func TestRouteGroup_MultipleMiddleware(t *testing.T) {
	router := NewRouter()

	var order []string

	mw1 := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		order = append(order, "mw1-before")
		result := next()
		order = append(order, "mw1-after")
		return result
	}

	mw2 := func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		order = append(order, "mw2-before")
		result := next()
		order = append(order, "mw2-after")
		return result
	}

	api := router.Group("/api")
	api.Use(mw1)
	api.Use(mw2)
	api.Get("/test", SimpleHandlerFunc(func(ctx *Context) Response {
		order = append(order, "handler")
		return Text("ok")
	}))

	req := NewRequest(MethodGet, "/api/test", EmptyBodyReader())
	ctx := NewContext(req)

	router.Handle(ctx)

	// Middleware should be applied in order: mw1 wraps mw2 wraps handler
	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	assert.Equal(t, expected, order)
}

func TestRouteGroup_NestedGroups(t *testing.T) {
	router := NewRouter()

	api := router.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("v1 users")
	}))

	v2 := api.Group("/v2")
	v2.Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("v2 users")
	}))

	// Test v1
	req1 := NewRequest(MethodGet, "/api/v1/users", EmptyBodyReader())
	ctx1 := NewContext(req1)
	result1 := router.Handle(ctx1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, "v1 users", string(result1.Unwrap().Body.Bytes()))

	// Test v2
	req2 := NewRequest(MethodGet, "/api/v2/users", EmptyBodyReader())
	ctx2 := NewContext(req2)
	result2 := router.Handle(ctx2)
	assert.True(t, result2.IsOk())
	assert.Equal(t, "v2 users", string(result2.Unwrap().Body.Bytes()))
}

func TestRouteGroup_NestedMiddleware(t *testing.T) {
	router := NewRouter()

	var calls []string

	api := router.Group("/api")
	api.Use(func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		calls = append(calls, "api-mw")
		return next()
	})

	v1 := api.Group("/v1")
	v1.Use(func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		calls = append(calls, "v1-mw")
		return next()
	})

	v1.Get("/test", SimpleHandlerFunc(func(ctx *Context) Response {
		calls = append(calls, "handler")
		return Text("ok")
	}))

	req := NewRequest(MethodGet, "/api/v1/test", EmptyBodyReader())
	ctx := NewContext(req)

	router.Handle(ctx)

	// Both middlewares should be called
	assert.Contains(t, calls, "api-mw")
	assert.Contains(t, calls, "v1-mw")
	assert.Contains(t, calls, "handler")
}

func TestRouteGroup_NoMiddleware(t *testing.T) {
	router := NewRouter()

	// Group without middleware
	api := router.Group("/api")
	api.Get("/test", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("no middleware")
	}))

	req := NewRequest(MethodGet, "/api/test", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	assert.Equal(t, "no middleware", string(result.Unwrap().Body.Bytes()))
}

func TestRouteGroup_ChainedMethods(t *testing.T) {
	router := NewRouter()

	// Test fluent API
	router.Group("/api").
		Use(func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
			return next()
		}).
		Get("/users", SimpleHandlerFunc(func(ctx *Context) Response { return Text("users") })).
		Post("/users", SimpleHandlerFunc(func(ctx *Context) Response { return Text("create") })).
		Get("/posts", SimpleHandlerFunc(func(ctx *Context) Response { return Text("posts") }))

	// Test all routes work
	tests := []struct {
		method Method
		path   string
		body   string
	}{
		{MethodGet, "/api/users", "users"},
		{MethodPost, "/api/users", "create"},
		{MethodGet, "/api/posts", "posts"},
	}

	for _, tc := range tests {
		req := NewRequest(tc.method, tc.path, EmptyBodyReader())
		ctx := NewContext(req)
		result := router.Handle(ctx)
		assert.True(t, result.IsOk())
		assert.Equal(t, tc.body, string(result.Unwrap().Body.Bytes()))
	}
}

func TestRouteGroup_Route(t *testing.T) {
	router := NewRouter()

	api := router.Group("/api")
	api.Route(MethodGet, "/custom", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("custom route")
	}))

	req := NewRequest(MethodGet, "/api/custom", EmptyBodyReader())
	ctx := NewContext(req)

	result := router.Handle(ctx)

	assert.True(t, result.IsOk())
	assert.Equal(t, "custom route", string(result.Unwrap().Body.Bytes()))
}
