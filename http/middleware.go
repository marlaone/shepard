package http

import (
	"fmt"
	"time"

	"github.com/marlaone/shepard"
)

// Logger creates a logging middleware
// The logFn receives a formatted log message for each request
func Logger(logFn func(msg string)) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
			start := time.Now()
			method := ctx.Method().String()
			path := ctx.Path()

			result := next.Handle(ctx)

			duration := time.Since(start)
			status := StatusOK
			if result.IsOk() {
				status = result.Unwrap().StatusCode
			}

			msg := fmt.Sprintf("%s %s -> %d (%v)", method, path, status.Int(), duration)
			logFn(msg)

			return result
		})
	}
}

// Recover creates a middleware that recovers from panics
// Returns a 500 Internal Server Error instead of crashing
func Recover() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) (result shepard.Result[Response, HttpError]) {
			defer func() {
				if r := recover(); r != nil {
					result = shepard.Ok[Response, HttpError](
						InternalErrorResponse(fmt.Sprintf("internal server error: %v", r)),
					)
				}
			}()

			return next.Handle(ctx)
		})
	}
}

// CORS creates a CORS middleware with the given allowed origins
func CORS(allowedOrigins ...string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
			// Handle preflight OPTIONS request
			if ctx.Method() == MethodOptions {
				resp := NewResponse(StatusNoContent)
				setCORSHeaders(&resp, allowedOrigins)
				return shepard.Ok[Response, HttpError](resp)
			}

			result := next.Handle(ctx)

			if result.IsOk() {
				resp := result.Unwrap()
				setCORSHeaders(&resp, allowedOrigins)
				return shepard.Ok[Response, HttpError](resp)
			}

			return result
		})
	}
}

func setCORSHeaders(resp *Response, allowedOrigins []string) {
	origin := "*"
	if len(allowedOrigins) > 0 {
		origin = allowedOrigins[0] // Simplified: just use first
	}

	resp.Headers.Set("Access-Control-Allow-Origin", origin)
	resp.Headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	resp.Headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// RequestID creates a middleware that adds a request ID header
func RequestID(generator func() string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
			id := generator()
			ctx.Request.Headers.Set("X-Request-ID", id)

			result := next.Handle(ctx)

			if result.IsOk() {
				resp := result.Unwrap()
				resp.Headers.Set("X-Request-ID", id)
				return shepard.Ok[Response, HttpError](resp)
			}

			return result
		})
	}
}

// Timeout creates a middleware that sets a timeout for the request context
// Note: This is a simplified version - full implementation would use context.WithTimeout
func Timeout(duration time.Duration) Middleware {
	return func(next Handler) Handler {
		// TODO: Implement proper timeout with context cancellation
		// For now, just a passthrough
		return next
	}
}

// Chain combines multiple middlewares into one
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		// Build the chain from inside out
		handler := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
