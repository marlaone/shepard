package http

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {

	assert := assert.New(t)

	testCases := []struct {
		name     string
		query    string
		expected shepard.Result[Values, error]
	}{
		{
			name:  "empty query",
			query: "",
			expected: shepard.Ok[Values, error](Values{
				values: hashmap.New[string, slice.Slice[string]](),
			}),
		},
		{
			name:  "single key-value pair",
			query: "key=value",
			expected: shepard.Ok[Values, error](Values{
				values: hashmap.From[string, slice.Slice[string]]([]hashmap.Pair[string, slice.Slice[string]]{
					{Key: "key", Value: slice.Init("value")},
				}),
			}),
		},
		{
			name:  "multiple key-value pairs",
			query: "key1=value1&key2=value2",
			expected: shepard.Ok[Values, error](Values{
				values: hashmap.From[string, slice.Slice[string]]([]hashmap.Pair[string, slice.Slice[string]]{
					{Key: "key1", Value: slice.Init("value1")},
					{Key: "key2", Value: slice.Init("value2")},
				}),
			}),
		},
		{
			name:  "multiple values for a single key",
			query: "key=value1&key=value2",
			expected: shepard.Ok[Values, error](Values{
				values: hashmap.From[string, slice.Slice[string]]([]hashmap.Pair[string, slice.Slice[string]]{
					{Key: "key", Value: slice.Init("value1", "value2")},
				}),
			}),
		},
		{
			name:  "multiple keys with multiple values",
			query: "key1=value1&key2=value2&key1=value3&key2=value4",
			expected: shepard.Ok[Values, error](Values{
				values: hashmap.From[string, slice.Slice[string]]([]hashmap.Pair[string, slice.Slice[string]]{
					{Key: "key1", Value: slice.Init("value1", "value3")},
					{Key: "key2", Value: slice.Init("value2", "value4")},
				}),
			}),
		},
		{
			name:     "empty key",
			query:    "=value",
			expected: shepard.Ok[Values, error](*NewValues()),
		},
		{
			name:  "empty value",
			query: "key=",
			expected: shepard.Ok[Values, error](Values{
				values: hashmap.From[string, slice.Slice[string]]([]hashmap.Pair[string, slice.Slice[string]]{
					{Key: "key", Value: slice.Init[string]("")},
				}),
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ParseQuery(tc.query)

			assert.Equal(tc.expected, actual)
		})
	}
}

func TestParseURL(t *testing.T) {

	assert := assert.New(t)

	testCases := []struct {
		name     string
		url      string
		expected shepard.Result[URL, error]
	}{
		{
			name: "empty url",
			url:  "",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "",
				Host:   "",
				Port:   0,
				Path:   "",
			}),
		},
		{
			name: "scheme",
			url:  "http://",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "http",
				Host:   "",
				Port:   80,
				Path:   "",
			}),
		},
		{
			name: "host",
			url:  "http://example.com",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "http",
				Host:   "example.com",
				Port:   80,
				Path:   "",
			}),
		},
		{
			name: "port",
			url:  "http://example.com:8080",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "http",
				Host:   "example.com",
				Port:   8080,
				Path:   "",
			}),
		},
		{
			name: "path",
			url:  "http://example.com:8080/path",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "http",
				Host:   "example.com",
				Port:   8080,
				Path:   "/path",
			}),
		},
		{
			name: "query",
			url:  "http://example.com:8080/path?key=value",
			expected: shepard.Ok[URL, error](URL{
				Scheme:   "http",
				Host:     "example.com",
				Port:     8080,
				Path:     "/path",
				RawQuery: "key=value",
			}),
		},
		{
			name: "fragment",
			url:  "http://example.com:8080/path?key=value#fragment",
			expected: shepard.Ok[URL, error](URL{
				Scheme:   "http",
				Host:     "example.com",
				Port:     8080,
				Path:     "/path",
				Fragment: "fragment",
				RawQuery: "key=value",
			}),
		},
		{
			name: "user",
			url:  "http://user:pass@example",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "http",
				Host:   "example",
				User: &Userinfo{
					Username: "user",
					Password: "pass",
				},
				Port: 80,
				Path: "",
			}),
		},
		{
			name: "user with port",
			url:  "http://user:pass@example:8080",
			expected: shepard.Ok[URL, error](URL{
				Scheme: "http",
				Host:   "example",
				User: &Userinfo{
					Username: "user",
					Password: "pass",
				},
				Port: 8080,
				Path: "",
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ParseURL(tc.url)

			assert.Equal(tc.expected, actual)
		})
	}
}
