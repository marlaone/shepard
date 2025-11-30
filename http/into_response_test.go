package http

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/stretchr/testify/assert"
)

func TestIntoResponse_Response(t *testing.T) {
	// Response itself implements IntoResponse (identity)
	resp := NewResponse(StatusOK).WithBody([]byte("hello"))

	result := resp.IntoResponse()

	assert.Equal(t, StatusOK, result.StatusCode)
	assert.Equal(t, "hello", string(result.Body.Bytes()))
}

func TestIntoResponse_String(t *testing.T) {
	s := "Hello, World!"

	resp := StringResponse(s).IntoResponse()

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, World!", string(resp.Body.Bytes()))

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "text/plain; charset=utf-8", *ct.Unwrap())
}

func TestIntoResponse_Bytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}

	resp := BytesResponse(data).IntoResponse()

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, data, resp.Body.Bytes())

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/octet-stream", *ct.Unwrap())
}

func TestIntoResponse_StatusString(t *testing.T) {
	// Tuple-like: (StatusCode, string)
	resp := StatusTextResponse(StatusCreated, "Resource created").IntoResponse()

	assert.Equal(t, StatusCreated, resp.StatusCode)
	assert.Equal(t, "Resource created", string(resp.Body.Bytes()))
}

func TestIntoResponse_Result_Ok(t *testing.T) {
	// Result[Response, HttpError] - Ok case
	result := shepard.Ok[Response, HttpError](NewResponse(StatusOK).WithBody([]byte("success")))

	resp := ResultResponse(result)

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "success", string(resp.Body.Bytes()))
}

func TestIntoResponse_Result_Err(t *testing.T) {
	// Result[Response, HttpError] - Err case
	result := shepard.Err[Response, HttpError](NotFound("resource not found"))

	resp := ResultResponse(result)

	assert.Equal(t, StatusNotFound, resp.StatusCode)
	assert.Contains(t, string(resp.Body.Bytes()), "resource not found")
}

func TestIntoResponse_HttpError(t *testing.T) {
	err := BadRequest("invalid input")

	resp := err.IntoResponse()

	assert.Equal(t, StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resp.Body.Bytes()), "invalid input")
}

func TestJsonResponse(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	jr := JsonBody(User{Name: "John"})
	resp := jr.IntoResponse()

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, `{"name":"John"}`, string(resp.Body.Bytes()))

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/json", *ct.Unwrap())
}

func TestJsonResponseWithStatus(t *testing.T) {
	jr := JsonBodyWithStatus(StatusCreated, map[string]int{"id": 1})
	resp := jr.IntoResponse()

	assert.Equal(t, StatusCreated, resp.StatusCode)
}

func TestHtmlResponse(t *testing.T) {
	hr := HtmlBody("<h1>Hello</h1>")
	resp := hr.IntoResponse()

	assert.Equal(t, StatusOK, resp.StatusCode)

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "text/html; charset=utf-8", *ct.Unwrap())
}

// FileResponse tests

func TestFileResponse_Attachment(t *testing.T) {
	data := []byte("file content")

	fr := File(data, "document.txt")
	resp := fr.IntoResponse()

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, data, resp.Body.Bytes())

	// Check Content-Type (defaults to application/octet-stream)
	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/octet-stream", *ct.Unwrap())

	// Check Content-Disposition (attachment)
	cd := resp.Headers.Get("Content-Disposition")
	assert.True(t, cd.IsSome())
	assert.Equal(t, `attachment; filename="document.txt"`, *cd.Unwrap())
}

func TestFileResponse_WithContentType(t *testing.T) {
	data := []byte("PDF content")

	fr := FileWithContentType(data, "report.pdf", "application/pdf")
	resp := fr.IntoResponse()

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/pdf", *ct.Unwrap())

	cd := resp.Headers.Get("Content-Disposition")
	assert.True(t, cd.IsSome())
	assert.Equal(t, `attachment; filename="report.pdf"`, *cd.Unwrap())
}

func TestFileResponse_Inline(t *testing.T) {
	data := []byte("<svg></svg>")

	fr := InlineFile(data, "image.svg", "image/svg+xml")
	resp := fr.IntoResponse()

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "image/svg+xml", *ct.Unwrap())

	// Inline files use inline disposition
	cd := resp.Headers.Get("Content-Disposition")
	assert.True(t, cd.IsSome())
	assert.Equal(t, `inline; filename="image.svg"`, *cd.Unwrap())
}

func TestFileResponse_FilenameEscaping(t *testing.T) {
	data := []byte("content")

	// Filenames with special characters should be properly escaped
	fr := File(data, `test"file.txt`)
	resp := fr.IntoResponse()

	cd := resp.Headers.Get("Content-Disposition")
	assert.True(t, cd.IsSome())
	// Quote should be escaped
	assert.Equal(t, `attachment; filename="test\"file.txt"`, *cd.Unwrap())
}
