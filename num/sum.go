package num

import "github.com/marlaone/shepard/iter"

// Sum takes an iterator and generates T from the elements by “summing up” the items.
func Sum[T Number](iterator iter.Iter[T]) T {
	var sum T
	iterator.Foreach(func(_ int, v T) {
		sum += v
	})
	return sum
}
