package slice

import "github.com/marlaone/shepard"

var _ shepard.Clone[Slice[int]] = &Slice[int]{}

// Clone returns a deep copy of the slice.
func (s *Slice[V]) Clone() Slice[V] {
	clone := make([]V, len(s.values))
	copy(clone, s.values)
	return Slice[V]{values: clone}
}
