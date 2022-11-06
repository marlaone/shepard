package hashmap

import (
	"github.com/marlaone/shepard"
)

type EntryOrInsertWithFunc[V any] func() V
type EntryOrInsertWithKeyFunc[K comparable, V any] func(*K) V
type EntryAndModifyFunc[V any] func(*V)

type Entry[K comparable, V any] struct {
	key      *K
	value    *V
	occupied bool
}

// Occupied is an occupied entry.
func Occupied[K comparable, V any](k *K, v *V) Entry[K, V] {
	return Entry[K, V]{
		key:      k,
		value:    v,
		occupied: true,
	}
}

// Vacant is a vacant entry.
func Vacant[K comparable, V any](k *K) Entry[K, V] {
	return Entry[K, V]{
		key:      k,
		value:    nil,
		occupied: false,
	}
}

func (e Entry[K, V]) IsOccupied() bool {
	return e.occupied
}

// OrInsert ensures a value is in the entry by inserting the default if empty, and returns a mutable reference to the value in the entry.
func (e *Entry[K, V]) OrInsert(defaultValue V) *V {
	if e.IsOccupied() {
		return e.value
	}
	e.value = &defaultValue
	e.occupied = true
	return e.value
}

// OrInsertWith ensures a value is in the entry by inserting the result of the default function if empty, and returns a mutable reference to the value in the entry.
func (e *Entry[K, V]) OrInsertWith(f EntryOrInsertWithFunc[V]) *V {
	if e.IsOccupied() {
		return e.value
	}
	defaultValue := f()
	e.value = &defaultValue
	e.occupied = true
	return e.value
}

// OrInsertWithKey Ensures a value is in the entry by inserting, if empty, the result of the default function.
// This method allows for generating key-derived values for insertion by providing the default function a reference to the key that was moved during the Entry(key) method call.
//
// The reference to the moved key is provided so that cloning or copying the key is unnecessary, unlike with OrInsertWith.
func (e *Entry[K, V]) OrInsertWithKey(f EntryOrInsertWithKeyFunc[K, V]) *V {
	if e.IsOccupied() {
		return e.value
	}
	defaultValue := f(e.key)
	e.value = &defaultValue
	e.occupied = true
	return e.value
}

// Key returns a reference to this Entry’s key.
func (e Entry[K, V]) Key() *K {
	return e.key
}

// Value returns a reference to this Entry’s value.
func (e Entry[K, V]) Value() *V {
	return e.value
}

// AndModify provides in-place mutable access to an occupied Entry before any potential inserts into the map.
func (e *Entry[K, V]) AndModify(f EntryAndModifyFunc[V]) *Entry[K, V] {
	if e.IsOccupied() {
		f(e.value)
	}
	return e
}

// OrDefault ensures a value is in the entry by inserting the default value if empty, and returns a mutable reference to the value in the entry.
func (e *Entry[K, V]) OrDefault() *V {
	if e.IsOccupied() {
		return e.value
	}
	defaultValue := shepard.GetDefault[V]()
	e.value = &defaultValue
	e.occupied = true
	return e.value
}
