package hashmap

import "github.com/marlaone/shepard"

// implement Clone[T] for HashMap[K, V]

var _ shepard.Clone[HashMap[string, any]] = (*HashMap[string, any])(nil)

func (m *HashMap[K, V]) Clone() HashMap[K, V] {
	return HashMap[K, V]{
		values: m.values,
		keys:   m.keys,
	}
}
