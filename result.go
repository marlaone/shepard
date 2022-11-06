package shepard

import (
	"fmt"
	"reflect"
)

type ResultUnwrapElseFunc[T any, E any] func(err E) T
type ResultAndThenFunc[T any, E any] func(val T) Result[T, E]
type ResultOrElseFunc[T any, E any] func(val E) Result[T, E]

// Result is a type that represents either success (Ok) or failure (Err)
type Result[T any, E any] struct {
	ok  *T
	err *Error[E]
}

// Ok (T, E) contains the success value
func Ok[T any, E any](val T) Result[T, E] {
	return Result[T, E]{
		ok:  &val,
		err: nil,
	}
}

// Err (E) contains the error value
func Err[T any, E any](err E) Result[T, E] {
	return Result[T, E]{
		ok:  nil,
		err: NewError[E](err),
	}
}

// Unwrap returns the contained Ok value, consuming the self value.
//
// Because this function may panic, its use is generally discouraged. Instead, prefer to use pattern matching and handle the Err case explicitly, or call UnwrapOr, UnwrapOrElse, or UnwrapOrDefault.
//
// Panics if the value is an Err, with a panic message provided by the Err’s value.
func (r Result[T, E]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Errorf("%v", r.err.Value()))
	}
	return *r.ok
}

// UnwrapOr returns the contained Ok value or a provided default.
//
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (r Result[T, E]) UnwrapOr(defaultValue T) T {
	if r.IsErr() {
		return defaultValue
	}
	return r.Unwrap()
}

// UnwrapOrElse returns the contained Ok value or computes it from a closure.
func (r Result[T, E]) UnwrapOrElse(elseFunc ResultUnwrapElseFunc[T, E]) T {
	if r.IsErr() {
		return elseFunc(r.Err().Unwrap())
	}
	return r.Unwrap()
}

// UnwrapOrDefault returns the contained Ok value or a default
//
// Consumes the self argument then, if Ok, returns the contained value, otherwise if Err, returns the default value for that type.
func (r Result[T, E]) UnwrapOrDefault() T {
	if r.IsErr() {
		return GetDefault[T]()
	}
	return r.Unwrap()
}

// IsOk returns true if the Result is Ok
func (r Result[T, E]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result is Err
func (r Result[T, E]) IsErr() bool {
	return r.err != nil
}

// Or returns res if the Result is Err, otherwise returns the Ok value of self.
//
// Arguments passed to or are eagerly evaluated; if you are passing the Result of a function call, it is recommended to use or_else, which is lazily evaluated.
func (r Result[T, E]) Or(res Result[T, E]) Result[T, E] {
	if !r.IsOk() {
		return res
	}
	return r
}

// OrElse calls op if the Result is Err, otherwise returns the Ok value of self.
//
// This function can be used for control flow based on Result values.
func (r Result[T, E]) OrElse(op ResultOrElseFunc[T, E]) Result[T, E] {
	if r.IsOk() {
		return r
	}
	return op(r.Err().Unwrap())
}

// And returns res if the Result is Ok, otherwise returns the Err value of self.
func (r Result[T, E]) And(res Result[T, E]) Result[T, E] {
	if r.IsOk() && res.IsOk() {
		return res
	}
	if r.IsErr() {
		return r
	}
	return res
}

// AndThen calls op if the Result is Ok, otherwise returns the Err value of self.
//
// This function can be used for control flow based on Result values.
func (r Result[T, E]) AndThen(op ResultAndThenFunc[T, E]) Result[T, E] {
	if r.IsOk() {
		return op(*r.ok)
	}
	return r
}

// Equal checks if two Result`s contain the same Result value
func (r Result[T, E]) Equal(res Result[T, E]) bool {
	if r.IsOk() && res.IsOk() {
		return reflect.DeepEqual(r.Unwrap(), res.Unwrap())
	}
	if r.IsErr() && res.IsErr() {
		return reflect.DeepEqual(r.Err(), res.Err())
	}
	return false
}

// Ok converts from Result[T, E] to Option[T].
// Converts self into an Option[T], consuming self, and discarding the error, if any
func (r Result[T, E]) Ok() Option[T] {
	if r.IsErr() {
		return None[T]()
	}
	return Some[T](r.Unwrap())
}

// Err converts from Result[T, E] to Option[E].
// Converts self into an Option[E], consuming self, and discarding the success value, if any.
func (r Result[T, E]) Err() Option[E] {
	if r.IsErr() {
		err := *r.err
		return Some[E](err.Value())
	}
	return None[E]()
}

// Expect returns the contained Ok value, consuming the self value.
//
// Because this function may panic, its use is generally discouraged. Instead, prefer to use pattern matching and handle the Err case explicitly, or call UnwrapOr, UnwrapOrElse, or UnwrapOrDefault.
//
// Panics if the value is an Err, with a panic message including the passed message, and the content of the Err.
func (r Result[T, E]) Expect(err E) T {
	if r.IsErr() {
		panic(fmt.Errorf("%v: %v", err, r.err.Value()))
	}
	return r.Unwrap()
}

// ExpectErr returns the contained Err value, consuming the self value.
//
// Panics if the value is an Ok, with a panic message including the passed message, and the content of the Ok.
func (r Result[T, E]) ExpectErr(err E) E {
	if r.IsOk() {
		panic(fmt.Errorf("%v: %v", err, r.Unwrap()))
	}
	return r.Err().Unwrap()
}

// UnwrapErr Returns the contained Err value, consuming the self value.
//
// Panics if the value is an Ok, with a custom panic message provided by the Ok’s value.
func (r Result[T, E]) UnwrapErr() E {
	if r.IsOk() {
		panic(fmt.Errorf("%v", r.Unwrap()))
	}
	return r.Err().Unwrap()
}

func (r Result[T, E]) AsMut() (*T, *Error[E]) {
	return r.ok, r.err
}
