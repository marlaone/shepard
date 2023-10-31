package http

import (
	"io"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
)

type RequestBody io.Reader

type RequestHead struct {
	Method  Method
	URI     URI
	Version Version
	Headers hashmap.HashMap[string, []string]
}

func (r RequestHead) Default() RequestHead {
	return RequestHead{
		Method:  r.Method.Default(),
		URI:     r.URI.Default(),
		Version: r.Version.Default(),
		Headers: hashmap.WithCapacity[string, []string](16),
	}
}

type Request[T RequestBody] struct {
	head RequestHead
	body T
}

func (r Request[T]) Default() Request[T] {
	var body T
	return Request[T]{
		head: RequestHead{}.Default(),
		body: body,
	}
}

type RequestBuilder[T RequestBody] struct {
	request Request[T]
}

func NewRequestBuilder[T RequestBody]() *RequestBuilder[T] {
	return &RequestBuilder[T]{
		request: Request[T]{}.Default(),
	}
}

func (r *RequestBuilder[T]) Method(method Method) *RequestBuilder[T] {
	r.request.head.Method = method
	return r
}

func (r *RequestBuilder[T]) URI(uri URI) *RequestBuilder[T] {
	r.request.head.URI = uri
	return r
}

func (r *RequestBuilder[T]) Version(version Version) *RequestBuilder[T] {
	r.request.head.Version = version
	return r
}

func (r *RequestBuilder[T]) Header(key string, values ...string) *RequestBuilder[T] {
	header := r.request.head.Headers

	if header.ContainsKey(key) {
		header.Entry(key).AndModify(func(s *[]string) {
			*s = append(*s, values...)
		})
	} else {
		header.Insert(key, values)
	}

	return r
}

func (r *RequestBuilder[T]) Body(body T) shepard.Result[Request[T], error] {
	r.request.body = body
	return shepard.Ok[Request[T], error](r.request)
}
