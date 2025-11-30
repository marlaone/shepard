package num_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCheckedAdd(t *testing.T) {
	x := int32(2147483327)
	y := int32(2147483327)

	assert.True(t, num.CheckedAdd(x, y).Equal(shepard.None[int32]()))

	x2 := int32(-2147483327)
	y2 := int32(-2147483327)
	assert.True(t, num.CheckedAdd(x2, y2).Equal(shepard.None[int32]()))

	x3 := int32(2)
	y3 := int32(2)
	assert.True(t, num.CheckedAdd(x3, y3).Equal(shepard.Some[int32](4)))

	x4 := math.MaxFloat64
	y4 := math.MaxFloat64
	assert.True(t, num.CheckedAdd(x4, y4).Equal(shepard.None[float64]()))

	x5 := 2.0
	y5 := 2.0
	assert.True(t, num.CheckedAdd(x5, y5).Equal(shepard.Some[float64](4.0)))
}

func TestCheckedAdd_Zero(t *testing.T) {
	// Adding zero should always succeed
	assert.True(t, num.CheckedAdd(0, 0).Equal(shepard.Some(0)))
	assert.True(t, num.CheckedAdd(100, 0).Equal(shepard.Some(100)))
	assert.True(t, num.CheckedAdd(0, 100).Equal(shepard.Some(100)))
	assert.True(t, num.CheckedAdd(-100, 0).Equal(shepard.Some(-100)))
}

func TestCheckedAdd_NegativePositive(t *testing.T) {
	// Negative + positive that doesn't overflow
	assert.True(t, num.CheckedAdd(int32(-100), int32(50)).Equal(shepard.Some[int32](-50)))
	assert.True(t, num.CheckedAdd(int32(100), int32(-50)).Equal(shepard.Some[int32](50)))
	assert.True(t, num.CheckedAdd(int32(-100), int32(100)).Equal(shepard.Some[int32](0)))
}

func TestCheckedAdd_MaxMin(t *testing.T) {
	// Adding to max/min values
	assert.True(t, num.CheckedAdd(int8(127), int8(1)).Equal(shepard.None[int8]()))
	assert.True(t, num.CheckedAdd(int8(-128), int8(-1)).Equal(shepard.None[int8]()))
	assert.True(t, num.CheckedAdd(int8(127), int8(0)).Equal(shepard.Some[int8](127)))
}
