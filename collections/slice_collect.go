package collections

import "github.com/marlaone/shepard/iter"

func CollectSlice[T any](iter iter.Iter[T]) Slice[T] {
	s := New[T]()
	iter.Foreach(func(_ int, v T) {
		s.Push(v)
	})
	return s
}
