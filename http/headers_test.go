package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders_New(t *testing.T) {
	h := NewHeaders()
	assert.NotNil(t, h)
	assert.Equal(t, 0, h.Len())
}

func TestHeaders_Set_Get(t *testing.T) {
	h := NewHeaders()

	// Set a header
	h.Set("Content-Type", "application/json")

	// Get it back
	result := h.Get("Content-Type")
	assert.True(t, result.IsSome())
	assert.Equal(t, "application/json", result.Unwrap())
}

func TestHeaders_Get_CaseInsensitive(t *testing.T) {
	h := NewHeaders()
	h.Set("Content-Type", "application/json")

	// Should be case-insensitive
	result := h.Get("content-type")
	assert.True(t, result.IsSome())
	assert.Equal(t, "application/json", result.Unwrap())

	result = h.Get("CONTENT-TYPE")
	assert.True(t, result.IsSome())
	assert.Equal(t, "application/json", result.Unwrap())
}

func TestHeaders_Get_NotFound(t *testing.T) {
	h := NewHeaders()
	result := h.Get("X-Not-Found")
	assert.True(t, result.IsNone())
}

func TestHeaders_Append(t *testing.T) {
	h := NewHeaders()

	// Set initial value
	h.Set("Accept", "text/html")

	// Append another value
	h.Append("Accept", "application/json")

	// Get all values
	values := h.GetAll("Accept")
	assert.Equal(t, 2, values.Len())
	assert.Equal(t, "text/html", *values.Get(0).Unwrap())
	assert.Equal(t, "application/json", *values.Get(1).Unwrap())
}

func TestHeaders_Remove(t *testing.T) {
	h := NewHeaders()
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "text/html")

	assert.Equal(t, 2, h.Len())

	h.Remove("Content-Type")

	assert.Equal(t, 1, h.Len())
	assert.True(t, h.Get("Content-Type").IsNone())
	assert.True(t, h.Get("Accept").IsSome())
}

func TestHeaders_Contains(t *testing.T) {
	h := NewHeaders()
	h.Set("Content-Type", "application/json")

	assert.True(t, h.Contains("Content-Type"))
	assert.True(t, h.Contains("content-type"))
	assert.False(t, h.Contains("X-Not-Found"))
}

func TestHeaders_Clone(t *testing.T) {
	h := NewHeaders()
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "text/html")

	// Clone
	cloned := h.Clone()

	// Modify original
	h.Set("Content-Type", "text/plain")

	// Clone should be unchanged
	result := cloned.Get("Content-Type")
	assert.True(t, result.IsSome())
	assert.Equal(t, "application/json", result.Unwrap())
}

func TestHeaders_Iter(t *testing.T) {
	h := NewHeaders()
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "text/html")

	count := 0
	h.Iter().Foreach(func(_ int, _ HeaderEntry) {
		count++
	})

	assert.Equal(t, 2, count)
}

func TestHeaders_CommonHeaders(t *testing.T) {
	// Test common header constants exist
	assert.Equal(t, "Content-Type", HeaderContentType)
	assert.Equal(t, "Content-Length", HeaderContentLength)
	assert.Equal(t, "Accept", HeaderAccept)
	assert.Equal(t, "Authorization", HeaderAuthorization)
	assert.Equal(t, "Host", HeaderHost)
	assert.Equal(t, "User-Agent", HeaderUserAgent)
}
