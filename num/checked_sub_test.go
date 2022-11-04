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

	x5 := 2.2
	y5 := 2.2
	assert.True(t, num.CheckedSub(x5, y5).Equal(shepard.Some[float64](0)))
}
