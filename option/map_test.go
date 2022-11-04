package option_test

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/option"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	maybeSomeString := shepard.Some[string]("Hello, World!")
	maybeSomeLen := option.Map[string, int](maybeSomeString, func(s string) int { return len(s) })
	assert.True(t, maybeSomeLen.Equal(shepard.Some[int](13)))
}

func TestMapOr(t *testing.T) {
	x := shepard.Some("foo")
	assert.Equal(t, 3, option.MapOr[string, int](x, 42, func(v string) int { return len(v) }))

	y := shepard.None[string]()
	assert.Equal(t, 42, option.MapOr[string, int](y, 42, func(v string) int { return len(v) }))
}

func TestMapOrElse(t *testing.T) {
	k := 21

	x := shepard.Some("foo")
	assert.Equal(t, 3, option.MapOrElse[string, int](x, func() int { return 2 * k }, func(v string) int { return len(v) }))

	y := shepard.None[string]()
	assert.Equal(t, 42, option.MapOrElse[string, int](y, func() int { return 2 * k }, func(v string) int { return len(v) }))
}

func TestMapOrDefault(t *testing.T) {
	x := shepard.Some("foo")
	assert.Equal(t, 3, option.MapOrDefault[string, int](x, func(v string) int { return len(v) }))

	y := shepard.None[string]()
	assert.Equal(t, 0, option.MapOrDefault[string, int](y, func(v string) int { return len(v) }))
}

func BenchmarkMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		option.Map[int, int](shepard.Some[int](n), func(v int) int { return v + 1 })
	}
}

func BenchmarkMap_Complex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		option.Map[int, string](shepard.Some[int](n).AndThen(func(v int) shepard.Option[int] { return shepard.Some(v * v) }), func(v int) string { return strconv.Itoa(v) })
	}
}
