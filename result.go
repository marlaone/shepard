package shepard

import (
	"fmt"
)

type ResultUnwrapElseFunc[T comparable] func(err error) T
type ResultAndThenFunc[T comparable] func(val *T) *Result[T]

type Result[T comparable] struct {
	ok  *T
	err error
}

func Ok[T comparable](val T) *Result[T] {
	return &Result[T]{
		ok:  &val,
		err: nil,
	}
}

func Err[T comparable](err error) *Result[T] {
	return &Result[T]{
		ok:  nil,
		err: err,
	}
}

func (r *Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Errorf("can't unwrap errored result: %v", r.Error()))
	}
	return *r.ok
}

func (r *Result[T]) UnwrapOr(defaultValue T) T {
	if r.IsErr() {
		return defaultValue
	}
	return r.Unwrap()
}

func (r *Result[T]) UnwrapOrElse(elseFunc ResultUnwrapElseFunc[T]) T {
	if r.IsErr() {
		return elseFunc(r.Error())
	}
	return r.Unwrap()
}

func (o *Result[T]) UnwrapOrDefault() T {
	if o.IsErr() {
		var zero T
		return zero
	}
	return o.Unwrap()
}

func (r *Result[T]) Error() error {
	return r.err
}

func (r *Result[T]) IsOk() bool {
	return r.err == nil
}

func (r *Result[T]) IsErr() bool {
	return r.err != nil
}

func (r *Result[T]) Or(orRes *Result[T]) *Result[T] {
	if !r.IsOk() {
		return orRes
	}
	return r
}

func (r *Result[T]) And(res *Result[T]) *Result[T] {
	if r.IsOk() && res.IsOk() {
		return res
	}
	if res.IsErr() {
		return res
	}
	return r
}

func (r *Result[T]) AndThen(op ResultAndThenFunc[T]) *Result[T] {
	if r.IsOk() {
		return op(r.ok)
	}
	return r
}

func (r *Result[T]) Equal(res *Result[T]) bool {
	if r.IsOk() && res.IsOk() {
		return r.Unwrap() == res.Unwrap()
	}
	if r.IsErr() && res.IsErr() {
		return r.Error().Error() == res.Error().Error()
	}
	return false
}
