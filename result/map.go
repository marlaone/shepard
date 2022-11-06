package result

import "github.com/marlaone/shepard"

type MapFunc[T any, U any] func(value T) U
type MapElseFunc[U any] func() U

// Map maps a shepard.Result[T, E] to shepard.Result[U, E] by applying a function to a contained shepard.Ok value, leaving a shepard.Err value untouched.
//
// This function can be used to compose the shepard.Result`s of two functions.
func Map[T any, U any, E any](res shepard.Result[T, E], f MapFunc[T, U]) shepard.Result[U, E] {
	if res.IsErr() {
		return shepard.Err[U, E](res.Err().Unwrap())
	}
	return shepard.Ok[U, E](f(res.Unwrap()))
}

// MapOr returns the provided default (if shepard.Err), or applies a function to the contained value (if shepard.Ok),
//
// Arguments passed to MapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func MapOr[T any, U any, E any](res shepard.Result[T, E], defaultValue U, f MapFunc[T, U]) U {
	if res.IsErr() {
		return defaultValue
	}
	return f(res.Unwrap())
}

// MapOrElse Maps a Result<T, E> to U by applying fallback function default to a contained Err value, or function f to a contained Ok value.
//
// This function can be used to unpack a successful result while handling an error.
func MapOrElse[T any, U any, E any](res shepard.Result[T, E], op MapElseFunc[U], f MapFunc[T, U]) U {
	if res.IsErr() {
		return op()
	}
	return f(res.Unwrap())
}

// MapOrDefault returns the types default (if shepard.Err), or applies a function to the contained value (if shepard.Ok).
func MapOrDefault[T any, U any, E any](res shepard.Result[T, E], f MapFunc[T, U]) U {
	if res.IsErr() {
		return shepard.GetDefault[U]()
	}
	return f(res.Unwrap())
}

// MapErr maps a shepard.Result[T, E] to shepard.Result[T, F] by applying a function to a contained shepard.Err value, leaving an shepard.Ok value untouched.
//
// This function can be used to pass through a successful result while handling an error.
func MapErr[T any, E any, F any](res shepard.Result[T, E], op MapFunc[E, F]) shepard.Result[T, F] {
	if res.IsOk() {
		return shepard.Ok[T, F](res.Unwrap())
	}
	return shepard.Err[T, F](op(res.Err().Unwrap()))
}
