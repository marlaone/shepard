package num_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCheckedDiv(t *testing.T) {

	x1 := math.MaxFloat64
	y1 := 0.5
	assert.True(t, num.CheckedDiv(x1, y1).Equal(shepard.None[float64]()))

	x2 := 4
	y2 := 2
	assert.True(t, num.CheckedDiv(x2, y2).Equal(shepard.Some(2)))

	x3 := 4
	y3 := 0
	assert.True(t, num.CheckedDiv(x3, y3).Equal(shepard.None[int]()))

}
