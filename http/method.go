package http

import (
	"strings"

	"github.com/marlaone/shepard"
)

// Method represents an HTTP method
type Method uint8

const (
	MethodGet Method = iota + 1
	MethodPost
	MethodPut
	MethodDelete
	MethodPatch
	MethodHead
	MethodOptions
	MethodTrace
	MethodConnect
)

var methodStrings = map[Method]string{
	MethodGet:     "GET",
	MethodPost:    "POST",
	MethodPut:     "PUT",
	MethodDelete:  "DELETE",
	MethodPatch:   "PATCH",
	MethodHead:    "HEAD",
	MethodOptions: "OPTIONS",
	MethodTrace:   "TRACE",
	MethodConnect: "CONNECT",
}

var stringToMethod = map[string]Method{
	"GET":     MethodGet,
	"POST":    MethodPost,
	"PUT":     MethodPut,
	"DELETE":  MethodDelete,
	"PATCH":   MethodPatch,
	"HEAD":    MethodHead,
	"OPTIONS": MethodOptions,
	"TRACE":   MethodTrace,
	"CONNECT": MethodConnect,
}

// String returns the string representation of the method
func (m Method) String() string {
	if s, ok := methodStrings[m]; ok {
		return s
	}
	return ""
}

// MethodFromString parses a string into a Method
// Returns Result[Method, error] - Ok if valid, Err if invalid
func MethodFromString(s string) shepard.Result[Method, error] {
	upper := strings.ToUpper(s)
	if m, ok := stringToMethod[upper]; ok {
		return shepard.Ok[Method, error](m)
	}
	return shepard.Err[Method, error](ErrInvalidMethod)
}

// IsIdempotent returns true if the method is idempotent
// Idempotent methods: GET, HEAD, PUT, DELETE, OPTIONS, TRACE
func (m Method) IsIdempotent() bool {
	switch m {
	case MethodGet, MethodHead, MethodPut, MethodDelete, MethodOptions, MethodTrace:
		return true
	default:
		return false
	}
}

// IsSafe returns true if the method is safe (read-only)
// Safe methods: GET, HEAD, OPTIONS, TRACE
func (m Method) IsSafe() bool {
	switch m {
	case MethodGet, MethodHead, MethodOptions, MethodTrace:
		return true
	default:
		return false
	}
}
