package iter

import (
	"github.com/marlaone/shepard"
)

type ForeachFunc[T any] func(index int, value T)
type FilterFunc[T any] func(val *T) bool
type FindFunc[T any] func(val *T) bool

type Iterator[T any] interface {
	Next() shepard.Option[T]
}

type Iter[T any] struct {
	values []T
	index  int
}

func New[T any](values []T) Iter[T] {
	return Iter[T]{
		values: values,
		index:  -1,
	}
}

func (i *Iter[T]) Next() shepard.Option[T] {
	if i.index >= len(i.values)-1 {
		return shepard.None[T]()
	}
	i.index++
	return shepard.Some(i.values[i.index])
}

func (i Iter[T]) Foreach(op ForeachFunc[T]) {
	for {
		next := i.Next()
		if next.IsNone() {
			break
		}
		op(i.index, next.Unwrap())
	}
}

// Take creates an iterator that yields the first n elements, or fewer if the underlying iterator ends sooner.
//
// take(n) yields elements until n elements are yielded or the end of the iterator is reached (whichever happens first).
// The returned iterator is a prefix of length n if the original iterator contains at least n elements,
// otherwise it contains all the (fewer than n) elements of the original iterator.
func (i Iter[T]) Take(n int) Iter[T] {
	valuesLen := len(i.values)
	if valuesLen > 0 {
		if n >= valuesLen {
			n = valuesLen
		}
		return New(i.values[:n])
	}
	return i
}

// Filter creates an iterator which uses a closure to determine if an element should be yielded.
//
// Given an element the closure must return true or false.
// The returned iterator will yield only the elements for which the closure returns true.
func (i Iter[T]) Filter(predicate FilterFunc[T]) Iter[T] {
	newValues := make([]T, 0, cap(i.values))
	for {
		next := i.Next()
		if next.IsNone() {
			break
		}
		v := next.Unwrap()
		if predicate(&v) {
			newValues = append(newValues, next.Unwrap())
		}
	}
	return New(newValues)
}

// Find searches for an element of an iterator that satisfies a predicate.
//
// find() takes a closure that returns true or false.
// It applies this closure to each element of the iterator, and if any of them return true, then find() returns shepard.Some(element).
// If they all return false, it returns shepard.None.
//
// find() is short-circuiting; in other words, it will stop processing as soon as the closure returns true.
func (i Iter[T]) Find(predicate FindFunc[T]) shepard.Option[T] {
	for {
		next := i.Next()
		if next.IsNone() {
			return next
		}
		v := next.Unwrap()
		if predicate(&v) {
			return next
		}
	}
}

func (i Iter[T]) Count() int {
	return len(i.values)
}
