package option_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
	"github.com/marlaone/shepard/option"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAnd(t *testing.T) {
	x := shepard.Some[int](2)
	y := shepard.None[string]()
	assert.True(t, option.And(x, y).Equal(shepard.None[string]()))

	x2 := shepard.None[int]()
	y2 := shepard.Some("foo")
	assert.True(t, option.And(x2, y2).Equal(shepard.None[string]()))

	x3 := shepard.Some[int](2)
	y3 := shepard.Some("foo")
	assert.True(t, option.And(x3, y3).Equal(shepard.Some("foo")))

	x4 := shepard.None[int]()
	y4 := shepard.None[string]()
	assert.True(t, option.And(x4, y4).Equal(shepard.None[string]()))
}

func TestAndThen(t *testing.T) {
	sqThenToString := option.AndThenFunc[int, string](func(x int) shepard.Option[string] {
		return option.Map[int, string](num.CheckedMul(x, x), func(sq int) string { return strconv.Itoa(sq) })
	})

	assert.True(t, option.AndThen(shepard.Some[int](2), sqThenToString).Equal(shepard.Some("4")))
	assert.True(t, option.AndThen(shepard.None[int](), sqThenToString).Equal(shepard.None[string]()))
}
