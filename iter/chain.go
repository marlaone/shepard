package iter

import (
	"github.com/marlaone/shepard"
)

// ChainIter is a lazy iterator that chains two iterators together.
type ChainIter[T any] struct {
	first    *Iter[T]
	second   *Iter[T]
	onSecond bool
}

// Chain creates a lazy iterator that chains this iterator with another.
// First yields all elements from self, then all elements from other.
func (i Iter[T]) Chain(other Iter[T]) ChainIter[T] {
	return ChainIter[T]{
		first:    &i,
		second:   &other,
		onSecond: false,
	}
}

// Next returns the next element from either the first or second iterator.
func (c *ChainIter[T]) Next() shepard.Option[T] {
	if !c.onSecond {
		opt := c.first.Next()
		if opt.IsSome() {
			return opt
		}
		c.onSecond = true
	}
	return c.second.Next()
}

// Foreach applies a function to each element from both iterators.
func (c ChainIter[T]) Foreach(op ForeachFunc[T]) {
	idx := 0
	// First iterator
	for i := c.first.index + 1; i < len(c.first.values); i++ {
		op(idx, c.first.values[i])
		idx++
	}
	// Second iterator
	for i := c.second.index + 1; i < len(c.second.values); i++ {
		op(idx, c.second.values[i])
		idx++
	}
}

// Count returns the total number of elements from both iterators.
func (c ChainIter[T]) Count() int {
	firstRemaining := len(c.first.values) - c.first.index - 1
	secondRemaining := len(c.second.values) - c.second.index - 1
	return firstRemaining + secondRemaining
}

// Collect evaluates the lazy iterator and returns all elements as an Iter.
func (c ChainIter[T]) Collect() Iter[T] {
	count := c.Count()
	if count <= 0 {
		return New[T](nil)
	}
	result := make([]T, 0, count)
	for i := c.first.index + 1; i < len(c.first.values); i++ {
		result = append(result, c.first.values[i])
	}
	for i := c.second.index + 1; i < len(c.second.values); i++ {
		result = append(result, c.second.values[i])
	}
	return New(result)
}
