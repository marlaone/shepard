package iter

import (
	"github.com/marlaone/shepard"
)

// MapFunc is a function that transforms an element.
type MapLazyFunc[T any, U any] func(item T) U

// MapIter is a lazy mapping iterator that transforms elements on-demand.
type MapIter[T any, U any] struct {
	upstream *Iter[T]
	mapper   MapLazyFunc[T, U]
}

// MapLazy creates a lazy mapping iterator.
// Unlike Map(), this does not allocate a new slice immediately.
// Elements are transformed on-demand when Next() or Collect() is called.
func MapLazy[T any, U any](it Iter[T], mapper MapLazyFunc[T, U]) MapIter[T, U] {
	return MapIter[T, U]{
		upstream: &it,
		mapper:   mapper,
	}
}

// Next returns the next mapped element.
func (m *MapIter[T, U]) Next() shepard.Option[U] {
	opt := m.upstream.Next()
	if opt.IsNone() {
		return shepard.None[U]()
	}
	return shepard.Some(m.mapper(opt.Unwrap()))
}

// Foreach applies a function to each mapped element.
func (m MapIter[T, U]) Foreach(op ForeachFunc[U]) {
	idx := 0
	for i := m.upstream.index + 1; i < len(m.upstream.values); i++ {
		mapped := m.mapper(m.upstream.values[i])
		op(idx, mapped)
		idx++
	}
}

// Count returns the number of elements remaining.
func (m MapIter[T, U]) Count() int {
	return len(m.upstream.values) - m.upstream.index - 1
}

// Collect evaluates the lazy iterator and returns all mapped elements as an Iter.
func (m MapIter[T, U]) Collect() Iter[U] {
	remaining := len(m.upstream.values) - m.upstream.index - 1
	if remaining <= 0 {
		return New[U](nil)
	}
	result := make([]U, remaining)
	for i := m.upstream.index + 1; i < len(m.upstream.values); i++ {
		result[i-m.upstream.index-1] = m.mapper(m.upstream.values[i])
	}
	return New(result)
}
