package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse_New(t *testing.T) {
	resp := NewResponse(StatusOK)

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Headers)
	assert.NotNil(t, resp.Body)
}

func TestResponse_WithStatus(t *testing.T) {
	resp := NewResponse(StatusOK).WithStatus(StatusCreated)

	assert.Equal(t, StatusCreated, resp.StatusCode)
}

func TestResponse_WithHeader(t *testing.T) {
	resp := NewResponse(StatusOK).WithHeader("Content-Type", "application/json")

	result := resp.Headers.Get("Content-Type")
	assert.True(t, result.IsSome())
	assert.Equal(t, "application/json", result.Unwrap())
}

func TestResponse_WithBody(t *testing.T) {
	resp := NewResponse(StatusOK).WithBody([]byte("Hello, World!"))

	assert.Equal(t, "Hello, World!", string(resp.Body.Bytes()))
}

func TestResponseBuilder(t *testing.T) {
	result := NewResponseBuilder().
		Status(StatusCreated).
		Header("Content-Type", "application/json").
		Header("X-Request-Id", "123").
		Body([]byte(`{"id": 1}`)).
		Build()

	assert.True(t, result.IsOk())
	resp := result.Unwrap()

	assert.Equal(t, StatusCreated, resp.StatusCode)
	assert.True(t, resp.Headers.Contains("Content-Type"))
	assert.True(t, resp.Headers.Contains("X-Request-Id"))
	assert.Equal(t, `{"id": 1}`, string(resp.Body.Bytes()))
}

func TestResponseBuilder_Json(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	result := NewResponseBuilder().
		Status(StatusOK).
		Json(User{Name: "John", Age: 30}).
		Build()

	assert.True(t, result.IsOk())
	resp := result.Unwrap()

	assert.Equal(t, `{"name":"John","age":30}`, string(resp.Body.Bytes()))

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "application/json", ct.Unwrap())
}

func TestResponseBuilder_Text(t *testing.T) {
	result := NewResponseBuilder().
		Status(StatusOK).
		Text("Hello, World!").
		Build()

	assert.True(t, result.IsOk())
	resp := result.Unwrap()

	assert.Equal(t, "Hello, World!", string(resp.Body.Bytes()))

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "text/plain; charset=utf-8", ct.Unwrap())
}

func TestResponseBuilder_Html(t *testing.T) {
	result := NewResponseBuilder().
		Status(StatusOK).
		Html("<h1>Hello</h1>").
		Build()

	assert.True(t, result.IsOk())
	resp := result.Unwrap()

	ct := resp.Headers.Get("Content-Type")
	assert.True(t, ct.IsSome())
	assert.Equal(t, "text/html; charset=utf-8", ct.Unwrap())
}

// Helper function tests

func TestText(t *testing.T) {
	resp := Text("Hello!")

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello!", string(resp.Body.Bytes()))
}

func TestTextWithStatus(t *testing.T) {
	resp := TextWithStatus(StatusCreated, "Created!")

	assert.Equal(t, StatusCreated, resp.StatusCode)
	assert.Equal(t, "Created!", string(resp.Body.Bytes()))
}

func TestJson(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := Json(data)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, `{"key":"value"}`, string(resp.Body.Bytes()))
}

func TestJsonWithStatus(t *testing.T) {
	data := map[string]int{"id": 1}
	result := JsonWithStatus(StatusCreated, data)

	assert.True(t, result.IsOk())
	resp := result.Unwrap()
	assert.Equal(t, StatusCreated, resp.StatusCode)
}

func TestHtml(t *testing.T) {
	resp := Html("<h1>Title</h1>")

	assert.Equal(t, StatusOK, resp.StatusCode)
	assert.Equal(t, "<h1>Title</h1>", string(resp.Body.Bytes()))
}

func TestRedirect(t *testing.T) {
	resp := Redirect("/new-location")

	assert.Equal(t, StatusFound, resp.StatusCode)
	loc := resp.Headers.Get("Location")
	assert.True(t, loc.IsSome())
	assert.Equal(t, "/new-location", loc.Unwrap())
}

func TestPermanentRedirect(t *testing.T) {
	resp := PermanentRedirect("/permanent")

	assert.Equal(t, StatusPermanentRedirect, resp.StatusCode)
	loc := resp.Headers.Get("Location")
	assert.True(t, loc.IsSome())
	assert.Equal(t, "/permanent", loc.Unwrap())
}

func TestNoContent(t *testing.T) {
	resp := NoContent()

	assert.Equal(t, StatusNoContent, resp.StatusCode)
	assert.Equal(t, 0, resp.Body.Len())
}

func TestResponse_IsSuccess(t *testing.T) {
	okResp := NewResponse(StatusOK)
	errResp := NewResponse(StatusBadRequest)

	assert.True(t, okResp.IsSuccess())
	assert.False(t, errResp.IsSuccess())
}

func TestResponse_IsError(t *testing.T) {
	okResp := NewResponse(StatusOK)
	errResp := NewResponse(StatusInternalServerError)

	assert.False(t, okResp.IsError())
	assert.True(t, errResp.IsError())
}
