package http

import (
	"bufio"
	"context"
	"io"

	"github.com/marlaone/shepard"
)

type RequestBody io.Reader

func ParseParams(body RequestBody) shepard.Result[Values, error] {
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()

		res := ParseQuery(line)
		if res.IsErr() {
			return shepard.Err[Values, error](res.Err().Unwrap())
		}
		return shepard.Ok[Values, error](res.Ok().Unwrap())
	}

	return shepard.Err[Values, error](scanner.Err())
}

type Request[T RequestBody] struct {
	Context context.Context
	Method  Method
	URL     URL
	Version Version
	Headers Headers
	params  *Values
	body    T
}

func (r Request[T]) Default() Request[T] {
	var body T
	return Request[T]{
		Context: context.Background(),
		Method:  r.Method.Default(),
		URL:     r.URL.Default(),
		Version: r.Version.Default(),
		Headers: r.Headers.Default(),
		body:    body,
	}
}

func (r *Request[T]) Params() *Values {
	return r.params
}

func (r *Request[T]) ParseForm() shepard.Result[any, error] {

	if r.params != nil {
		return shepard.Ok[any, error](nil)
	}

	params := ParseParams(r.body)
	if params.IsErr() {
		return shepard.Err[any, error](params.Err().Unwrap())
	}

	values := params.Unwrap()
	r.params = &values
	return shepard.Ok[any, error](nil)
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
	r.request.Method = method
	return r
}

func (r *RequestBuilder[T]) URL(url URL) *RequestBuilder[T] {
	r.request.URL = url
	return r
}

func (r *RequestBuilder[T]) Version(version Version) *RequestBuilder[T] {
	r.request.Version = version
	return r
}

func (r *RequestBuilder[T]) Header(key string, values ...string) *RequestBuilder[T] {
	r.request.Headers.Set(key, values...)

	return r
}

func (r *RequestBuilder[T]) Body(body T) shepard.Result[Request[T], error] {
	r.request.body = body
	return shepard.Ok[Request[T], error](r.request)
}
