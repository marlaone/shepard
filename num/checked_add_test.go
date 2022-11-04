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

	x5 := 2.2
	y5 := 2.2
	assert.True(t, num.CheckedAdd(x5, y5).Equal(shepard.Some[float64](4.4)))
}
