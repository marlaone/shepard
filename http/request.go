package http

import (
	"context"
	"strconv"
	"strings"

	"github.com/marlaone/shepard"
)

// Request represents an HTTP request (non-generic)
type Request struct {
	Context context.Context
	Method  Method
	URL     URL
	Headers *Headers
	Body    *BodyReader
}

// NewRequest creates a new Request with the given method, path, and body
func NewRequest(method Method, path string, body *BodyReader) Request {
	urlResult := ParseURL(path)
	url := URL{Path: path}
	if urlResult.IsOk() {
		url = urlResult.Unwrap()
	}

	return Request{
		Context: context.Background(),
		Method:  method,
		URL:     url,
		Headers: NewHeaders(),
		Body:    body,
	}
}

// WithContext returns a copy of the request with the given context
func (r Request) WithContext(ctx context.Context) Request {
	r.Context = ctx
	return r
}

// WithHeader returns a copy of the request with an added header
func (r Request) WithHeader(key, value string) Request {
	r.Headers.Set(key, value)
	return r
}

// Header returns the value of the given header
// Returns Option[string] - Some if found, None if not
func (r Request) Header(key string) shepard.Option[string] {
	return r.Headers.Get(key)
}

// Query returns the value of a query parameter
// Returns Option[string] - Some if found, None if not
func (r Request) Query(key string) shepard.Option[string] {
	queryResult := r.URL.Query()
	if queryResult.IsErr() {
		return shepard.None[string]()
	}
	qv := queryResult.Unwrap()
	result := qv.Get(key)
	if result.IsNone() {
		return shepard.None[string]()
	}
	return shepard.Some(*result.Unwrap())
}

// ContentType returns the Content-Type header value
func (r Request) ContentType() shepard.Option[string] {
	return r.Header(HeaderContentType)
}

// ContentLength returns the Content-Length header as int64
func (r Request) ContentLength() shepard.Option[int64] {
	lengthOpt := r.Header(HeaderContentLength)
	if lengthOpt.IsNone() {
		return shepard.None[int64]()
	}
	length, err := strconv.ParseInt(lengthOpt.Unwrap(), 10, 64)
	if err != nil {
		return shepard.None[int64]()
	}
	return shepard.Some(length)
}

// IsJSON returns true if the Content-Type indicates JSON
func (r Request) IsJSON() bool {
	ct := r.ContentType()
	if ct.IsNone() {
		return false
	}
	return strings.Contains(ct.Unwrap(), "application/json")
}

// IsForm returns true if the Content-Type indicates form data
func (r Request) IsForm() bool {
	ct := r.ContentType()
	if ct.IsNone() {
		return false
	}
	return strings.Contains(ct.Unwrap(), "application/x-www-form-urlencoded")
}

// IsMultipart returns true if the Content-Type indicates multipart form data
func (r Request) IsMultipart() bool {
	ct := r.ContentType()
	if ct.IsNone() {
		return false
	}
	return strings.Contains(ct.Unwrap(), "multipart/form-data")
}

// RequestBuilder provides a fluent API for building requests
type RequestBuilder struct {
	method  Method
	path    string
	headers *Headers
	body    *BodyReader
	ctx     context.Context
}

// NewRequestBuilder creates a new RequestBuilder
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		method:  MethodGet,
		path:    "/",
		headers: NewHeaders(),
		body:    EmptyBodyReader(),
		ctx:     context.Background(),
	}
}

// Method sets the HTTP method
func (b *RequestBuilder) Method(method Method) *RequestBuilder {
	b.method = method
	return b
}

// Path sets the request path
func (b *RequestBuilder) Path(path string) *RequestBuilder {
	b.path = path
	return b
}

// Header adds a header
func (b *RequestBuilder) Header(key, value string) *RequestBuilder {
	b.headers.Set(key, value)
	return b
}

// Body sets the request body
func (b *RequestBuilder) Body(body *BodyReader) *RequestBuilder {
	b.body = body
	return b
}

// Context sets the request context
func (b *RequestBuilder) Context(ctx context.Context) *RequestBuilder {
	b.ctx = ctx
	return b
}

// Build creates the Request
func (b *RequestBuilder) Build() Request {
	urlResult := ParseURL(b.path)
	url := URL{Path: b.path}
	if urlResult.IsOk() {
		url = urlResult.Unwrap()
	}

	return Request{
		Context: b.ctx,
		Method:  b.method,
		URL:     url,
		Headers: b.headers,
		Body:    b.body,
	}
}
