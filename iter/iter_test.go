package iter_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIter_Next(t *testing.T) {
	it := iter.New([]int{1, 2, 3})

	assert.True(t, it.Next().Equal(shepard.Some(1)))
	assert.True(t, it.Next().Equal(shepard.Some(2)))
	assert.True(t, it.Next().Equal(shepard.Some(3)))
}

func TestIter_Foreach(t *testing.T) {
	var indices []int
	var values []int

	expectedIndices := []int{0, 1, 2}
	expectedValues := []int{1, 2, 3}

	it := iter.New([]int{1, 2, 3})

	it.Foreach(func(index int, value int) {
		indices = append(indices, index)
		values = append(values, value)
	})

	assert.Equal(t, expectedIndices, indices)
	assert.Equal(t, expectedValues, values)
}

func TestIter_Take(t *testing.T) {
	values := []int{1, 2, 3}

	it := iter.New(values).Take(2)

	assert.True(t, it.Next().Equal(shepard.Some(1)))
	assert.True(t, it.Next().Equal(shepard.Some(2)))
	assert.True(t, it.Next().Equal(shepard.None[int]()))
}

func TestIter_Filter(t *testing.T) {
	it := iter.New([]int{1, 2, 3, 4, 5, 6}).Filter(func(v *int) bool { return *v%2 == 0 })

	assert.True(t, it.Next().Equal(shepard.Some(2)))
	assert.True(t, it.Next().Equal(shepard.Some(4)))
	assert.True(t, it.Next().Equal(shepard.Some(6)))
	assert.True(t, it.Next().Equal(shepard.None[int]()))
}

func TestIter_Find(t *testing.T) {
	item := iter.New([]int{1, 2, 3, 4, 5, 6}).Find(func(v *int) bool { return *v%2 == 0 })
	assert.True(t, item.Equal(shepard.Some(2)))

	item = iter.New([]int{1, 2, 3, 4, 5, 6}).Find(func(v *int) bool { return *v == 100 })
	assert.True(t, item.Equal(shepard.None[int]()))
}
