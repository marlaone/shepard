package result_test

import (
	"errors"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
	"github.com/marlaone/shepard/result"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAnd(t *testing.T) {
	x := shepard.Ok[int, string](2)
	y := shepard.Err[string, string]("failed")
	assert.True(t, result.And(x, y).Equal(shepard.Err[string, string]("failed")))

	x2 := shepard.Err[int, string]("failed")
	y2 := shepard.Ok[string, string]("foo")
	assert.True(t, result.And(x2, y2).Equal(shepard.Err[string]("failed")))

	x3 := shepard.Ok[int, string](2)
	y3 := shepard.Ok[string, string]("foo")
	assert.True(t, result.And(x3, y3).Equal(shepard.Ok[string, string]("foo")))

	x4 := shepard.Err[int, string]("not a 2")
	y4 := shepard.Err[string, string]("late error")
	assert.True(t, result.And(x4, y4).Equal(shepard.Err[string, string]("not a 2")))
}

func TestAndThen(t *testing.T) {
	sqThenToString := result.AndThenFunc[int, string, error](func(x int) shepard.Result[string, error] {
		return result.Map[int, string](num.CheckedMul(x, x).OkOr(errors.New("failed")), func(sq int) string { return strconv.Itoa(sq) })
	})

	assert.True(t, result.AndThen(shepard.Ok[int, error](2), sqThenToString).Equal(shepard.Ok[string, error]("4")))
	assert.True(t, result.AndThen(shepard.Err[int, error](errors.New("failed")), sqThenToString).Equal(shepard.Err[string, error](errors.New("failed"))))
}
