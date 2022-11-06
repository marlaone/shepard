package hashmap

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
	"golang.org/x/exp/constraints"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type HashMap[K constraints.Ordered, V any] struct {
	values []*Entry[K, V]
	keys   []K
}

func New[K constraints.Ordered, V any]() HashMap[K, V] {
	return HashMap[K, V]{
		values: make([]*Entry[K, V], 0),
		keys:   make([]K, 0),
	}
}

func WithCapacity[K constraints.Ordered, V any](capacity int) HashMap[K, V] {
	return HashMap[K, V]{
		values: make([]*Entry[K, V], 0, capacity),
		keys:   make([]K, 0, capacity),
	}
}

func From[K constraints.Ordered, V any](pairs []Pair[K, V]) HashMap[K, V] {
	hashmap := New[K, V]()
	for _, p := range pairs {
		hashmap.Insert(p.Key, p.Value)
	}
	return hashmap
}

// Capacity returns the number of elements the map can hold without reallocating.
func (m HashMap[K, V]) Capacity() int {
	return cap(m.keys)
}

// Keys returns an iter.Iter[K] visiting all keys in arbitrary order.
func (m HashMap[K, V]) Keys() iter.Iter[K] {
	return iter.New(m.keys)
}

// Values returns an iter.Iter[V] visiting all values in arbitrary order.
func (m HashMap[K, V]) Values() iter.Iter[V] {
	var values []V
	for _, e := range m.values {
		values = append(values, *e.OrDefault())
	}
	return iter.New(values)
}

// ValuesMut returns an iter.Iter[*V] visiting all values mutably in arbitrary order.
func (m HashMap[K, V]) ValuesMut() iter.Iter[*V] {
	var values []*V
	for _, e := range m.values {
		values = append(values, e.OrDefault())
	}
	return iter.New(values)
}

// Iter returns an iter.Iter[Pair[*K, *V]] visiting all key-value pairs in arbitrary order.
func (m HashMap[K, V]) Iter() iter.Iter[Pair[*K, *V]] {
	var values []Pair[*K, *V]
	for i, k := range m.keys {
		v := m.values[i]
		values = append(values, Pair[*K, *V]{Key: &k, Value: v.OrDefault()})
	}
	return iter.New(values)
}

// Len returns the number of elements in the map.
func (m HashMap[K, V]) Len() int {
	return len(m.keys)
}

// IsEmpty returns true if the map contains no elements.
func (m HashMap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Clear clears the map, removing all key-value pairs. Keeps the allocated memory for reuse.
func (m *HashMap[K, V]) Clear() {
	m.values = make([]*Entry[K, V], 0)
	m.keys = make([]K, 0)
}

func (m HashMap[K, V]) keyIndex(key K) (int, bool) {
	for i, j := 0, m.Len()-1; i < m.Len() && j >= 0; i, j = i+1, j-1 {
		if m.keys[i] == key {
			return i, true
		} else if m.keys[j] == key {
			return j, true
		}
	}
	return 0, false
}

// Entry gets the given key’s corresponding Entry[K, V] in the map for in-place manipulation.
func (m *HashMap[K, V]) Entry(key K) *Entry[K, V] {
	i, ok := m.keyIndex(key)
	if !ok {
		e := Vacant[K, V](&key)
		m.values = append(m.values, &e)
		m.keys = append(m.keys, key)
		return &e
	}
	e := m.values[i]
	return e
}

func (m HashMap[K, V]) Get(k K) shepard.Option[*V] {
	e := m.Entry(k)
	if e.IsOccupied() {
		return shepard.Some(e.Value())
	}
	return shepard.None[*V]()
}

// GetKeyValue returns the key-value Pair[*K, *V] corresponding to the supplied key.
func (m HashMap[K, V]) GetKeyValue(k K) shepard.Option[Pair[*K, *V]] {
	e := m.Entry(k)
	if e.IsOccupied() {
		return shepard.Some(Pair[*K, *V]{Key: e.Key(), Value: e.Value()})
	}
	return shepard.None[Pair[*K, *V]]()
}

// ContainsKey returns true if the map contains a value for the specified key.
func (m HashMap[K, V]) ContainsKey(k K) bool {
	return m.Entry(k).IsOccupied()
}

// Insert inserts a key-value pair into the map.
//
// If the map did not have this key present, shepard.None is returned.
//
// If the map did have this key present, the value is updated, and the old value is returned.
// The key is not updated, though; this matters for types that can be == without being identical.
func (m *HashMap[K, V]) Insert(k K, v V) shepard.Option[V] {
	i, ok := m.keyIndex(k)
	e := Occupied[K, V](&k, &v)
	if ok {
		m.values[i] = &e
		return shepard.Some(*e.OrDefault())
	} else {
		m.values = append(m.values, &e)
		m.keys = append(m.keys, k)
	}
	return shepard.None[V]()
}

// Remove removes a key from the map, returning the value at the key if the key was previously in the map.
func (m *HashMap[K, V]) Remove(k K) shepard.Option[V] {
	i, ok := m.keyIndex(k)
	if ok {
		e := m.values[i]

		m.values = append(m.values[:i], m.values[i+1:]...)
		m.keys = append(m.keys[:i], m.keys[i+1:]...)

		return shepard.Some[V](*e.Value())
	}
	return shepard.None[V]()
}

// RemoveEntry removes a key from the map, returning the stored key and value if the key was previously in the map.
func (m *HashMap[K, V]) RemoveEntry(k K) shepard.Option[Pair[K, V]] {
	opt := m.Remove(k)
	if opt.IsSome() {
		return shepard.Some[Pair[K, V]](Pair[K, V]{Key: k, Value: opt.Unwrap()})
	}
	return shepard.None[Pair[K, V]]()
}
