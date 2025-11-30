package http

import (
	"github.com/marlaone/shepard"
)

// --- Result-returning response helpers ---
// These functions return shepard.Result[Response, HttpError] for use directly in handlers,
// eliminating the need for verbose shepard.Ok[Response, HttpError](...) wrapping.

// JsonOk returns a successful 200 JSON response
func JsonOk(value any) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](JsonBody(value).IntoResponse())
}

// JsonCreated returns a 201 Created JSON response
func JsonCreated(value any) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](JsonBodyWithStatus(StatusCreated, value).IntoResponse())
}

// JsonAccepted returns a 202 Accepted JSON response
func JsonAccepted(value any) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](JsonBodyWithStatus(StatusAccepted, value).IntoResponse())
}

// JsonWithStatusOk returns a JSON response with a custom status code
func JsonWithStatusOk(status StatusCode, value any) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](JsonBodyWithStatus(status, value).IntoResponse())
}

// TextOk returns a successful 200 text response
func TextOk(text string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](Text(text))
}

// TextWithStatusOk returns a text response with a custom status code
func TextWithStatusOk(status StatusCode, text string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](TextWithStatus(status, text))
}

// HtmlOk returns a successful 200 HTML response
func HtmlOk(html string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](Html(html))
}

// HtmlWithStatusOk returns an HTML response with a custom status code
func HtmlWithStatusOk(status StatusCode, html string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](HtmlWithStatus(status, html))
}

// NoContentOk returns a 204 No Content response
func NoContentOk() shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](NoContent())
}

// CreatedOk returns a 201 Created response with optional location header
func CreatedOk(location string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](Created(location))
}

// RedirectOk returns a 302 Found redirect response
func RedirectOk(url string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](Redirect(url))
}

// PermanentRedirectOk returns a 308 Permanent Redirect response
func PermanentRedirectOk(url string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](PermanentRedirect(url))
}

// TemporaryRedirectOk returns a 307 Temporary Redirect response
func TemporaryRedirectOk(url string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](TemporaryRedirect(url))
}

// StatusOk returns a response with just a status code
func StatusOk(status StatusCode) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](NewResponse(status))
}

// ResponseOk wraps any Response in a Result
func ResponseOk(resp Response) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](resp)
}

// FileOk returns a successful file download response
func FileOk(data []byte, filename string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](File(data, filename).IntoResponse())
}

// FileWithContentTypeOk returns a file download response with custom content type
func FileWithContentTypeOk(data []byte, filename, contentType string) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](FileWithContentType(data, filename, contentType).IntoResponse())
}

// BytesOk returns a successful response with raw bytes
func BytesOk(data []byte) shepard.Result[Response, HttpError] {
	return shepard.Ok[Response, HttpError](BytesResponse(data).IntoResponse())
}
