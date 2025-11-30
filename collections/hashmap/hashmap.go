package hashmap

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
	"golang.org/x/exp/constraints"
	"slices"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type HashMap[K constraints.Ordered, V any] struct {
	data map[K]*Entry[K, V]
}

func New[K constraints.Ordered, V any]() HashMap[K, V] {
	return HashMap[K, V]{
		data: make(map[K]*Entry[K, V]),
	}
}

func WithCapacity[K constraints.Ordered, V any](capacity int) HashMap[K, V] {
	return HashMap[K, V]{
		data: make(map[K]*Entry[K, V], capacity),
	}
}

func From[K constraints.Ordered, V any](pairs []Pair[K, V]) HashMap[K, V] {
	hashmap := WithCapacity[K, V](len(pairs))
	for _, p := range pairs {
		hashmap.Insert(p.Key, p.Value)
	}
	return hashmap
}

// Capacity returns the number of elements the map can hold without reallocating.
// Note: Go maps don't expose capacity, this returns current length as a reasonable approximation.
func (m HashMap[K, V]) Capacity() int {
	return len(m.data)
}

// Keys returns an iter.Iter[K] visiting all keys in sorted order.
func (m HashMap[K, V]) Keys() iter.Iter[K] {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return iter.New(keys)
}

// Values returns an iter.Iter[V] visiting all values in key-sorted order.
func (m HashMap[K, V]) Values() iter.Iter[V] {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	values := make([]V, 0, len(m.data))
	for _, k := range keys {
		e := m.data[k]
		values = append(values, *e.OrDefault())
	}
	return iter.New(values)
}

// ValuesMut returns an iter.Iter[*V] visiting all values mutably in key-sorted order.
func (m HashMap[K, V]) ValuesMut() iter.Iter[*V] {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	values := make([]*V, 0, len(m.data))
	for _, k := range keys {
		e := m.data[k]
		values = append(values, e.OrDefault())
	}
	return iter.New(values)
}

// Iter returns an iter.Iter[Pair[*K, *V]] visiting all key-value pairs in key-sorted order.
func (m HashMap[K, V]) Iter() iter.Iter[Pair[*K, *V]] {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	pairs := make([]Pair[*K, *V], 0, len(m.data))
	for i := range keys {
		e := m.data[keys[i]]
		pairs = append(pairs, Pair[*K, *V]{Key: &keys[i], Value: e.OrDefault()})
	}
	return iter.New(pairs)
}

// Len returns the number of elements in the map.
func (m HashMap[K, V]) Len() int {
	return len(m.data)
}

// IsEmpty returns true if the map contains no elements.
func (m HashMap[K, V]) IsEmpty() bool {
	return len(m.data) == 0
}

// Clear clears the map, removing all key-value pairs.
func (m *HashMap[K, V]) Clear() {
	m.data = make(map[K]*Entry[K, V])
}

// Entry gets the given key's corresponding Entry[K, V] in the map for in-place manipulation.
func (m *HashMap[K, V]) Entry(key K) *Entry[K, V] {
	if e, ok := m.data[key]; ok {
		return e
	}
	e := Vacant[K, V](&key)
	m.data[key] = &e
	return &e
}

// Get returns the value corresponding to the key, without modifying the map.
func (m HashMap[K, V]) Get(k K) shepard.Option[*V] {
	if e, ok := m.data[k]; ok && e.IsOccupied() {
		return shepard.Some(e.Value())
	}
	return shepard.None[*V]()
}

// GetKeyValue returns the key-value Pair[*K, *V] corresponding to the supplied key.
func (m HashMap[K, V]) GetKeyValue(k K) shepard.Option[Pair[*K, *V]] {
	if e, ok := m.data[k]; ok && e.IsOccupied() {
		return shepard.Some(Pair[*K, *V]{Key: e.Key(), Value: e.Value()})
	}
	return shepard.None[Pair[*K, *V]]()
}

// ContainsKey returns true if the map contains a value for the specified key.
func (m HashMap[K, V]) ContainsKey(k K) bool {
	e, ok := m.data[k]
	return ok && e.IsOccupied()
}

// Insert inserts a key-value pair into the map.
//
// If the map did not have this key present, shepard.None is returned.
//
// If the map did have this key present, the value is updated, and the old value is returned.
// The key is not updated, though; this matters for types that can be == without being identical.
func (m *HashMap[K, V]) Insert(k K, v V) shepard.Option[V] {
	if existingEntry, ok := m.data[k]; ok && existingEntry.IsOccupied() {
		oldValue := *existingEntry.Value()
		e := Occupied[K, V](&k, &v)
		m.data[k] = &e
		return shepard.Some(oldValue)
	}
	e := Occupied[K, V](&k, &v)
	m.data[k] = &e
	return shepard.None[V]()
}

// Remove removes a key from the map, returning the value at the key if the key was previously in the map.
func (m *HashMap[K, V]) Remove(k K) shepard.Option[V] {
	if e, ok := m.data[k]; ok && e.IsOccupied() {
		value := *e.Value()
		delete(m.data, k)
		return shepard.Some(value)
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
