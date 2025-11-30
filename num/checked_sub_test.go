package num_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCheckedSub(t *testing.T) {
	x := int32(2147483327)
	y := int32(2147483327)

	assert.True(t, num.CheckedSub(x, y).Equal(shepard.Some[int32](0)))

	x2 := int32(-2147483327)
	y2 := int32(2147483327)
	assert.True(t, num.CheckedSub(x2, y2).Equal(shepard.None[int32]()))

	x3 := int32(2)
	y3 := int32(2)
	assert.True(t, num.CheckedSub(x3, y3).Equal(shepard.Some[int32](0)))

	x4 := math.MaxFloat64
	y4 := -math.MaxFloat64
	assert.True(t, num.CheckedSub(x4, y4).Equal(shepard.None[float64]()))

	x5 := 2.0
	y5 := 2.0
	assert.True(t, num.CheckedSub(x5, y5).Equal(shepard.Some[float64](0)))
}

func TestCheckedSub_Zero(t *testing.T) {
	// Subtracting zero should always succeed
	assert.True(t, num.CheckedSub(0, 0).Equal(shepard.Some(0)))
	assert.True(t, num.CheckedSub(100, 0).Equal(shepard.Some(100)))
	assert.True(t, num.CheckedSub(-100, 0).Equal(shepard.Some(-100)))
}

func TestCheckedSub_NegativeNegative(t *testing.T) {
	// Negative - negative
	assert.True(t, num.CheckedSub(int32(-100), int32(-50)).Equal(shepard.Some[int32](-50)))
	assert.True(t, num.CheckedSub(int32(-50), int32(-100)).Equal(shepard.Some[int32](50)))
}

func TestCheckedSub_MaxMin(t *testing.T) {
	// Subtracting from max/min values
	assert.True(t, num.CheckedSub(int8(-128), int8(1)).Equal(shepard.None[int8]()))
	assert.True(t, num.CheckedSub(int8(127), int8(-1)).Equal(shepard.None[int8]()))
	assert.True(t, num.CheckedSub(int8(127), int8(0)).Equal(shepard.Some[int8](127)))
}
