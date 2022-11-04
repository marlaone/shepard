package result_test

import (
	"fmt"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/result"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	maybeSomeString := shepard.Ok[string, error]("Hello, World!")
	maybeSomeLen := result.Map[string, int](maybeSomeString, func(s string) int { return len(s) })
	assert.True(t, maybeSomeLen.Equal(shepard.Ok[int, error](13)))
}

func TestMapOr(t *testing.T) {
	x := shepard.Ok[string, string]("foo")
	assert.Equal(t, 3, result.MapOr[string, int](x, 42, func(v string) int { return len(v) }))

	y := shepard.Err[string, string]("failed")
	assert.Equal(t, 42, result.MapOr[string, int](y, 42, func(v string) int { return len(v) }))
}

func TestMapOrElse(t *testing.T) {
	k := 21

	x := shepard.Ok[string, string]("foo")
	assert.Equal(t, 3, result.MapOrElse[string, int](x, func() int { return 2 * k }, func(v string) int { return len(v) }))

	y := shepard.Err[string]("failed")
	assert.Equal(t, 42, result.MapOrElse[string, int](y, func() int { return 2 * k }, func(v string) int { return len(v) }))
}

func TestMapOrDefault(t *testing.T) {
	x := shepard.Ok[string, string]("foo")
	assert.Equal(t, 3, result.MapOrDefault[string, int](x, func(v string) int { return len(v) }))

	y := shepard.Err[string]("failed")
	assert.Equal(t, 0, result.MapOrDefault[string, int](y, func(v string) int { return len(v) }))
}

func TestMapErr(t *testing.T) {
	stringify := result.MapFunc[int, string](func(x int) string {
		return fmt.Sprintf("error code: %d", x)
	})

	x := shepard.Ok[int, int](2)
	assert.True(t, result.MapErr(x, stringify).Equal(shepard.Ok[int, string](2)))

	x2 := shepard.Err[int, int](13)
	assert.True(t, result.MapErr(x2, stringify).Equal(shepard.Err[int, string]("error code: 13")))
}

func BenchmarkMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result.Map[int, int](shepard.Ok[int, string](n), func(v int) int { return v + 1 })
	}
}

func BenchmarkMap_Complex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result.Map[int, string](shepard.Ok[int, string](n).AndThen(func(v int) shepard.Result[int, string] { return shepard.Ok[int, string](v * v) }), func(v int) string { return strconv.Itoa(v) })
	}
}
