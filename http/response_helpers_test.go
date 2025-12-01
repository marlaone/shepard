package http

import (
	"testing"
)

func TestJsonOk(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := JsonOk(data)

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode.Int())
	}

	contentType := resp.Headers.Get(HeaderContentType)
	if contentType.IsNone() {
		t.Fatal("expected Content-Type header")
	}
	if contentType.Unwrap() != "application/json" {
		t.Errorf("expected application/json, got %s", contentType.Unwrap())
	}
}

func TestJsonCreated(t *testing.T) {
	data := map[string]int{"id": 123}
	result := JsonCreated(data)

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode.Int())
	}
}

func TestJsonAccepted(t *testing.T) {
	result := JsonAccepted(map[string]string{"status": "processing"})

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusAccepted {
		t.Errorf("expected status 202, got %d", resp.StatusCode.Int())
	}
}

func TestJsonWithStatusOk(t *testing.T) {
	result := JsonWithStatusOk(StatusPartialContent, []int{1, 2, 3})

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusPartialContent {
		t.Errorf("expected status 206, got %d", resp.StatusCode.Int())
	}
}

func TestTextOk(t *testing.T) {
	result := TextOk("Hello, World!")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode.Int())
	}

	body := string(resp.Body.Bytes())
	if body != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got '%s'", body)
	}
}

func TestTextWithStatusOk(t *testing.T) {
	result := TextWithStatusOk(StatusNotFound, "Not found")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode.Int())
	}
}

func TestHtmlOk(t *testing.T) {
	result := HtmlOk("<h1>Hello</h1>")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode.Int())
	}

	contentType := resp.Headers.Get(HeaderContentType)
	if contentType.IsNone() {
		t.Fatal("expected Content-Type header")
	}
	if contentType.Unwrap() != "text/html; charset=utf-8" {
		t.Errorf("expected text/html, got %s", contentType.Unwrap())
	}
}

func TestNoContentOk(t *testing.T) {
	result := NoContentOk()

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode.Int())
	}
}

func TestCreatedOk(t *testing.T) {
	result := CreatedOk("/users/123")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode.Int())
	}

	location := resp.Headers.Get(HeaderLocation)
	if location.IsNone() {
		t.Fatal("expected Location header")
	}
	if location.Unwrap() != "/users/123" {
		t.Errorf("expected /users/123, got %s", location.Unwrap())
	}
}

func TestRedirectOk(t *testing.T) {
	result := RedirectOk("https://example.com")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusFound {
		t.Errorf("expected status 302, got %d", resp.StatusCode.Int())
	}

	location := resp.Headers.Get(HeaderLocation)
	if location.IsNone() {
		t.Fatal("expected Location header")
	}
	if location.Unwrap() != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", location.Unwrap())
	}
}

func TestPermanentRedirectOk(t *testing.T) {
	result := PermanentRedirectOk("https://new-site.com")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusPermanentRedirect {
		t.Errorf("expected status 308, got %d", resp.StatusCode.Int())
	}
}

func TestStatusOk(t *testing.T) {
	result := StatusOk(StatusAccepted)

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusAccepted {
		t.Errorf("expected status 202, got %d", resp.StatusCode.Int())
	}
}

func TestResponseOk(t *testing.T) {
	resp := NewResponse(StatusOK)
	resp.Headers.Set("X-Custom", "value")

	result := ResponseOk(resp)

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	wrapped := result.Unwrap()
	if wrapped.StatusCode != StatusOK {
		t.Errorf("expected status 200, got %d", wrapped.StatusCode.Int())
	}
}

func TestFileOk(t *testing.T) {
	data := []byte("file content")
	result := FileOk(data, "test.txt")

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode.Int())
	}

	disposition := resp.Headers.Get("Content-Disposition")
	if disposition.IsNone() {
		t.Fatal("expected Content-Disposition header")
	}
}

func TestBytesOk(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	result := BytesOk(data)

	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	resp := result.Unwrap()
	if resp.StatusCode != StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode.Int())
	}

	body := resp.Body.Bytes()
	if len(body) != 3 {
		t.Errorf("expected 3 bytes, got %d", len(body))
	}
}
