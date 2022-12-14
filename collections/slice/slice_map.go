package slice

import "github.com/marlaone/shepard/iter"

type MapFunc[T any, U any] func(item T) U

// Map takes a closure and creates an iterator which calls that closure on each element.
func Map[T any, U any](s Slice[T], f MapFunc[T, U]) iter.Iter[U] {
	return iter.Map[T, U](s.Iter(), iter.MapFunc[T, U](f))
}
