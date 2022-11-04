package option

import "github.com/marlaone/shepard"

type MapFunc[T any, U any] func(value T) U
type MapElseFunc[U any] func() U

// Map maps a shepard.Option[T] to shepard.Option[U] by applying a function to a contained value.
func Map[T any, U any](opt shepard.Option[T], f MapFunc[T, U]) shepard.Option[U] {
	if opt.IsNone() {
		return shepard.None[U]()
	}
	return shepard.Some[U](f(opt.Unwrap()))
}

// MapOr returns the provided default result (if shepard.None), or applies a function to the contained value (if shepard.Some).
//
// Arguments passed to MapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func MapOr[T any, U any](opt shepard.Option[T], defaultValue U, f MapFunc[T, U]) U {
	if opt.IsNone() {
		return defaultValue
	}
	return f(opt.Unwrap())
}

// MapOrElse computes a default function result (if shepard.None), or applies a different function to the contained value (if shepard.Some).
func MapOrElse[T any, U any](opt shepard.Option[T], op MapElseFunc[U], f MapFunc[T, U]) U {
	if opt.IsNone() {
		return op()
	}
	return f(opt.Unwrap())
}

// MapOrDefault returns the types default (if shepard.None), or applies a function to the contained value (if shepard.Some).
func MapOrDefault[T any, U any](opt shepard.Option[T], f MapFunc[T, U]) U {
	if opt.IsNone() {
		var valType U
		defaulter, ok := any(valType).(shepard.Default[U])
		if ok {
			return defaulter.Default()
		}
		return valType
	}
	return f(opt.Unwrap())
}
