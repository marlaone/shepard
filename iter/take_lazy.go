package iter

import (
	"github.com/marlaone/shepard"
)

// TakeIter is a lazy iterator that yields at most n elements.
type TakeIter[T any] struct {
	upstream  *Iter[T]
	remaining int
}

// TakeLazy creates a lazy take iterator.
// Unlike Take(), this does not create a new slice immediately.
// It wraps the upstream iterator and limits the number of elements.
func (i Iter[T]) TakeLazy(n int) TakeIter[T] {
	return TakeIter[T]{
		upstream:  &i,
		remaining: n,
	}
}

// Next returns the next element if remaining count > 0.
func (t *TakeIter[T]) Next() shepard.Option[T] {
	if t.remaining <= 0 {
		return shepard.None[T]()
	}
	opt := t.upstream.Next()
	if opt.IsSome() {
		t.remaining--
	}
	return opt
}

// Foreach applies a function to each element (up to n).
func (t TakeIter[T]) Foreach(op ForeachFunc[T]) {
	count := 0
	for i := t.upstream.index + 1; i < len(t.upstream.values) && count < t.remaining; i++ {
		op(count, t.upstream.values[i])
		count++
	}
}

// Count returns the number of elements that will be yielded.
func (t TakeIter[T]) Count() int {
	available := len(t.upstream.values) - t.upstream.index - 1
	if available < t.remaining {
		return available
	}
	return t.remaining
}

// Collect evaluates the lazy iterator and returns elements as an Iter.
func (t TakeIter[T]) Collect() Iter[T] {
	count := t.Count()
	if count <= 0 {
		return New[T](nil)
	}
	result := make([]T, count)
	for i := 0; i < count; i++ {
		result[i] = t.upstream.values[t.upstream.index+1+i]
	}
	return New(result)
}
