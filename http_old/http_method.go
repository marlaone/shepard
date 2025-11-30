package http

import (
	"fmt"
	"strings"

	"github.com/marlaone/shepard"
)

type Method string

func (m Method) String() string {
	return string(m)
}

func (m Method) Default() Method {
	return MethodGet
}

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH"
	MethodDelete  Method = "DELETE"
	MethodConnect Method = "CONNECT"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"
)

func TryMethodFromString(s string) shepard.Result[Method, error] {
	s = strings.TrimSpace(s)
	switch s {
	case MethodGet.String():
		return shepard.Ok[Method, error](MethodGet)
	case MethodHead.String():
		return shepard.Ok[Method, error](MethodHead)
	case MethodPost.String():
		return shepard.Ok[Method, error](MethodPost)
	case MethodPut.String():
		return shepard.Ok[Method, error](MethodPut)
	case MethodDelete.String():
		return shepard.Ok[Method, error](MethodDelete)
	case MethodConnect.String():
		return shepard.Ok[Method, error](MethodConnect)
	case MethodOptions.String():
		return shepard.Ok[Method, error](MethodOptions)
	case MethodTrace.String():
		return shepard.Ok[Method, error](MethodTrace)
	default:
		return shepard.Err[Method, error](fmt.Errorf("invalid method: %s", s))
	}
}
