package shepard_test

import (
	"errors"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOk(t *testing.T) {
	ok1 := shepard.Ok[int, int](1)
	ok2 := shepard.Ok[int, int](1)
	ok3 := shepard.Ok[int, int](2)
	assert.True(t, ok1.Equal(ok2), "Ok(1) and Ok(1) should be equal")
	assert.False(t, ok1.Equal(ok3), "Ok(1) and Ok(2) should not be equal")
}

func TestErr(t *testing.T) {
	err1 := shepard.Err[int, int](1)
	err2 := shepard.Err[int, int](1)
	err3 := shepard.Err[int, int](2)
	assert.True(t, err1.Equal(err2), "Err(1) and Err(1) should be equal")
	assert.False(t, err1.Equal(err3), "Err(1) and Err(2) should not be equal")
}

func TestResult_IsOk(t *testing.T) {
	x := shepard.Ok[int, int](-3)
	assert.True(t, x.IsOk(), "Ok(-3) should be a valid result")

	x2 := shepard.Err[int, string]("Some error message")
	assert.False(t, x2.IsOk(), "Err(\"Some error message\") should not be a valid result")
}

func TestResult_IsErr(t *testing.T) {
	x := shepard.Ok[int, int](-3)
	assert.False(t, x.IsErr(), "Ok(-3) should not be an error")

	x2 := shepard.Err[int, string]("Some error message")
	assert.True(t, x2.IsErr(), "Err(\"Some error message\") should be an error")
}

func TestResult_Ok(t *testing.T) {
	x := shepard.Ok[int, int](2)
	assert.True(t, x.Ok().Equal(shepard.Some[int](2)), "Ok(2) should be Some(2)")

	x2 := shepard.Err[string, string]("Nothing here")
	assert.True(t, x2.Ok().Equal(shepard.None[string]()), "Err(\"Nothing here\") should be None")
}

func TestResult_Err(t *testing.T) {
	x := shepard.Ok[int, int](2)
	assert.True(t, x.Err().Equal(shepard.None[int]()), "Ok(2) should be None")

	x2 := shepard.Err[string, string]("Nothing here")
	assert.True(t, x2.Err().Equal(shepard.Some[string]("Nothing here")), "Err(\"Nothing here\") should be None")
}

func TestResult_Expect(t *testing.T) {
	x := shepard.Err[string, string]("emergency failure")
	assert.PanicsWithError(t, "Testing expect: emergency failure", func() {
		x.Expect("Testing expect")
	})
}

func TestResult_Unwrap(t *testing.T) {
	x := shepard.Ok[int, int](2)
	assert.Equal(t, 2, x.Unwrap(), "Ok(2) should unwrap to 2")

	x2 := shepard.Err[int, string]("emergency failure")
	assert.PanicsWithError(t, "emergency failure", func() {
		x2.Unwrap()
	})
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	x := shepard.Ok[int, int](1909)
	x2 := shepard.Err[int, string]("emergency failure")

	assert.Equal(t, 1909, x.UnwrapOrDefault())
	assert.Equal(t, 0, x2.UnwrapOrDefault())
}

func TestResult_ExpectErr(t *testing.T) {
	x := shepard.Ok[int, string](10)
	assert.PanicsWithError(t, "Testing ExpectErr: 10", func() {
		x.ExpectErr("Testing ExpectErr")
	})
}

func TestResult_UnwrapErr(t *testing.T) {
	x := shepard.Ok[int, string](2)
	assert.PanicsWithError(t, "2", func() {
		x.UnwrapErr()
	})

	x2 := shepard.Err[int, string]("emergency failure")
	assert.Equal(t, "emergency failure", x2.UnwrapErr())
}

func TestResult_And(t *testing.T) {
	x := shepard.Ok[int, string](2)
	y := shepard.Err[int, string]("late error")
	assert.True(t, x.And(y).Equal(shepard.Err[int, string]("late error")))

	x2 := shepard.Err[string, string]("early error")
	y2 := shepard.Ok[string, string]("foo")
	assert.True(t, x2.And(y2).Equal(shepard.Err[string, string]("early error")))

	x3 := shepard.Err[string, string]("not a 2")
	y3 := shepard.Err[string, string]("late error")
	assert.True(t, x3.And(y3).Equal(shepard.Err[string, string]("not a 2")))

	x4 := shepard.Ok[int, string](2)
	y4 := shepard.Ok[int, string](10)
	assert.True(t, x4.And(y4).Equal(shepard.Ok[int, string](10)))
}

func TestResult_AndThen(t *testing.T) {
	sq := shepard.ResultAndThenFunc[int, string](func(val int) shepard.Result[int, string] {
		return shepard.Ok[int, string](val * val)
	})

	assert.True(t, shepard.Ok[int, string](2).AndThen(sq).Equal(shepard.Ok[int, string](4)))
	assert.True(t, shepard.Err[int, string]("not a number").AndThen(sq).Equal(shepard.Err[int, string]("not a number")))
}

func TestResult_Or(t *testing.T) {
	x := shepard.Ok[int, string](2)
	y := shepard.Err[int, string]("late error")
	assert.True(t, x.Or(y).Equal(shepard.Ok[int, string](2)))

	x2 := shepard.Err[int, string]("early error")
	y2 := shepard.Ok[int, string](2)
	assert.True(t, x2.Or(y2).Equal(shepard.Ok[int, string](2)))

	x3 := shepard.Err[string, string]("not a 2")
	y3 := shepard.Err[string, string]("late error")
	assert.True(t, x3.Or(y3).Equal(shepard.Err[string, string]("late error")))

	x4 := shepard.Ok[int, string](2)
	y4 := shepard.Ok[int, string](10)
	assert.True(t, x4.Or(y4).Equal(shepard.Ok[int, string](2)))
}

func TestResult_OrElse(t *testing.T) {
	sq := shepard.ResultOrElseFunc[int, int](func(val int) shepard.Result[int, int] {
		return shepard.Ok[int, int](val * val)
	})
	err := shepard.ResultOrElseFunc[int, int](func(val int) shepard.Result[int, int] {
		return shepard.Err[int, int](val)
	})

	assert.True(t, shepard.Ok[int, int](2).OrElse(sq).OrElse(sq).Equal(shepard.Ok[int, int](2)))
	assert.True(t, shepard.Ok[int, int](2).OrElse(err).OrElse(sq).Equal(shepard.Ok[int, int](2)))
	assert.True(t, shepard.Err[int, int](3).OrElse(sq).OrElse(err).Equal(shepard.Ok[int, int](9)))
	assert.True(t, shepard.Err[int, int](3).OrElse(err).OrElse(err).Equal(shepard.Err[int, int](3)))
}

func TestResult_UnwrapOr(t *testing.T) {
	defaultVal := 2

	x := shepard.Ok[int, int](9)
	assert.Equal(t, x.UnwrapOr(defaultVal), 9)

	x2 := shepard.Err[int, string]("error")
	assert.Equal(t, x2.UnwrapOr(defaultVal), defaultVal)
}

func TestResult_UnwrapOrElse(t *testing.T) {
	count := shepard.ResultUnwrapElseFunc[int, string](func(x string) int {
		return len(x)
	})

	assert.Equal(t, shepard.Ok[int, string](2).UnwrapOrElse(count), 2)
	assert.Equal(t, shepard.Err[int, string]("foo").UnwrapOrElse(count), 3)
}

func TestResult_Equal(t *testing.T) {
	assert.True(t, shepard.Ok[int, int](2).Equal(shepard.Ok[int, int](2)))
	assert.False(t, shepard.Ok[int, int](3).Equal(shepard.Ok[int, int](2)))
	assert.False(t, shepard.Err[int, int](2).Equal(shepard.Ok[int, int](2)))
}

func TestResult_AsMut(t *testing.T) {
	res := shepard.Ok[int, int](5)
	ok, _ := res.AsMut()
	*ok = 7
	assert.Equal(t, 7, res.Unwrap())
}

// ============================================================================
// Comparison Benchmarks: Result[T, E] vs Native Go (value, error) Pattern
// ============================================================================

// --- Ok/Unwrap vs Value Check ---

func BenchmarkResult_OkUnwrap_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res := shepard.Ok[int, string](n)
		_ = res.Unwrap()
	}
}

