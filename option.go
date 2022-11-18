package shepard

import (
	"errors"
	"fmt"
	"reflect"
)

type OptionUnwrapElseFunc[T any] func() T
type OptionOkOrElseFunc[T any] func() T
type OptionAndThenFunc[T any] func(val T) Option[T]
type OptionFilterFunc[T any] func(val *T) bool
type OptionOrElseFunc[T any] func() Option[T]
type OptionGetOrInsertWithFunc[T any] func() T

type Option[T any] struct {
	v *T
}

func (o Option[T]) Default() Option[T] {
	return None[T]()
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		v: &val,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		v: nil,
	}
}

// Unwrap Returns the contained Some value, consuming the self value.
//
// Because this function may panic, its use is generally discouraged. Instead, prefer to use pattern matching and handle the None case explicitly, or call UnwrapOr, UnwrapOrElse, or UnwrapOrDefault.
//
// Panics if the self value equals None.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic(errors.New("unwrap on None Option"))
	}
	return *o.v
}

// IsSome returns true if the Option is a Some value.
func (o Option[T]) IsSome() bool {
	return o.v != nil
}

// IsNone returns true if the Option is a None value.
func (o Option[T]) IsNone() bool {
	return o.v == nil
}

// UnwrapOr returns the contained Some value or a provided default.
//
// Arguments passed to unwrap_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (o Option[T]) UnwrapOr(value T) T {
	if o.IsSome() {
		return o.Unwrap()
	}
	return value
}

// UnwrapOrElse returns the contained Some value or computes it from a closure.
func (o Option[T]) UnwrapOrElse(elseFunc OptionUnwrapElseFunc[T]) T {
	if o.IsSome() {
		return o.Unwrap()
	}
	return elseFunc()
}

// UnwrapOrDefault returns the contained Some value or a default.
//
// Consumes the self argument then, if Some, returns the contained value, otherwise if None, returns the default value for that type.
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		return GetDefault[T]()
	}
	return o.Unwrap()
}

// OkOr transforms the Option[T] into a Result[T, T], mapping Some(v) to Ok(v) and None to Err(err).
//
// Arguments passed to ok_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use OkOrElse, which is lazily evaluated.
func (o Option[T]) OkOr(err error) Result[T, error] {
	if o.IsSome() {
		return Ok[T, error](o.Unwrap())
	}
	return Err[T, error](err)
}

// OkOrElse transforms the Option[T] into a Result[T, T], mapping Some(v) to Ok(v) and None to Err(err()).
func (o Option[T]) OkOrElse(err OptionOkOrElseFunc[T]) Result[T, T] {
	if o.IsSome() {
		return Ok[T, T](o.Unwrap())
	}
	return Err[T, T](err())
}

// And returns None if the option is None, otherwise returns opt.
func (o Option[T]) And(opt Option[T]) Option[T] {
	if o.IsSome() && opt.IsSome() {
		return opt
	}
	return None[T]()
}

// AndThen returns None if the option is None, otherwise calls f with the wrapped value and returns the result.
//
// Some languages call this operation flatmap.
func (o Option[T]) AndThen(op OptionAndThenFunc[T]) Option[T] {
	if o.IsSome() {
		return op(*o.v)
	}
	return o
}

// Filter returns None if the option is None, otherwise calls predicate with the wrapped value and returns:
//
// Some(t) if predicate returns true (where t is the wrapped value), and
// None if predicate returns false.
func (o Option[T]) Filter(predicate OptionFilterFunc[T]) Option[T] {
	if o.IsSome() && predicate(o.v) {
		return o
	}
	return None[T]()
}

// Or returns the option if it contains a value, otherwise returns opt.
//
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use or_else, which is lazily evaluated.
func (o Option[T]) Or(opt Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	if opt.IsSome() {
		return opt
	}
	return None[T]()
}

// OrElse returns the option if it contains a value, otherwise calls f and returns the result.
func (o Option[T]) OrElse(f OptionOrElseFunc[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

// Xor returns Some if exactly one of self, opt is Some, otherwise returns None.
func (o Option[T]) Xor(opt Option[T]) Option[T] {
	if o.IsSome() && opt.IsSome() {
		return None[T]()
	}
	if o.IsSome() {
		return o
	}
	if opt.IsSome() {
		return opt
	}
	return None[T]()
}

// Insert inserts value into the option, then returns a mutable reference to it.
//
// If the option already contains a value, the old value is dropped.
//
// See also GetOrInsert, which doesn't update the value if the option already contains Some.
func (o *Option[T]) Insert(val T) *T {
	o.v = &val
	return o.v
}

// GetOrInsert inserts value into the option if it is None, then returns a mutable reference to the contained value.
//
// See also Insert, which updates the value even if the option already contains Some.
func (o *Option[T]) GetOrInsert(val T) *T {
	if o.IsSome() {
		return o.v
	}
	o.v = &val
	return o.v
}

// GetOrInsertDefault inserts the Default value into the option if it is None, then returns a mutable reference to the contained value.
func (o *Option[T]) GetOrInsertDefault() *T {
	if o.IsSome() {
		return o.v
	}

	defaultValue := GetDefault[T]()
	o.v = &defaultValue

	return o.v
}

// GetOrInsertWith inserts a value computed from f into the Option if it is None, then returns a mutable reference to the contained value.
func (o *Option[T]) GetOrInsertWith(f OptionGetOrInsertWithFunc[T]) *T {
	if o.IsSome() {
		return o.v
	}
	val := f()
	o.v = &val
	return o.v
}

// Take takes the value out of the Option, leaving a None in its place.
func (o *Option[T]) Take() Option[T] {
	old := *o
	o.v = nil
	return old
}

// Replace replaces the actual value in the Option by the value given in parameter, returning the old value if present, leaving a Some in its place without deinitializing either one.
func (o *Option[T]) Replace(value T) Option[T] {
	old := *o
	o.v = &value
	return old
}

// Equal checks if two Option`s values are equal
func (o Option[T]) Equal(opt Option[T]) bool {
	if o.IsSome() && opt.IsSome() {
		return reflect.DeepEqual(o.Unwrap(), o.Unwrap())
	}
	return o.IsNone() == opt.IsNone()
}

// Expect returns the contained Some value, consuming the self value.
//
// Panics if the value is a None with a custom panic message provided by msg.
func (o Option[T]) Expect(err any) T {
	if o.IsNone() {
		panic(fmt.Errorf("%v", err))
	}
	return o.Unwrap()
}
