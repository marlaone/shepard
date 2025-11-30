package iter

import "github.com/marlaone/shepard"

type MapFunc[T any, U any] func(item T) U
type FilterMapFunc[T any, U any] func(item T) shepard.Option[U]

// Map takes a closure and creates an iterator which calls that closure on each element.
func Map[T any, U any](it Iter[T], f MapFunc[T, U]) Iter[U] {
	remaining := len(it.values) - it.index - 1
	newValues := make([]U, remaining)
	for idx := it.index + 1; idx < len(it.values); idx++ {
		newValues[idx-it.index-1] = f(it.values[idx])
	}
	return New[U](newValues)
}

// FilterMap an iterator that uses f to both filter and map elements from iter.
// Combines filter and map operations in a single pass for better performance.
func FilterMap[T any, U any](it Iter[T], f FilterMapFunc[T, U]) Iter[U] {
	remaining := len(it.values) - it.index - 1
	newValues := make([]U, 0, remaining)
	for idx := it.index + 1; idx < len(it.values); idx++ {
		if result := f(it.values[idx]); result.IsSome() {
			newValues = append(newValues, result.Unwrap())
		}
	}
	return New[U](newValues)
}
