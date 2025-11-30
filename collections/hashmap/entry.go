package hashmap

import (
	"github.com/marlaone/shepard"
	"golang.org/x/exp/constraints"
)

type EntryOrInsertWithFunc[V any] func() V
type EntryOrInsertWithKeyFunc[K comparable, V any] func(*K) V
type EntryAndModifyFunc[V any] func(*V)

// Entry represents either an occupied or vacant entry in a HashMap.
// Like Rust's Entry enum, it enables efficient in-place operations on map entries.
// The entry is returned by value to avoid heap allocations (Rust-aligned optimization).
// Key is stored by value to prevent heap escape in vacant entries.
type Entry[K constraints.Ordered, V any] struct {
	key      K                // Store by value to prevent heap escape
	value    *V
	hashmap  *HashMap[K, V]   // Back-reference for lazy insertion (Rust-aligned: VacantEntry holds reference to map)
	isVacant bool             // Explicit state tracking (vacant entries have nil value but need distinction)
}

// Occupied creates an occupied entry (key exists with a value).
func Occupied[K constraints.Ordered, V any](k K, v *V) Entry[K, V] {
	return Entry[K, V]{
		key:      k,
		value:    v,
		hashmap:  nil, // Occupied entries don't need back-reference
		isVacant: false,
	}
}

// Vacant creates a vacant entry (key does not exist yet).
// The hashmap parameter is needed for lazy insertion when OrInsert is called.
func Vacant[K constraints.Ordered, V any](k K, hashmap *HashMap[K, V]) Entry[K, V] {
	return Entry[K, V]{
		key:      k,
		value:    nil,
		hashmap:  hashmap, // Back-reference for lazy insertion
		isVacant: true,
	}
}

func (e Entry[K, V]) IsOccupied() bool {
	return !e.isVacant && e.value != nil
}

// OrInsert ensures a value is in the entry by inserting the default if empty, and returns a mutable reference to the value in the entry.
// For vacant entries, this performs lazy insertion into the HashMap (Rust-aligned behavior).
func (e *Entry[K, V]) OrInsert(defaultValue V) *V {
	if e.IsOccupied() {
		return e.value
	}
	e.value = &defaultValue
	e.isVacant = false
	// Lazy insertion: only store in map when actually inserting a value
	if e.hashmap != nil {
		newEntry := Entry[K, V]{key: e.key, value: e.value, isVacant: false}
		e.hashmap.data[e.key] = &newEntry
	}
	return e.value
}

// OrInsertWith ensures a value is in the entry by inserting the result of the default function if empty, and returns a mutable reference to the value in the entry.
// For vacant entries, this performs lazy insertion into the HashMap (Rust-aligned behavior).
func (e *Entry[K, V]) OrInsertWith(f EntryOrInsertWithFunc[V]) *V {
	if e.IsOccupied() {
		return e.value
	}
	defaultValue := f()
	e.value = &defaultValue
	e.isVacant = false
	// Lazy insertion: only store in map when actually inserting a value
	if e.hashmap != nil {
		newEntry := Entry[K, V]{key: e.key, value: e.value, isVacant: false}
		e.hashmap.data[e.key] = &newEntry
	}
	return e.value
}

// OrInsertWithKey Ensures a value is in the entry by inserting, if empty, the result of the default function.
// This method allows for generating key-derived values for insertion by providing the default function a reference to the key that was moved during the Entry(key) method call.
//
// The reference to the moved key is provided so that cloning or copying the key is unnecessary, unlike with OrInsertWith.
// For vacant entries, this performs lazy insertion into the HashMap (Rust-aligned behavior).
func (e *Entry[K, V]) OrInsertWithKey(f EntryOrInsertWithKeyFunc[K, V]) *V {
	if e.IsOccupied() {
		return e.value
	}
	defaultValue := f(&e.key)
	e.value = &defaultValue
	e.isVacant = false
	// Lazy insertion: only store in map when actually inserting a value
	if e.hashmap != nil {
		newEntry := Entry[K, V]{key: e.key, value: e.value, isVacant: false}
		e.hashmap.data[e.key] = &newEntry
	}
	return e.value
}

// Key returns this Entry's key by value.
func (e Entry[K, V]) Key() K {
	return e.key
}

// Value returns a reference to this Entry’s value.
func (e Entry[K, V]) Value() *V {
	return e.value
}

// AndModify provides in-place mutable access to an occupied Entry before any potential inserts into the map.
// The closure receives *V (mutable reference to the value).
func (e *Entry[K, V]) AndModify(f EntryAndModifyFunc[V]) *Entry[K, V] {
	if e.IsOccupied() {
		f(e.value)
	}
	return e
}

// OrDefault ensures a value is in the entry by inserting the default value if empty, and returns a mutable reference to the value in the entry.
// For vacant entries, this performs lazy insertion into the HashMap (Rust-aligned behavior).
func (e *Entry[K, V]) OrDefault() *V {
	if e.IsOccupied() {
		return e.value
	}
	defaultValue := shepard.GetDefault[V]()
	e.value = &defaultValue
	e.isVacant = false
	// Lazy insertion: only store in map when actually inserting a value
	if e.hashmap != nil {
		newEntry := Entry[K, V]{key: e.key, value: e.value, isVacant: false}
		e.hashmap.data[e.key] = &newEntry
	}
	return e.value
}
