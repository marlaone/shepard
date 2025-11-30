package iter

import (
	"github.com/marlaone/shepard"
)

// UniqueIter is a lazy iterator that yields only unique elements.
// It requires the element type to be comparable.
type UniqueIter[T comparable] struct {
	upstream *Iter[T]
	seen     map[T]struct{}
}

// Unique creates a lazy iterator that yields only unique elements.
// Elements are deduplicated on-demand when consumed.
func Unique[T comparable](it Iter[T]) UniqueIter[T] {
	return UniqueIter[T]{
		upstream: &it,
		seen:     make(map[T]struct{}),
	}
}

// Next returns the next unique element.
func (u *UniqueIter[T]) Next() shepard.Option[T] {
	for {
		opt := u.upstream.Next()
		if opt.IsNone() {
			return shepard.None[T]()
		}
		val := opt.Unwrap()
		if _, exists := u.seen[val]; !exists {
			u.seen[val] = struct{}{}
			return shepard.Some(val)
		}
	}
}

// Foreach applies a function to each unique element.
func (u UniqueIter[T]) Foreach(op ForeachFunc[T]) {
	seen := make(map[T]struct{})
	idx := 0
	for i := u.upstream.index + 1; i < len(u.upstream.values); i++ {
		val := u.upstream.values[i]
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			op(idx, val)
			idx++
		}
	}
}

// Count returns the number of unique elements.
func (u UniqueIter[T]) Count() int {
	seen := make(map[T]struct{})
	for i := u.upstream.index + 1; i < len(u.upstream.values); i++ {
		seen[u.upstream.values[i]] = struct{}{}
	}
	return len(seen)
}

// Collect evaluates the lazy iterator and returns all unique elements as an Iter.
func (u UniqueIter[T]) Collect() Iter[T] {
	seen := make(map[T]struct{})
	var result []T
	for i := u.upstream.index + 1; i < len(u.upstream.values); i++ {
		val := u.upstream.values[i]
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return New(result)
}
