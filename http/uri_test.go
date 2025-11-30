package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected URL
		isOk     bool
	}{
		{
			name:  "simple path",
			input: "/users",
			expected: URL{
				Path: "/users",
			},
			isOk: true,
		},
		{
			name:  "path with query",
			input: "/users?page=1&limit=10",
			expected: URL{
				Path:     "/users",
				RawQuery: "page=1&limit=10",
			},
			isOk: true,
		},
		{
			name:  "full URL",
			input: "https://example.com:8080/api/users?id=123",
			expected: URL{
				Scheme:   "https",
				Host:     "example.com",
				Port:     "8080",
				Path:     "/api/users",
				RawQuery: "id=123",
			},
			isOk: true,
		},
		{
			name:  "URL with fragment",
			input: "/page#section",
			expected: URL{
				Path:     "/page",
				Fragment: "section",
			},
			isOk: true,
		},
		{
			name:  "root path",
			input: "/",
			expected: URL{
				Path: "/",
			},
			isOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseURL(tt.input)
			if tt.isOk {
				assert.True(t, result.IsOk())
				url := result.Unwrap()
				assert.Equal(t, tt.expected.Scheme, url.Scheme)
				assert.Equal(t, tt.expected.Host, url.Host)
				assert.Equal(t, tt.expected.Port, url.Port)
				assert.Equal(t, tt.expected.Path, url.Path)
				assert.Equal(t, tt.expected.RawQuery, url.RawQuery)
				assert.Equal(t, tt.expected.Fragment, url.Fragment)
			} else {
				assert.True(t, result.IsErr())
			}
		})
	}
}

func TestURL_String(t *testing.T) {
	tests := []struct {
		name     string
		url      URL
		expected string
	}{
		{
			name:     "simple path",
			url:      URL{Path: "/users"},
			expected: "/users",
		},
		{
			name:     "path with query",
			url:      URL{Path: "/users", RawQuery: "page=1"},
			expected: "/users?page=1",
		},
		{
			name:     "full URL",
			url:      URL{Scheme: "https", Host: "example.com", Port: "8080", Path: "/api"},
			expected: "https://example.com:8080/api",
		},
		{
			name:     "with fragment",
			url:      URL{Path: "/page", Fragment: "top"},
			expected: "/page#top",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.url.String())
		})
	}
}

func TestQueryValues_Get(t *testing.T) {
	qv := NewQueryValues()
	qv.Set("page", "1")
	qv.Set("limit", "10")

	// Get existing key
	result := qv.Get("page")
	assert.True(t, result.IsSome())
	assert.Equal(t, "1", *result.Unwrap())

	// Get non-existing key
	result = qv.Get("unknown")
	assert.True(t, result.IsNone())
}

func TestQueryValues_GetAll(t *testing.T) {
	qv := NewQueryValues()
	qv.Add("tag", "go")
	qv.Add("tag", "rust")
	qv.Add("tag", "python")

	values := qv.GetAll("tag")
	assert.Equal(t, 3, values.Len())
}

func TestQueryValues_Parse(t *testing.T) {
	result := ParseQuery("page=1&limit=10&tag=go&tag=rust")
	assert.True(t, result.IsOk())

	qv := result.Unwrap()

	page := qv.Get("page")
	assert.True(t, page.IsSome())
	assert.Equal(t, "1", *page.Unwrap())

	limit := qv.Get("limit")
	assert.True(t, limit.IsSome())
	assert.Equal(t, "10", *limit.Unwrap())

	tags := qv.GetAll("tag")
	assert.Equal(t, 2, tags.Len())
}

func TestQueryValues_Encode(t *testing.T) {
	qv := NewQueryValues()
	qv.Set("page", "1")
	qv.Set("name", "John Doe")

	encoded := qv.Encode()
	// Should contain both params (order may vary)
	assert.Contains(t, encoded, "page=1")
	assert.Contains(t, encoded, "name=John+Doe")
}

func TestURL_Query(t *testing.T) {
	url := URL{
		Path:     "/search",
		RawQuery: "q=hello&page=1",
	}

	result := url.Query()
	assert.True(t, result.IsOk())

	qv := result.Unwrap()
	q := qv.Get("q")
	assert.True(t, q.IsSome())
	assert.Equal(t, "hello", *q.Unwrap())
}
