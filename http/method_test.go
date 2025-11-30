package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethod_String(t *testing.T) {
	tests := []struct {
		method   Method
		expected string
	}{
		{MethodGet, "GET"},
		{MethodPost, "POST"},
		{MethodPut, "PUT"},
		{MethodDelete, "DELETE"},
		{MethodPatch, "PATCH"},
		{MethodHead, "HEAD"},
		{MethodOptions, "OPTIONS"},
		{MethodTrace, "TRACE"},
		{MethodConnect, "CONNECT"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method.String())
		})
	}
}

func TestMethodFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected Method
		isOk     bool
	}{
		{"GET", MethodGet, true},
		{"POST", MethodPost, true},
		{"PUT", MethodPut, true},
		{"DELETE", MethodDelete, true},
		{"PATCH", MethodPatch, true},
		{"HEAD", MethodHead, true},
		{"OPTIONS", MethodOptions, true},
		{"TRACE", MethodTrace, true},
		{"CONNECT", MethodConnect, true},
		{"get", MethodGet, true},      // case insensitive
		{"post", MethodPost, true},    // case insensitive
		{"INVALID", Method(0), false}, // invalid method
		{"", Method(0), false},        // empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MethodFromString(tt.input)
			if tt.isOk {
				assert.True(t, result.IsOk())
				assert.Equal(t, tt.expected, result.Unwrap())
			} else {
				assert.True(t, result.IsErr())
			}
		})
	}
}

func TestMethod_IsIdempotent(t *testing.T) {
	idempotent := []Method{MethodGet, MethodHead, MethodPut, MethodDelete, MethodOptions, MethodTrace}
	notIdempotent := []Method{MethodPost, MethodPatch, MethodConnect}

	for _, m := range idempotent {
		t.Run(m.String()+"_idempotent", func(t *testing.T) {
			assert.True(t, m.IsIdempotent())
		})
	}

	for _, m := range notIdempotent {
		t.Run(m.String()+"_not_idempotent", func(t *testing.T) {
			assert.False(t, m.IsIdempotent())
		})
	}
}

func TestMethod_IsSafe(t *testing.T) {
	safe := []Method{MethodGet, MethodHead, MethodOptions, MethodTrace}
	notSafe := []Method{MethodPost, MethodPut, MethodDelete, MethodPatch, MethodConnect}

	for _, m := range safe {
		t.Run(m.String()+"_safe", func(t *testing.T) {
			assert.True(t, m.IsSafe())
		})
	}

	for _, m := range notSafe {
		t.Run(m.String()+"_not_safe", func(t *testing.T) {
			assert.False(t, m.IsSafe())
		})
	}
}
