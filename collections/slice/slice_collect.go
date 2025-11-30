package slice

import "github.com/marlaone/shepard/iter"

func Collect[T any](iter iter.Iter[T]) Slice[T] {
	s := WithCapacity[T](iter.Count())
	iter.Foreach(func(_ int, v T) {
		s.Push(v)
	})
	return s
}
