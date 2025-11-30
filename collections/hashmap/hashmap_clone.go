package hashmap

import "github.com/marlaone/shepard"

// implement Clone[T] for HashMap[K, V]

var _ shepard.Clone[HashMap[string, any]] = (*HashMap[string, any])(nil)

func (m *HashMap[K, V]) Clone() HashMap[K, V] {
	clone := HashMap[K, V]{
		data: make(map[K]*Entry[K, V], len(m.data)),
	}
	for k, e := range m.data {
		if e.IsOccupied() {
			v := *e.Value()
			newEntry := Occupied[K, V](k, &v)
			clone.data[k] = &newEntry
		}
		// Note: We don't clone vacant entries - they shouldn't exist in the map
		// after the lazy insertion optimization. If somehow there's a vacant entry,
		// we skip it since it represents "no value" anyway.
	}
	return clone
}
