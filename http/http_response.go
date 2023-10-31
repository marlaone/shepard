package http

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/iter"
)

// header parts

type Version string

func (v Version) String() string {
	return string(v)
}

func (v Version) Default() Version {
	return "1.1"
}

// response header

type ResponseHead struct {
	StatusCode StatusCode
	Version    Version
	Headers    hashmap.HashMap[string, []string]
}

func (r ResponseHead) Default() ResponseHead {
	return ResponseHead{
		StatusCode: r.StatusCode.Default(),
		Version:    r.Version.Default(),
		Headers:    hashmap.New[string, []string](),
	}
}

// response

type Body interface {
	Write() chan<- []byte
	Read() <-chan []byte
	Finish()
}

type Headers struct {
	headers hashmap.HashMap[string, slice.Slice[string]]
}

func (h Headers) Default() Headers {
	return Headers{
		headers: hashmap.WithCapacity[string, slice.Slice[string]](16),
	}
}

func (h *Headers) Has(key string) bool {
	return h.headers.ContainsKey(key)
}

func (h *Headers) Set(key string, values ...string) {
	h.headers.Entry(key).AndModify(func(s *slice.Slice[string]) {
		s.Clear()
		for _, v := range values {
			s.Push(v)
		}
	}).OrInsert(slice.Init(values...))
}

func (h *Headers) Get(key string) slice.Slice[string] {
	entry := h.headers.Entry(key)
	if entry.IsOccupied() {
		return *entry.Value()
	}
	return slice.New[string]()
}

func (h *Headers) Iter() iter.Iter[hashmap.Pair[*string, *slice.Slice[string]]] {
	return h.headers.Iter()
}

type Response[T Body] interface {
	SetStatusCode(statusCode StatusCode)
	StatusCode() StatusCode
	SetVersion(version Version)
	Version() Version
	SetHeader(key string, values ...string)
	Headers() *Headers
	SetBody(body T)
	Body() Body
}

// response builder

type ResponseBuilder[T Response[Body]] struct {
	response T
}

func NewResponseBuilder[T Response[Body]](res T) *ResponseBuilder[T] {
	return &ResponseBuilder[T]{
		response: res,
	}
}

func (r *ResponseBuilder[T]) Status(statusCode StatusCode) *ResponseBuilder[T] {
	r.response.SetStatusCode(statusCode)
	return r
}

func (r *ResponseBuilder[T]) Version(version Version) *ResponseBuilder[T] {
	r.response.SetVersion(version)
	return r
}

func (r *ResponseBuilder[T]) Header(key string, values ...string) *ResponseBuilder[T] {
	r.response.Headers().Set(key, values...)

	return r
}

func (r *ResponseBuilder[T]) Body(body Body) shepard.Result[T, error] {
	r.response.SetBody(body)
	return shepard.Ok[T, error](r.response)
}
