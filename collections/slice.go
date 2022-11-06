package collections

import (
	"fmt"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
)

type RetainFunc[T any] func(e *T) bool

type Slice[T any] struct {
	values []T
}

func Init[T any](values ...T) Slice[T] {
	s := WithCapacity[T](len(values))
	for _, v := range values {
		s.Push(v)
	}
	return s
}

func New[T any]() Slice[T] {
	return Slice[T]{
		values: make([]T, 0),
	}
}

func WithCapacity[T any](capacity int) Slice[T] {
	return Slice[T]{
		values: make([]T, 0, capacity),
	}
}

// Capacity returns the number of elements the vector can hold without reallocating.
func (s Slice[T]) Capacity() int {
	return cap(s.values)
}

// Reserve reserves capacity for at least additional more elements to be inserted in the given Slice[T].
func (s *Slice[T]) Reserve(capacity int) {
	newSlice := make([]T, len(s.values), cap(s.values)+capacity)
	s.values = newSlice
}

// Truncate shortens the vector, keeping the first len elements and dropping the rest.
//
// If len is greater than the vector’s current length, this has no effect.
//
// The drain method can emulate truncate, but causes the excess elements to be returned instead of dropped.
//
// Note that this method has no effect on the allocated capacity of the slice.
func (s *Slice[T]) Truncate(length uint) {
	s.values = s.values[:length]
}

// Insert inserts an element at position index within the vector, shifting all elements after it to the right.
//
// Panics if index > len.
func (s *Slice[T]) Insert(index int, element T) {
	if index > len(s.values) {
		panic(fmt.Errorf("index %d out of range for []%T of length %d", index, element, len(s.values)))
	} else if index == len(s.values) {
		s.Push(element)
	} else {
		s.values = append(s.values[:index+1], s.values[index:]...)
		s.values[index] = element
	}
}

// Remove removes and returns the element at position index within the slice, shifting all elements after it to the left.
//
// Panics if index is out of bounds.
func (s *Slice[T]) Remove(index int) T {
	if index >= len(s.values) {
		var zero T
		panic(fmt.Errorf("index %d out of range for []%T of length %d", index, zero, len(s.values)))
	}
	element := s.values[index]
	s.values = append(s.values[:index], s.values[index+1:]...)
	return element
}

// Retain retains only the elements specified by the predicate.
//
// In other words, remove all elements e for which f(&e) returns false.
// This method operates in place, visiting each element exactly once in the original order,
// and preserves the order of the retained elements.
func (s *Slice[T]) Retain(f RetainFunc[T]) {
	*s = CollectSlice[T](s.Iter().Filter(iter.FilterFunc[T](f)))
}

// Push appends an element to the back of a collection.
func (s *Slice[T]) Push(value T) {
	s.values = append(s.values, value)
}

// Pop removes the last element from a slice and returns it, or shepard.None if it is empty.
func (s *Slice[T]) Pop() shepard.Option[T] {
	if len(s.values) > 0 {
		value := s.values[len(s.values)-1]
		s.values = s.values[:len(s.values)-1]
		return shepard.Some(value)
	}
	return shepard.None[T]()
}

// Append moves all the elements of other into self, leaving other empty.
func (s *Slice[T]) Append(other *Slice[T]) {
	s.values = append(s.values, other.values...)
	other.Clear()
}

// Clear clears the vector, removing all values.
//
// Note that this method has no effect on the allocated capacity of the slice.
func (s *Slice[T]) Clear() {
	s.values = s.values[:0]
}

// Len returns the number of elements in the slice, also referred to as its ‘length’.
func (s Slice[T]) Len() int {
	return len(s.values)
}

// IsEmpty returns true if the slice has a length of 0.
func (s Slice[T]) IsEmpty() bool {
	return len(s.values) == 0
}

// First returns the first element of the slice, or shepard.None if it is empty.
func (s Slice[T]) First() shepard.Option[*T] {
	if len(s.values) > 0 {
		return shepard.Some(&s.values[0])
	}
	return shepard.None[*T]()
}

// Last returns the last element of the slice, or shepard.None if it is empty.
func (s Slice[T]) Last() shepard.Option[*T] {
	sliceLen := len(s.values)
	if sliceLen > 0 {
		return shepard.Some(&s.values[sliceLen-1])
	}
	return shepard.None[*T]()
}

// Get returns a reference to an element of index. Returns shepard.None if slice is empty or index is out of bounce.
func (s Slice[T]) Get(index int) shepard.Option[*T] {
	sliceLen := len(s.values)
	if sliceLen == 0 || index >= sliceLen {
		return shepard.None[*T]()
	}
	return shepard.Some(&s.values[index])
}

// Swap swaps two elements in the slice.
//
// Panics if a or b are out of bounds.
func (s *Slice[T]) Swap(a int, b int) {
	sliceLen := len(s.values)
	if a >= sliceLen {
		var zero T
		panic(fmt.Sprintf("index %d is out of bounce for []%T with length %d", a, zero, sliceLen))
	} else if b >= sliceLen {
		var zero T
		panic(fmt.Sprintf("index %d is out of bounce for []%T with length %d", b, zero, sliceLen))
	}
	s.values[a], s.values[b] = s.values[b], s.values[a]
}

// Reverse reverses the order of elements in the slice, in place.
func (s *Slice[T]) Reverse() {
	for i, j := 0, len(s.values)-1; i < j; i, j = i+1, j-1 {
		s.values[i], s.values[j] = s.values[j], s.values[i]
	}
}

// Iter returns an iterator over the slice.
//
// The iterator yields all items from start to end.
func (s Slice[T]) Iter() iter.Iter[T] {
	return iter.New(s.values)
}
