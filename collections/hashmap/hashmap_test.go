package hashmap_test

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/iter"
	"github.com/stretchr/testify/assert"
)

func TestHashMap_Capacity(t *testing.T) {
	m := hashmap.WithCapacity[string, int](10)
	assert.Equal(t, 10, m.Capacity())
}

func TestHashMap_Keys(t *testing.T) {
	m := hashmap.From[string, int]([]hashmap.Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}})

	assert.Equal(t, iter.New[string]([]string{"a", "b", "c"}), m.Keys())
}

func TestHashMap_Values(t *testing.T) {
	m := hashmap.From[string, int]([]hashmap.Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}})

	assert.Equal(t, iter.New[int]([]int{1, 2, 3}), m.Values())
}

func TestHashMap_ValuesMut(t *testing.T) {
	m := hashmap.From[string, int]([]hashmap.Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}})

	m.ValuesMut().Foreach(func(_ int, value *int) {
		*value = *value * 2
	})

	assert.EqualValues(t, iter.New[int]([]int{2, 4, 6}), m.Values())
}

func TestHashMap_Iter(t *testing.T) {
	pairs := []hashmap.Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
	m := hashmap.From[string, int](pairs)

	var expected []hashmap.Pair[*string, *int]
	for i := 0; i < len(pairs); i++ {
		p := pairs[i]
		expected = append(expected, hashmap.Pair[*string, *int]{Key: &p.Key, Value: &p.Value})
	}

	assert.EqualValues(
		t,
		iter.Map[hashmap.Pair[*string, *int], int](iter.New[hashmap.Pair[*string, *int]](expected), func(v hashmap.Pair[*string, *int]) int { return *v.Value }),
		iter.Map[hashmap.Pair[*string, *int], int](m.Iter(), func(v hashmap.Pair[*string, *int]) int { return *v.Value }),
	)
}

func TestHashMap_Len(t *testing.T) {
	m := hashmap.New[int, string]()
	assert.Equal(t, 0, m.Len())
	m.Insert(1, "a")
	assert.Equal(t, 1, m.Len())
}

func TestHashMap_IsEmpty(t *testing.T) {
	m := hashmap.New[int, string]()
	assert.True(t, m.IsEmpty())
	m.Insert(1, "a")
	assert.False(t, m.IsEmpty())
}

func TestHashMap_Clear(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Insert(1, "a")
	assert.False(t, m.IsEmpty())
	m.Clear()
	assert.True(t, m.IsEmpty())
}

func TestHashMap_Entry(t *testing.T) {
	letters := hashmap.New[rune, int]()

	for _, ch := range []rune("a short treatise on fungi") {
		letters.Entry(ch).AndModify(func(counter *int) {
			*counter += 1
		}).OrInsert(1)
	}

	assert.Equal(t, 2, *letters.Get('s').Unwrap())
	assert.Equal(t, 3, *letters.Get('t').Unwrap())
	assert.Equal(t, 1, *letters.Get('u').Unwrap())
	assert.True(t, letters.Get('y').Equal(shepard.None[*int]()))
}

func TestHashMap_Get(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Insert(1, "a")
	expected := "a"
	assert.True(t, m.Get(1).Equal(shepard.Some[*string](&expected)))
}

func TestHashMap_GetKeyValue(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Insert(1, "a")

	expectedKey := 1
	expectedValue := "a"
	assert.True(t, m.GetKeyValue(1).Equal(shepard.Some[hashmap.Pair[*int, *string]](hashmap.Pair[*int, *string]{&expectedKey, &expectedValue})))
	assert.True(t, m.GetKeyValue(2).Equal(shepard.None[hashmap.Pair[*int, *string]]()))
}

func TestHashMap_Insert(t *testing.T) {
	m := hashmap.New[int, string]()

	assert.True(t, m.Insert(37, "a").Equal(shepard.None[string]()))
	assert.False(t, m.IsEmpty())
	m.Insert(37, "b")
	assert.True(t, m.Insert(37, "c").Equal(shepard.Some[string]("b")))
	assert.Equal(t, "c", *m.Get(37).Unwrap())
}

func TestHashMap_Remove(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Insert(1, "a")
	assert.True(t, m.Remove(1).Equal(shepard.Some[string]("a")))
	assert.True(t, m.Remove(1).Equal(shepard.None[string]()))
}

func TestHashMap_RemoveEntry(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Insert(1, "a")
	assert.True(t, m.RemoveEntry(1).Equal(shepard.Some[hashmap.Pair[int, string]](hashmap.Pair[int, string]{1, "a"})))
	assert.True(t, m.RemoveEntry(1).Equal(shepard.None[hashmap.Pair[int, string]]()))
	assert.True(t, m.Remove(1).Equal(shepard.None[string]()))
}

func TestHashMap_ContainsKey(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Insert(1, "a")
	assert.True(t, m.ContainsKey(1))
	assert.False(t, m.ContainsKey(2))
}

func BenchmarkHashMap_Get(b *testing.B) {
	m := hashmap.WithCapacity[int, int](b.N)
	for i := 0; i < b.N; i++ {
		m.Insert(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i)
	}
	b.ReportAllocs()
}

func BenchmarkGoMap(b *testing.B) {
	m := make(map[int]int, b.N)
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ := m[i]
		void(v)
	}
	b.ReportAllocs()
}

func void(v int) {}
