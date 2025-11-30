package http

import (
	"encoding/json"

	"github.com/marlaone/shepard"
)

// Response represents an HTTP response (non-generic, concrete struct)
type Response struct {
	StatusCode StatusCode
	Headers    *Headers
	Body       *BodyWriter
}

// NewResponse creates a new Response with the given status code
func NewResponse(status StatusCode) Response {
	return Response{
		StatusCode: status,
		Headers:    NewHeaders(),
		Body:       NewBodyWriter(),
	}
}

// WithStatus returns a copy of the response with a new status code
func (r Response) WithStatus(status StatusCode) Response {
	r.StatusCode = status
	return r
}

// WithHeader returns a copy of the response with an added header
func (r Response) WithHeader(key, value string) Response {
	r.Headers.Set(key, value)
	return r
}

// WithBody returns a copy of the response with the given body
func (r Response) WithBody(data []byte) Response {
	r.Body.Write(data)
	return r
}

// IsSuccess returns true if the status code is 2xx
func (r Response) IsSuccess() bool {
	return r.StatusCode.IsSuccess()
}

// IsError returns true if the status code is 4xx or 5xx
func (r Response) IsError() bool {
	return r.StatusCode.IsError()
}

// IsRedirect returns true if the status code is 3xx
func (r Response) IsRedirect() bool {
	return r.StatusCode.IsRedirect()
}

// ResponseBuilder provides a fluent API for building responses
type ResponseBuilder struct {
	status  StatusCode
	headers *Headers
	body    *BodyWriter
	err     error
}

// NewResponseBuilder creates a new ResponseBuilder
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		status:  StatusOK,
		headers: NewHeaders(),
		body:    NewBodyWriter(),
		err:     nil,
	}
}

// Status sets the HTTP status code
func (b *ResponseBuilder) Status(status StatusCode) *ResponseBuilder {
	b.status = status
	return b
}

// Header adds a header
func (b *ResponseBuilder) Header(key, value string) *ResponseBuilder {
	b.headers.Set(key, value)
	return b
}

// Body sets the response body from bytes
func (b *ResponseBuilder) Body(data []byte) *ResponseBuilder {
	if b.err != nil {
		return b
	}
	result := b.body.Write(data)
	if result.IsErr() {
		b.err = result.UnwrapErr()
	}
	return b
}

// Json sets the response body as JSON and sets Content-Type header
func (b *ResponseBuilder) Json(value any) *ResponseBuilder {
	if b.err != nil {
		return b
	}
	data, err := json.Marshal(value)
	if err != nil {
		b.err = err
		return b
	}
	b.headers.Set(HeaderContentType, "application/json")
	b.body.Write(data)
	return b
}

// Text sets the response body as plain text and sets Content-Type header
func (b *ResponseBuilder) Text(text string) *ResponseBuilder {
	if b.err != nil {
		return b
	}
	b.headers.Set(HeaderContentType, "text/plain; charset=utf-8")
	b.body.WriteString(text)
	return b
}

// Html sets the response body as HTML and sets Content-Type header
func (b *ResponseBuilder) Html(html string) *ResponseBuilder {
	if b.err != nil {
		return b
	}
	b.headers.Set(HeaderContentType, "text/html; charset=utf-8")
	b.body.WriteString(html)
	return b
}

// File sets the response body for file download with default content type
func (b *ResponseBuilder) File(data []byte, filename string) *ResponseBuilder {
	return b.FileWithContentType(data, filename, "application/octet-stream")
}

// FileWithContentType sets the response body for file download with custom content type
func (b *ResponseBuilder) FileWithContentType(data []byte, filename, contentType string) *ResponseBuilder {
	if b.err != nil {
		return b
	}
	b.headers.Set(HeaderContentType, contentType)
	b.headers.Set("Content-Disposition", `attachment; filename="`+escapeFilename(filename)+`"`)
	b.body.Write(data)
	return b
}

// Build creates the Response
func (b *ResponseBuilder) Build() shepard.Result[Response, error] {
	if b.err != nil {
		return shepard.Err[Response, error](b.err)
	}
	return shepard.Ok[Response, error](Response{
		StatusCode: b.status,
		Headers:    b.headers,
		Body:       b.body,
	})
}

// --- Helper functions for quick response creation ---

// Text creates a 200 OK text response
func Text(text string) Response {
	return NewResponseBuilder().
		Status(StatusOK).
		Text(text).
		Build().
		UnwrapOrDefault()
}

// TextWithStatus creates a text response with a custom status
func TextWithStatus(status StatusCode, text string) Response {
	return NewResponseBuilder().
		Status(status).
		Text(text).
		Build().
		UnwrapOrDefault()
}

// Json creates a 200 OK JSON response
func Json(value any) shepard.Result[Response, error] {
	return NewResponseBuilder().
		Status(StatusOK).
		Json(value).
		Build()
}

// JsonWithStatus creates a JSON response with a custom status
func JsonWithStatus(status StatusCode, value any) shepard.Result[Response, error] {
	return NewResponseBuilder().
		Status(status).
		Json(value).
		Build()
}

// Html creates a 200 OK HTML response
func Html(html string) Response {
	return NewResponseBuilder().
		Status(StatusOK).
		Html(html).
		Build().
		UnwrapOrDefault()
}

// HtmlWithStatus creates an HTML response with a custom status
func HtmlWithStatus(status StatusCode, html string) Response {
	return NewResponseBuilder().
		Status(status).
		Html(html).
		Build().
		UnwrapOrDefault()
}

// Redirect creates a 302 Found redirect response
func Redirect(location string) Response {
	return NewResponse(StatusFound).
		WithHeader(HeaderLocation, location)
}

// PermanentRedirect creates a 308 Permanent Redirect response
func PermanentRedirect(location string) Response {
	return NewResponse(StatusPermanentRedirect).
		WithHeader(HeaderLocation, location)
}

// TemporaryRedirect creates a 307 Temporary Redirect response
func TemporaryRedirect(location string) Response {
	return NewResponse(StatusTemporaryRedirect).
		WithHeader(HeaderLocation, location)
}

// NoContent creates a 204 No Content response
func NoContent() Response {
	return NewResponse(StatusNoContent)
}

// Created creates a 201 Created response with optional location
func Created(location string) Response {
	resp := NewResponse(StatusCreated)
	if location != "" {
		resp = resp.WithHeader(HeaderLocation, location)
	}
	return resp
}

// BadRequestResponse creates a 400 Bad Request response
func BadRequestResponse(message string) Response {
	return TextWithStatus(StatusBadRequest, message)
}

// NotFoundResponse creates a 404 Not Found response
func NotFoundResponse(message string) Response {
	return TextWithStatus(StatusNotFound, message)
}

// InternalErrorResponse creates a 500 Internal Server Error response
func InternalErrorResponse(message string) Response {
	return TextWithStatus(StatusInternalServerError, message)
}
