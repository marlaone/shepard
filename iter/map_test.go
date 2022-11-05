package iter_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
	"github.com/marlaone/shepard/num"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	i := iter.New([]int{1, 2, 3})

	assert.Equal(t, iter.New[string]([]string{"1", "2", "3"}), iter.Map[int, string](i, func(x int) string { return strconv.Itoa(x) }))
}

func TestFilterMap(t *testing.T) {
	i := iter.FilterMap[string, uint8](iter.New([]string{"1", "two", "NaN", "four", "5", "-5"}), func(x string) shepard.Option[uint8] { return num.ParseString[uint8](x).Ok() })
	assert.True(t, i.Next().Equal(shepard.Some[uint8](1)))
	assert.True(t, i.Next().Equal(shepard.Some[uint8](5)))
	assert.True(t, i.Next().Equal(shepard.None[uint8]()))
}
