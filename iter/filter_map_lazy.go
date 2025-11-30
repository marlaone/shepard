package iter

import (
	"github.com/marlaone/shepard"
)

// FilterMapLazyFunc is a function that optionally transforms an element.
type FilterMapLazyFunc[T any, U any] func(item T) shepard.Option[U]

// FilterMapIter is a lazy iterator that filters and maps elements on-demand.
type FilterMapIter[T any, U any] struct {
	upstream *Iter[T]
	mapper   FilterMapLazyFunc[T, U]
}

// FilterMapLazy creates a lazy filter-map iterator.
// Unlike FilterMap(), this does not allocate a new slice immediately.
// Elements are filtered and transformed on-demand when consumed.
func FilterMapLazy[T any, U any](it Iter[T], mapper FilterMapLazyFunc[T, U]) FilterMapIter[T, U] {
	return FilterMapIter[T, U]{
		upstream: &it,
		mapper:   mapper,
	}
}

// Next returns the next element that passes the filter and is mapped.
func (fm *FilterMapIter[T, U]) Next() shepard.Option[U] {
	for fm.upstream.index++; fm.upstream.index < len(fm.upstream.values); fm.upstream.index++ {
		result := fm.mapper(fm.upstream.values[fm.upstream.index])
		if result.IsSome() {
			return result
		}
	}
	return shepard.None[U]()
}

// Foreach applies a function to each filtered and mapped element.
func (fm FilterMapIter[T, U]) Foreach(op ForeachFunc[U]) {
	idx := 0
	for i := fm.upstream.index + 1; i < len(fm.upstream.values); i++ {
		result := fm.mapper(fm.upstream.values[i])
		if result.IsSome() {
			op(idx, result.Unwrap())
			idx++
		}
	}
}

// Count returns the number of elements that pass the filter.
func (fm FilterMapIter[T, U]) Count() int {
	count := 0
	for i := fm.upstream.index + 1; i < len(fm.upstream.values); i++ {
		result := fm.mapper(fm.upstream.values[i])
		if result.IsSome() {
			count++
		}
	}
	return count
}

// Collect evaluates the lazy iterator and returns all elements as an Iter.
func (fm FilterMapIter[T, U]) Collect() Iter[U] {
	var result []U
	for i := fm.upstream.index + 1; i < len(fm.upstream.values); i++ {
		opt := fm.mapper(fm.upstream.values[i])
		if opt.IsSome() {
			result = append(result, opt.Unwrap())
		}
	}
	return New(result)
}
