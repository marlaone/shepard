package hashmap_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/iter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntry_OrInsert(t *testing.T) {
	k := "poneyland"
	e := hashmap.Vacant[string, int](&k)

	assert.Equal(t, 3, *e.OrInsert(3))
	assert.Equal(t, 3, *e.OrInsert(4))
	assert.Equal(t, 3, *e.Value())
}

func TestEntry_OrInsertWith(t *testing.T) {
	k := "poneyland"
	v := "hoho"
	e := hashmap.Vacant[string, string](&k)

	assert.Equal(t, "hoho", *e.OrInsertWith(func() string { return v }))
	assert.Equal(t, "hoho", *e.OrInsertWith(func() string { return "hihi" }))
}

func TestEntry_OrInsertWithKey(t *testing.T) {
	k := "poneyland"
	e := hashmap.Vacant[string, int](&k)

	assert.Equal(t, 9, *e.OrInsertWithKey(func(k *string) int { return iter.New[rune]([]rune(*k)).Count() }))
	assert.Equal(t, 9, *e.OrInsertWithKey(func(k *string) int { return 123 }))
}

func TestEntry_Key(t *testing.T) {
	k := "poneyland"
	e := hashmap.Vacant[string, string](&k)
	assert.Equal(t, "poneyland", *e.Key())
}

func TestEntry_AndModify(t *testing.T) {
	k := "poneyland"
	e := hashmap.Vacant[string, int](&k)

	e.AndModify(func(v *int) { *v += 1 })
	e.OrInsert(42)
	assert.Equal(t, 42, *e.Value())

	e.AndModify(func(v *int) { *v += 1 })
	e.OrInsert(42)
	assert.Equal(t, 43, *e.Value())
}

func TestEntry_OrDefault(t *testing.T) {
	k := "poneyland"
	e := hashmap.Vacant[string, shepard.Option[int]](&k)
	e.OrDefault()
	assert.True(t, e.Value().Equal(shepard.None[int]()))
}
