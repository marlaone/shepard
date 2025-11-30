package http

import (
	"encoding/json"

	"github.com/marlaone/shepard"
)

// IntoResponse is a trait (interface) for types that can be converted into an HTTP Response
// This mirrors Rust's axum IntoResponse trait
type IntoResponse interface {
	IntoResponse() Response
}

// --- Response implements IntoResponse (identity) ---

func (r Response) IntoResponse() Response {
	return r
}

// --- HttpError implements IntoResponse ---

func (e HttpError) IntoResponse() Response {
	return TextWithStatus(e.StatusCode, e.Error())
}

// --- StringResponse wraps a string as IntoResponse ---

type StringResponse string

func (s StringResponse) IntoResponse() Response {
	return Text(string(s))
}

// --- BytesResponse wraps []byte as IntoResponse ---

type BytesResponse []byte

func (b BytesResponse) IntoResponse() Response {
	resp := NewResponse(StatusOK)
	resp.Headers.Set(HeaderContentType, "application/octet-stream")
	resp.Body.Write(b)
	return resp
}

// --- StatusText is a tuple-like (StatusCode, string) ---

type StatusText struct {
	Status StatusCode
	Text   string
}

func (st StatusText) IntoResponse() Response {
	return TextWithStatus(st.Status, st.Text)
}

// StatusTextResponse creates a StatusText from status and text
// This is a constructor function since Go doesn't have tuple syntax
func StatusTextResponse(status StatusCode, text string) StatusText {
	return StatusText{Status: status, Text: text}
}

// --- JsonBody wraps a value to be serialized as JSON ---

type JsonBodyValue struct {
	Status StatusCode
	Value  any
}

func JsonBody(value any) JsonBodyValue {
	return JsonBodyValue{Status: StatusOK, Value: value}
}

func JsonBodyWithStatus(status StatusCode, value any) JsonBodyValue {
	return JsonBodyValue{Status: status, Value: value}
}

func (j JsonBodyValue) IntoResponse() Response {
	resp := NewResponse(j.Status)
	resp.Headers.Set(HeaderContentType, "application/json")

	data, err := json.Marshal(j.Value)
	if err != nil {
		// If JSON marshaling fails, return an error response
		return InternalErrorResponse("failed to serialize JSON")
	}
	resp.Body.Write(data)
	return resp
}

// --- HtmlBody wraps HTML content ---

type HtmlBodyValue struct {
	Status StatusCode
	Html   string
}

func HtmlBody(html string) HtmlBodyValue {
	return HtmlBodyValue{Status: StatusOK, Html: html}
}

func HtmlBodyWithStatus(status StatusCode, html string) HtmlBodyValue {
	return HtmlBodyValue{Status: status, Html: html}
}

func (h HtmlBodyValue) IntoResponse() Response {
	return HtmlWithStatus(h.Status, h.Html)
}

// --- ResultResponse converts Result[Response, HttpError] to Response ---

// ResultResponse converts a Result into a Response
// If Ok, returns the Response; if Err, converts the error to a Response
func ResultResponse(result shepard.Result[Response, HttpError]) Response {
	if result.IsOk() {
		return result.Unwrap()
	}
	return result.UnwrapErr().IntoResponse()
}

// --- Optional Response wrappers for common patterns ---

// WithHeaders wraps a response with additional headers
type WithHeaders struct {
	Inner   IntoResponse
	Headers map[string]string
}

func (w WithHeaders) IntoResponse() Response {
	resp := w.Inner.IntoResponse()
	for k, v := range w.Headers {
		resp.Headers.Set(k, v)
	}
	return resp
}

// WithStatus wraps a response with a different status code
type WithStatus struct {
	Status StatusCode
	Inner  IntoResponse
}

func (w WithStatus) IntoResponse() Response {
	resp := w.Inner.IntoResponse()
	resp.StatusCode = w.Status
	return resp
}

// --- FileResponse for file downloads ---

// FileResponse represents a file download response
type FileResponse struct {
	Data        []byte
	Filename    string
	ContentType string
	Inline      bool // true = inline (display in browser), false = attachment (download)
}

// File creates a FileResponse for download (attachment disposition)
// Uses application/octet-stream as the default content type
func File(data []byte, filename string) FileResponse {
	return FileResponse{
		Data:        data,
		Filename:    filename,
		ContentType: "application/octet-stream",
		Inline:      false,
	}
}

// FileWithContentType creates a FileResponse with a custom content type
func FileWithContentType(data []byte, filename, contentType string) FileResponse {
	return FileResponse{
		Data:        data,
		Filename:    filename,
		ContentType: contentType,
		Inline:      false,
	}
}

// InlineFile creates a FileResponse for inline display (displayed in browser)
func InlineFile(data []byte, filename, contentType string) FileResponse {
	return FileResponse{
		Data:        data,
		Filename:    filename,
		ContentType: contentType,
		Inline:      true,
	}
}

// escapeFilename escapes special characters in a filename for Content-Disposition header
func escapeFilename(filename string) string {
	// Escape backslashes and quotes for the quoted-string format
	var result []byte
	for i := 0; i < len(filename); i++ {
		c := filename[i]
		if c == '"' || c == '\\' {
			result = append(result, '\\')
		}
		result = append(result, c)
	}
	return string(result)
}

func (f FileResponse) IntoResponse() Response {
	resp := NewResponse(StatusOK)
	resp.Headers.Set(HeaderContentType, f.ContentType)

	// Set Content-Disposition header
	disposition := "attachment"
	if f.Inline {
		disposition = "inline"
	}
	escapedFilename := escapeFilename(f.Filename)
	resp.Headers.Set("Content-Disposition", disposition+`; filename="`+escapedFilename+`"`)

	resp.Body.Write(f.Data)
	return resp
}
