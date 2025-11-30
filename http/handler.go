package http

import (
	"github.com/marlaone/shepard"
)

// Handler is the core interface for request handlers
// Any type that can handle a request and produce a Result implements this
type Handler interface {
	Handle(ctx *Context) shepard.Result[Response, HttpError]
}

// HandlerFunc is the most common handler type that returns Result[Response, HttpError]
// This is the recommended handler type for most use cases
type HandlerFunc func(ctx *Context) shepard.Result[Response, HttpError]

func (f HandlerFunc) Handle(ctx *Context) shepard.Result[Response, HttpError] {
	return f(ctx)
}

// SimpleHandlerFunc is for handlers that just return a Response (no error case)
// Useful for simple handlers that can't fail
type SimpleHandlerFunc func(ctx *Context) Response

func (f SimpleHandlerFunc) Handle(ctx *Context) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](f(ctx))
}

// WriterHandlerFunc is for handlers that write directly to a ResponseWriter
// Useful for streaming responses, SSE, or when you need fine-grained control
type WriterHandlerFunc func(ctx *Context, w *ResponseWriter)

func (f WriterHandlerFunc) Handle(ctx *Context) shepard.Result[Response, HttpError] {
	w := NewResponseWriter()
	f(ctx, w)
	return w.IntoResponse()
}

// ResponseWriter allows handlers to write responses directly
// Similar to http.ResponseWriter but with Result-based error handling
type ResponseWriter struct {
	statusCode StatusCode
	headers    *Headers
	body       *BodyWriter
	headersSent bool
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		statusCode: StatusOK, // Default to 200 OK
		headers:    NewHeaders(),
		body:       NewBodyWriter(),
		headersSent: false,
	}
}

// WriteHeader sets the status code
// Must be called before Write if you want a non-200 status
func (w *ResponseWriter) WriteHeader(status StatusCode) {
	if !w.headersSent {
		w.statusCode = status
	}
}

// Header returns the headers for modification
func (w *ResponseWriter) Header() *Headers {
	return w.headers
}

// Write writes data to the response body
func (w *ResponseWriter) Write(data []byte) shepard.Result[int, error] {
	w.headersSent = true
	return w.body.Write(data)
}

// WriteString writes a string to the response body
func (w *ResponseWriter) WriteString(s string) shepard.Result[int, error] {
	w.headersSent = true
	return w.body.WriteString(s)
}

// Flush is a no-op for now (buffered writer)
// In the future, this could support streaming
func (w *ResponseWriter) Flush() {
	// TODO: Implement streaming support
}

// IntoResponse converts the ResponseWriter into a Response
func (w *ResponseWriter) IntoResponse() shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](Response{
		StatusCode: w.statusCode,
		Headers:    w.headers,
		Body:       w.body,
	})
}

// --- Handler adapters for common patterns ---

// IntoResponseHandler wraps a function that returns IntoResponse
func IntoResponseHandler(f func(ctx *Context) IntoResponse) Handler {
	return HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
		return shepard.Ok[Response, HttpError](f(ctx).IntoResponse())
	})
}

// FromFunc creates a Handler from any function that returns Result[IntoResponse, HttpError]
func FromFunc[T IntoResponse](f func(ctx *Context) shepard.Result[T, HttpError]) Handler {
	return HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
		result := f(ctx)
		if result.IsErr() {
			return shepard.Err[Response, HttpError](result.UnwrapErr())
		}
		return shepard.Ok[Response, HttpError](result.Unwrap().IntoResponse())
	})
}
