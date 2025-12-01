package http

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
)

// Context wraps a Request with additional extracted data
// This is passed to handlers and provides convenience methods for accessing request data
type Context struct {
	Request    Request
	pathParams *hashmap.HashMap[string, string] // nil until first path param is set
	state      *hashmap.HashMap[string, any]    // nil until first state is set
}

// NewContext creates a new Context from a Request
func NewContext(req Request) *Context {
	return &Context{
		Request: req,
		// pathParams and state are lazy-initialized on first use
	}
}

// SetState stores a value in the context state
func (c *Context) SetState(key string, value any) {
	if c.state == nil {
		hm := hashmap.New[string, any]()
		c.state = &hm
	}
	c.state.Insert(key, value)
}

// GetState retrieves a value from the context state
func (c *Context) GetState(key string) shepard.Option[any] {
	if c.state == nil {
		return shepard.None[any]()
	}
	result := c.state.Get(key)
	if result.IsNone() {
		return shepard.None[any]()
	}
	return shepard.Some[any](*result.Unwrap())
}

// SetPathParam sets a path parameter (called by router during routing)
func (c *Context) SetPathParam(key, value string) {
	if c.pathParams == nil {
		hm := hashmap.New[string, string]()
		c.pathParams = &hm
	}
	c.pathParams.Insert(key, value)
}

// Param returns a path parameter by name
// Returns Option[string] - Some if found, None if not
func (c *Context) Param(name string) shepard.Option[string] {
	if c.pathParams == nil {
		return shepard.None[string]()
	}
	result := c.pathParams.Get(name)
	if result.IsNone() {
		return shepard.None[string]()
	}
	return shepard.Some(*result.Unwrap())
}

// ClearPathParams resets the path parameters for context reuse
func (c *Context) ClearPathParams() {
	if c.pathParams != nil {
		c.pathParams.Clear()
	}
}

// Query returns a query parameter by name
// Returns Option[string] - Some if found, None if not
func (c *Context) Query(name string) shepard.Option[string] {
	return c.Request.Query(name)
}

// Header returns a header value by name
// Returns Option[string] - Some if found, None if not
func (c *Context) Header(name string) shepard.Option[string] {
	return c.Request.Header(name)
}

// Method returns the HTTP method
func (c *Context) Method() Method {
	return c.Request.Method
}

// Path returns the request path
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// ContentType returns the Content-Type header
func (c *Context) ContentType() shepard.Option[string] {
	return c.Request.ContentType()
}

// IsJSON returns true if the Content-Type is JSON
func (c *Context) IsJSON() bool {
	return c.Request.IsJSON()
}

// IsForm returns true if the Content-Type is form data
func (c *Context) IsForm() bool {
	return c.Request.IsForm()
}

// Body returns the request body reader
func (c *Context) Body() *BodyReader {
	return c.Request.Body
}

// BodyJson deserializes the request body as JSON into type T
func ContextBodyJson[T any](c *Context) shepard.Result[T, error] {
	return BodyJson[T](c.Request.Body)
}

// ContextBodyForm parses URL-encoded form data from the request body
func ContextBodyForm(c *Context) shepard.Result[FormValues, error] {
	return BodyForm(c.Request.Body)
}

// ContextBodyMultipart parses multipart form data from the request body
func ContextBodyMultipart(c *Context) shepard.Result[MultipartForm, error] {
	contentType := c.ContentType()
	if contentType.IsNone() {
		return shepard.Err[MultipartForm, error](&multipartError{msg: "missing Content-Type header"})
	}
	return BodyMultipart(c.Request.Body, contentType.Unwrap())
}

// ContextBodyMultipartWithLimit parses multipart form data with a custom memory limit
func ContextBodyMultipartWithLimit(c *Context, maxMemory int64) shepard.Result[MultipartForm, error] {
	contentType := c.ContentType()
	if contentType.IsNone() {
		return shepard.Err[MultipartForm, error](&multipartError{msg: "missing Content-Type header"})
	}
	return BodyMultipartWithLimit(c.Request.Body, contentType.Unwrap(), maxMemory)
}