func BenchmarkNative_ValueCheck_Int(b *testing.B) {
	b.ReportAllocs()
	getValue := func(n int) (int, error) {
		return n, nil
	}
	for n := 0; n < b.N; n++ {
		val, err := getValue(n)
		if err != nil {
			panic(err)
		}
		_ = val
	}
}

func BenchmarkResult_OkUnwrap_CustomType(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res := shepard.Ok[testutils.TestType, string](testutils.TestType{Val: "test"})
		_ = res.Unwrap()
	}
}

func BenchmarkNative_ValueCheck_CustomType(b *testing.B) {
	b.ReportAllocs()
	getValue := func() (testutils.TestType, error) {
		return testutils.TestType{Val: "test"}, nil
	}
	for n := 0; n < b.N; n++ {
		val, err := getValue()
		if err != nil {
			panic(err)
		}
		_ = val
	}
}

// --- Err Creation ---

func BenchmarkResult_ErrCreation_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res := shepard.Err[int, string]("error")
		_ = res.IsErr()
	}
}

func BenchmarkNative_ErrorCreation_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		err := errors.New("error")
		_ = err != nil
	}
}

// --- UnwrapOr vs Error Fallback ---

func BenchmarkResult_UnwrapOr_Ok_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res := shepard.Ok[int, string](n)
		_ = res.UnwrapOr(0)
	}
}

func BenchmarkResult_UnwrapOr_Err_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res := shepard.Err[int, string]("error")
		_ = res.UnwrapOr(42)
	}
}

