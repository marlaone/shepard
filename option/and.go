package option

import "github.com/marlaone/shepard"

type AndThenFunc[T any, U any] func(v T) shepard.Option[U]

// And returns shepard.None if option is shepard.None, otherwise returns optb.
//
// Arguments passed to and are eagerly evaluated; if you are passing the result of a function call, it is recommended to use AndThen, which is lazily evaluated.
func And[T any, U any](opt shepard.Option[T], optb shepard.Option[U]) shepard.Option[U] {
	if opt.IsSome() && optb.IsSome() {
		return optb
	}
	return shepard.None[U]()
}

// AndThen returns None if the shepard.Option is shepard.None, otherwise calls f with the wrapped value and returns the result.
//
// Some languages call this operation flatmap.
func AndThen[T any, U any](opt shepard.Option[T], f AndThenFunc[T, U]) shepard.Option[U] {
	if opt.IsNone() {
		return shepard.None[U]()
	}
	return f(opt.Unwrap())
}
