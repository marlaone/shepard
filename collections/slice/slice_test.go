package slice_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice_Capacity(t *testing.T) {
	s := slice.WithCapacity[int](5)
	assert.Equal(t, 5, s.Capacity())
}

func TestSlice_Reserve(t *testing.T) {
	s := slice.New[int]()
	s.Push(1)
	assert.Equal(t, 1, s.Capacity())
	s.Reserve(10)
	assert.Equal(t, 11, s.Capacity())
}

func TestSlice_Truncate(t *testing.T) {
	s := slice.Init(1, 2, 3)
	s.Truncate(2)
	assert.Equal(t, slice.Init(1, 2), s)
}

func TestSlice_Insert(t *testing.T) {
	s := slice.Init(1, 2, 3)
	s.Insert(1, 4)
	assert.Equal(t, slice.Init(1, 4, 2, 3), s)
	s.Insert(4, 5)
	assert.Equal(t, slice.Init(1, 4, 2, 3, 5), s)
}

func TestSlice_Remove(t *testing.T) {
	s := slice.Init(1, 2, 3)
	assert.Equal(t, 2, s.Remove(1))
	assert.Equal(t, slice.Init(1, 3), s)
}

func TestSlice_Retain(t *testing.T) {
	s := slice.Init(1, 2, 3, 4)
	s.Retain(func(x *int) bool { return *x%2 == 0 })
	assert.Equal(t, slice.Init(2, 4), s)
}

func TestSlice_Push(t *testing.T) {
	s := slice.Init(1, 2)
	s.Push(3)
	assert.Equal(t, slice.Init(1, 2, 3), s)
}

func TestSlice_Pop(t *testing.T) {
	s := slice.Init(1, 2, 3)
	assert.True(t, s.Pop().Equal(shepard.Some(3)))
	assert.Equal(t, slice.Init(1, 2), s)
}

func TestSlice_Append(t *testing.T) {
	s := slice.Init(1, 2, 3)
	s2 := slice.Init(4, 5, 6)
	s.Append(&s2)

	assert.Equal(t, slice.Init(1, 2, 3, 4, 5, 6), s)
	assert.Equal(t, slice.New[int](), s2)
}

func TestSlice_Clear(t *testing.T) {
	s := slice.Init(1, 2, 3)
	s.Clear()
	assert.True(t, s.IsEmpty())
}

func TestSlice_Len(t *testing.T) {
	s := slice.Init(1, 2, 3)
	assert.Equal(t, 3, s.Len())
}

func TestSlice_IsEmpty(t *testing.T) {
	s := slice.New[int]()
	assert.True(t, s.IsEmpty())
	s.Push(1)
	assert.False(t, s.IsEmpty())
}

func TestSlice_First(t *testing.T) {
	s := slice.Init(10, 40, 30)
	v := 10
	assert.True(t, s.First().Equal(shepard.Some[*int](&v)))

	s = slice.New[int]()
	assert.True(t, s.First().Equal(shepard.None[*int]()))
}

func TestSlice_Last(t *testing.T) {
	s := slice.Init(10, 40, 30)
	v := 30
	assert.True(t, s.Last().Equal(shepard.Some[*int](&v)))

	s = slice.New[int]()
	assert.True(t, s.Last().Equal(shepard.None[*int]()))
}

func TestSlice_Get(t *testing.T) {
	s := slice.Init(10, 40, 30)
	v := 40
	assert.True(t, s.Get(1).Equal(shepard.Some[*int](&v)))
	assert.True(t, s.Get(3).Equal(shepard.None[*int]()))
}

func TestSlice_Swap(t *testing.T) {
	v := slice.Init("a", "b", "c", "d", "e")
	v.Swap(2, 4)
	assert.Equal(t, slice.Init("a", "b", "e", "d", "c"), v)
}

func TestSlice_Reverse(t *testing.T) {
	v := slice.Init(1, 2, 3)
	v.Reverse()
	assert.Equal(t, slice.Init(3, 2, 1), v)
}

func TestSlice_Iter(t *testing.T) {
	x := slice.Init(1, 2, 4)
	iter := x.Iter()

	assert.True(t, iter.Next().Equal(shepard.Some(1)))
	assert.True(t, iter.Next().Equal(shepard.Some(2)))
	assert.True(t, iter.Next().Equal(shepard.Some(4)))
	assert.True(t, iter.Next().Equal(shepard.None[int]()))
}