func BenchmarkNative_ErrorFallback_Ok_Int(b *testing.B) {
	b.ReportAllocs()
	getValue := func(n int) (int, error) {
		return n, nil
	}
	for n := 0; n < b.N; n++ {
		val, err := getValue(n)
		if err != nil {
			val = 0
		}
		_ = val
	}
}

func BenchmarkNative_ErrorFallback_Err_Int(b *testing.B) {
	b.ReportAllocs()
	getValue := func() (int, error) {
		return 0, errors.New("error")
	}
	for n := 0; n < b.N; n++ {
		val, err := getValue()
		if err != nil {
			val = 42
		}
		_ = val
	}
}

// --- AndThen vs If Err Chain ---

func BenchmarkResult_AndThen_Int(b *testing.B) {
	b.ReportAllocs()
	square := shepard.ResultAndThenFunc[int, string](func(val int) shepard.Result[int, string] {
		return shepard.Ok[int, string](val * val)
	})
	for n := 0; n < b.N; n++ {
		res := shepard.Ok[int, string](n)
		_ = res.AndThen(square)
	}
}

func BenchmarkNative_IfErrChain_Int(b *testing.B) {
	b.ReportAllocs()
	getValue := func(n int) (int, error) {
		return n, nil
	}
	square := func(val int) (int, error) {
		return val * val, nil
	}
	for n := 0; n < b.N; n++ {
		val, err := getValue(n)
		if err == nil {
			val, err = square(val)
		}
		_, _ = val, err
	}
}

// --- OrElse vs Error Recovery ---

func BenchmarkResult_OrElse_Int(b *testing.B) {
	b.ReportAllocs()
	recover := shepard.ResultOrElseFunc[int, string](func(err string) shepard.Result[int, string] {
		return shepard.Ok[int, string](0)
	})
	for n := 0; n < b.N; n++ {
		res := shepard.Err[int, string]("error")
		_ = res.OrElse(recover)
	}
}

func BenchmarkNative_ErrorRecover_Int(b *testing.B) {
	b.ReportAllocs()
	getValue := func() (int, error) {
		return 0, errors.New("error")
	}
	recover := func(err error) (int, error) {
		return 0, nil
	}
	for n := 0; n < b.N; n++ {
		val, err := getValue()
		if err != nil {
			val, err = recover(err)
		}
		_, _ = val, err
	}
}

// --- Equal vs Direct Comparison ---

func BenchmarkResult_Equal_Int(b *testing.B) {
	b.ReportAllocs()
	res1 := shepard.Ok[int, string](42)
	res2 := shepard.Ok[int, string](42)
	for n := 0; n < b.N; n++ {
		_ = res1.Equal(res2)
	}
}

func BenchmarkNative_DirectEqual_Result_Int(b *testing.B) {
	b.ReportAllocs()
	val1, err1 := 42, error(nil)
	val2, err2 := 42, error(nil)
	for n := 0; n < b.N; n++ {
		var equal bool
		if err1 == nil && err2 == nil {
			equal = val1 == val2
		} else if err1 != nil && err2 != nil {
			equal = err1.Error() == err2.Error()
		} else {
			equal = false
		}
		_ = equal
	}
}

func BenchmarkResult_Equal_CustomType(b *testing.B) {
	b.ReportAllocs()
	res1 := shepard.Ok[testutils.TestType, string](testutils.TestType{Val: "test"})
	res2 := shepard.Ok[testutils.TestType, string](testutils.TestType{Val: "test"})
	for n := 0; n < b.N; n++ {
		_ = res1.Equal(res2)
	}
}

func BenchmarkNative_DirectEqual_Result_CustomType(b *testing.B) {
	b.ReportAllocs()
	val1, err1 := testutils.TestType{Val: "test"}, error(nil)
	val2, err2 := testutils.TestType{Val: "test"}, error(nil)
	for n := 0; n < b.N; n++ {
		var equal bool
		if err1 == nil && err2 == nil {
			equal = val1.Val == val2.Val
		} else if err1 != nil && err2 != nil {
			equal = err1.Error() == err2.Error()
		} else {
			equal = false
		}
		_ = equal
	}
}

// --- Chained Error Propagation ---

func BenchmarkResult_ChainedOps_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		res := shepard.Ok[int, string](n)
		res = res.AndThen(func(v int) shepard.Result[int, string] {
			if v%2 == 0 {
				return shepard.Ok[int, string](v * 2)
			}
			return shepard.Err[int, string]("odd number")
		})
		res = res.OrElse(func(err string) shepard.Result[int, string] {
			return shepard.Ok[int, string](0)
		})
		_ = res.UnwrapOr(-1)
	}
}

func BenchmarkNative_ChainedErrorOps_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		val := n
		var err error
		// AndThen: double if even, else error
		if err == nil {
			if val%2 == 0 {
				val = val * 2
			} else {
				err = errors.New("odd number")
			}
		}
		// OrElse: recover from error
		if err != nil {
			val = 0
			err = nil
		}
		// UnwrapOr
		if err != nil {
			val = -1
		}
		_ = val
	}
}
