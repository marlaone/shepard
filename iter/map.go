package iter

import "github.com/marlaone/shepard"

type MapFunc[T any, U any] func(item T) U
type FilterMapFunc[T any, U any] func(item T) shepard.Option[U]

// Map takes a closure and creates an iterator which calls that closure on each element.
func Map[T any, U any](iter Iter[T], f MapFunc[T, U]) Iter[U] {
	var newValues []U

	iter.Foreach(func(_ int, item T) {
		newValues = append(newValues, f(item))
	})

	return New[U](newValues)
}

// FilterMap an iterator that uses f to both filter and map elements from iter.
func FilterMap[T any, U any](iter Iter[T], f FilterMapFunc[T, U]) Iter[U] {
	var mappedValues []shepard.Option[U]

	iter.Foreach(func(_ int, item T) {
		mappedValues = append(mappedValues, f(item))
	})

	return Map[shepard.Option[U], U](New[shepard.Option[U]](mappedValues).Filter(func(i *shepard.Option[U]) bool { return i.IsSome() }), func(item shepard.Option[U]) U { return item.Unwrap() })
}
