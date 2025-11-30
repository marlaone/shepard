package hashmap_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/iter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntry_OrInsert(t *testing.T) {
	// Pass nil for hashmap since we're not testing map insertion here
	e := hashmap.Vacant[string, int]("poneyland", nil)

	assert.Equal(t, 3, *e.OrInsert(3))
	assert.Equal(t, 3, *e.OrInsert(4))
	assert.Equal(t, 3, *e.Value())
}

func TestEntry_OrInsertWith(t *testing.T) {
	v := "hoho"
	e := hashmap.Vacant[string, string]("poneyland", nil)

	assert.Equal(t, "hoho", *e.OrInsertWith(func() string { return v }))
	assert.Equal(t, "hoho", *e.OrInsertWith(func() string { return "hihi" }))
}

func TestEntry_OrInsertWithKey(t *testing.T) {
	e := hashmap.Vacant[string, int]("poneyland", nil)

	assert.Equal(t, 9, *e.OrInsertWithKey(func(k *string) int { return iter.New[rune]([]rune(*k)).Count() }))
	assert.Equal(t, 9, *e.OrInsertWithKey(func(k *string) int { return 123 }))
}

func TestEntry_Key(t *testing.T) {
	e := hashmap.Vacant[string, string]("poneyland", nil)
	assert.Equal(t, "poneyland", e.Key()) // Key now returns value, not pointer
}

func TestEntry_AndModify(t *testing.T) {
	e := hashmap.Vacant[string, int]("poneyland", nil)

	e.AndModify(func(v *int) { *v += 1 }).OrInsert(42)
	assert.Equal(t, 42, *e.Value())

	e.AndModify(func(v *int) { *v += 1 }).OrInsert(42)
	assert.Equal(t, 43, *e.Value())
}

func TestEntry_OrDefault(t *testing.T) {
	e := hashmap.Vacant[string, shepard.Option[int]]("poneyland", nil)
	e.OrDefault()
	assert.True(t, e.Value().Equal(shepard.None[int]()))
}

// Test Entry with actual HashMap to verify lazy insertion
func TestEntry_LazyInsertion(t *testing.T) {
	m := hashmap.New[string, int]()

	// Get entry for non-existing key
	entry := m.Entry("key1")
	assert.False(t, entry.IsOccupied())

	// Map should still be empty (lazy insertion)
	assert.Equal(t, 0, m.Len())

	// Now insert
	entry.OrInsert(42)

	// Now map should have the entry
	assert.Equal(t, 1, m.Len())
	assert.True(t, m.ContainsKey("key1"))

	opt := m.Get("key1")
	assert.True(t, opt.IsSome())
	assert.Equal(t, 42, *opt.Unwrap())
}

// Test Entry for existing key
func TestEntry_ExistingKey(t *testing.T) {
	m := hashmap.New[string, int]()
	m.Insert("key1", 10)

	// Get entry for existing key
	entry := m.Entry("key1")
	assert.True(t, entry.IsOccupied())
	assert.Equal(t, 10, *entry.Value())

	// Modify via AndModify
	entry.AndModify(func(v *int) { *v += 5 })
	assert.Equal(t, 15, *entry.Value())

	// Verify map was updated
	opt := m.Get("key1")
	assert.True(t, opt.IsSome())
	assert.Equal(t, 15, *opt.Unwrap())
}
