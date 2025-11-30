package slice

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
)

// Iterator is an interface for lazy iterators that can be collected.
type Iterator[T any] interface {
	Next() shepard.Option[T]
}

// SizedIterator is an iterator that knows its length.
type SizedIterator[T any] interface {
	Iterator[T]
	Count() int
}

// Collect converts an Iter into a Slice using zero-copy when possible.
func Collect[T any](it iter.Iter[T]) Slice[T] {
	values := it.Values()
	if len(values) == 0 {
		return New[T]()
	}
	s := Slice[T]{values: make([]T, len(values))}
	copy(s.values, values)
	return s
}

// CollectFilter converts a FilterIter into a Slice in a single pass.
func CollectFilter[T any](it iter.FilterIter[T]) Slice[T] {
	collected := it.Collect()
	if collected.Count() == 0 {
		return New[T]()
	}
	return Slice[T]{values: collected.Values()}
}

// CollectIter collects any lazy iterator into a Slice.
// If the iterator implements SizedIterator, it preallocates capacity.
func CollectIter[T any](it Iterator[T]) Slice[T] {
	// Check if we can get a size hint
	if sized, ok := it.(SizedIterator[T]); ok {
		count := sized.Count()
		if count <= 0 {
			return New[T]()
		}
		values := make([]T, 0, count)
		for {
			opt := it.Next()
			if opt.IsNone() {
				break
			}
			values = append(values, opt.Unwrap())
		}
		return Slice[T]{values: values}
	}

	// No size hint, grow dynamically
	var values []T
	for {
		opt := it.Next()
		if opt.IsNone() {
			break
		}
		values = append(values, opt.Unwrap())
	}
	if len(values) == 0 {
		return New[T]()
	}
	return Slice[T]{values: values}
}

// CollectInto collects a lazy iterator into the provided slice, reusing capacity.
// Returns the new slice with appended elements.
func CollectInto[T any](it Iterator[T], dst []T) []T {
	dst = dst[:0] // Reset length, keep capacity
	for {
		opt := it.Next()
		if opt.IsNone() {
			break
		}
		dst = append(dst, opt.Unwrap())
	}
	return dst
}

// CollectMapIter collects a MapIter into a Slice.
func CollectMapIter[T any, U any](it iter.MapIter[T, U]) Slice[U] {
	collected := it.Collect()
	if collected.Count() == 0 {
		return New[U]()
	}
	return Slice[U]{values: collected.Values()}
}

// CollectTakeIter collects a TakeIter into a Slice.
func CollectTakeIter[T any](it iter.TakeIter[T]) Slice[T] {
	collected := it.Collect()
	if collected.Count() == 0 {
		return New[T]()
	}
	return Slice[T]{values: collected.Values()}
}

// CollectChainIter collects a ChainIter into a Slice.
func CollectChainIter[T any](it iter.ChainIter[T]) Slice[T] {
	collected := it.Collect()
	if collected.Count() == 0 {
		return New[T]()
	}
	return Slice[T]{values: collected.Values()}
}

// CollectFilterMapIter collects a FilterMapIter into a Slice.
func CollectFilterMapIter[T any, U any](it iter.FilterMapIter[T, U]) Slice[U] {
	collected := it.Collect()
	if collected.Count() == 0 {
		return New[U]()
	}
	return Slice[U]{values: collected.Values()}
}

// CollectUniqueIter collects a UniqueIter into a Slice.
func CollectUniqueIter[T comparable](it iter.UniqueIter[T]) Slice[T] {
	collected := it.Collect()
	if collected.Count() == 0 {
		return New[T]()
	}
	return Slice[T]{values: collected.Values()}
}
