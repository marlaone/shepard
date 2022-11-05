package num_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseString_Int(t *testing.T) {
	assert.True(t, num.ParseString[int]("1").Ok().Equal(shepard.Some[int](1)))
	assert.True(t, num.ParseString[int8]("1").Ok().Equal(shepard.Some[int8](1)))
	assert.True(t, num.ParseString[int16]("1").Ok().Equal(shepard.Some[int16](1)))
	assert.True(t, num.ParseString[int32]("1").Ok().Equal(shepard.Some[int32](1)))
	assert.True(t, num.ParseString[int64]("1").Ok().Equal(shepard.Some[int64](1)))

	assert.True(t, num.ParseString[int8]("one").Ok().Equal(shepard.None[int8]()))
}

func TestParseString_Float(t *testing.T) {
	assert.True(t, num.ParseString[float32]("1").Ok().Equal(shepard.Some[float32](1)))
	assert.True(t, num.ParseString[float32]("1.2").Ok().Equal(shepard.Some[float32](1.2)))
	assert.True(t, num.ParseString[float64]("1.2").Ok().Equal(shepard.Some[float64](1.2)))

	assert.True(t, num.ParseString[float64]("one").Ok().Equal(shepard.None[float64]()))
}

func TestParseString_Uint(t *testing.T) {
	assert.True(t, num.ParseString[uint]("1").Ok().Equal(shepard.Some[uint](1)))
	assert.True(t, num.ParseString[uint8]("1").Ok().Equal(shepard.Some[uint8](1)))
	assert.True(t, num.ParseString[uint16]("1").Ok().Equal(shepard.Some[uint16](1)))
	assert.True(t, num.ParseString[uint32]("1").Ok().Equal(shepard.Some[uint32](1)))
	assert.True(t, num.ParseString[uint64]("1").Ok().Equal(shepard.Some[uint64](1)))

	assert.True(t, num.ParseString[uint8]("-1").Ok().Equal(shepard.None[uint8]()))
	assert.True(t, num.ParseString[uint8]("one").Ok().Equal(shepard.None[uint8]()))
}
