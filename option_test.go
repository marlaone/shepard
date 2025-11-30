package shepard_test

import (
	"errors"
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOption_IsSome(t *testing.T) {
	x := shepard.Some[int](2)
	assert.True(t, x.IsSome())

	x2 := shepard.None[int]()
	assert.False(t, x2.IsSome())
}

func TestOption_IsNone(t *testing.T) {
	x := shepard.Some[int](2)
	assert.False(t, x.IsNone())

	x2 := shepard.None[int]()
	assert.True(t, x2.IsNone())
}

func TestOption_Expect(t *testing.T) {
	x := shepard.Some[string]("value")
	assert.Equal(t, "value", x.Expect("fruits are healthy"))
	assert.PanicsWithError(t, "fruits are healthy", func() {
		x2 := shepard.None[string]()
		x2.Expect("fruits are healthy")
	})
}

func TestOption_Unwrap(t *testing.T) {
	x := shepard.Some[string]("air")
	assert.Equal(t, "air", x.Unwrap())

	assert.PanicsWithError(t, "unwrap on None Option", func() {
		x2 := shepard.None[string]()
		x2.Unwrap()
	})
}

func TestOption_UnwrapOr(t *testing.T) {
	defaultVal := "bike"
	assert.Equal(t, "car", shepard.Some[string]("car").UnwrapOr(defaultVal))
	assert.Equal(t, defaultVal, shepard.None[string]().UnwrapOr(defaultVal))
}

func TestOption_UnwrapOrElse(t *testing.T) {
	k := 10
	double := shepard.OptionUnwrapElseFunc[int](func() int {
		return 2 * k
	})

	assert.Equal(t, 4, shepard.Some[int](4).UnwrapOrElse(double))
	assert.Equal(t, 20, shepard.None[int]().UnwrapOrElse(double))
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	year := shepard.Some("1909")
	badYear := shepard.None[string]()

	assert.Equal(t, "1909", year.UnwrapOrDefault())
	assert.Equal(t, "", badYear.UnwrapOrDefault())
	opt := shepard.None[testutils.TestType]()
	assert.Equal(t, "test", opt.GetOrInsertDefault().Val)
}

func TestOption_OkOr(t *testing.T) {
	x := shepard.Some("foo")
	assert.True(t, x.OkOr(errors.New("bar")).Equal(shepard.Ok[string, error]("foo")))

	x2 := shepard.None[string]()
	assert.True(t, x2.OkOr(errors.New("bar")).Equal(shepard.Err[string, error](errors.New("bar"))))
}

func TestOption_OkOrElse(t *testing.T) {
	barIt := shepard.OptionOkOrElseFunc[string](func() string {
		return "bar"
	})

	x := shepard.Some("foo")
	assert.True(t, x.OkOrElse(barIt).Equal(shepard.Ok[string, string]("foo")))

	x2 := shepard.None[string]()
	assert.True(t, x2.OkOrElse(barIt).Equal(shepard.Err[string, string]("bar")))
}

func TestOption_And(t *testing.T) {
	x := shepard.Some[int](2)
	y := shepard.None[int]()
	assert.True(t, x.And(y).Equal(shepard.None[int]()))

	x2 := shepard.None[string]()
	y2 := shepard.Some[string]("foo")
	assert.True(t, x2.And(y2).Equal(shepard.None[string]()))

	x3 := shepard.None[string]()
	y3 := shepard.None[string]()
	assert.True(t, x3.And(y3).Equal(shepard.None[string]()))

	x4 := shepard.Some[int](2)
	y4 := shepard.Some[int](10)
	assert.True(t, x4.And(y4).Equal(shepard.Some[int](10)))
}

func TestOption_AndThen(t *testing.T) {
	square := shepard.OptionAndThenFunc[int](func(val int) shepard.Option[int] {
		return shepard.Some[int](val * val)
	})

	fail := shepard.OptionAndThenFunc[int](func(val int) shepard.Option[int] {
		return shepard.None[int]()
	})

	assert.True(t, shepard.Some[int](2).AndThen(square).Equal(shepard.Some[int](4)))
	assert.True(t, shepard.Some[int](2).AndThen(square).AndThen(fail).Equal(shepard.None[int]()))
	assert.True(t, shepard.None[int]().AndThen(square).Equal(shepard.None[int]()))
}

func TestOption_Filter(t *testing.T) {
	isEven := shepard.OptionFilterFunc[int](func(n *int) bool {
		return *n%2 == 0
	})

	assert.True(t, shepard.None[int]().Filter(isEven).Equal(shepard.None[int]()))
	assert.True(t, shepard.Some[int](3).Filter(isEven).Equal(shepard.None[int]()))
	assert.True(t, shepard.Some[int](4).Filter(isEven).Equal(shepard.Some[int](4)))
}

func TestOption_Or(t *testing.T) {
	x := shepard.Some[int](2)
	y := shepard.None[int]()
	assert.True(t, x.Or(y).Equal(shepard.Some[int](2)))

	x2 := shepard.None[int]()
	y2 := shepard.Some[int](2)
	assert.True(t, x2.Or(y2).Equal(shepard.Some[int](2)))

	x3 := shepard.None[string]()
	y3 := shepard.None[string]()
	assert.True(t, x3.Or(y3).Equal(shepard.None[string]()))

	x4 := shepard.Some[int](2)
	y4 := shepard.Some[int](10)
	assert.True(t, x4.Or(y4).Equal(shepard.Some[int](2)))
}

func TestOption_OrElse(t *testing.T) {
	nobody := shepard.OptionOrElseFunc[string](func() shepard.Option[string] {
		return shepard.None[string]()
	})
	vikings := shepard.OptionOrElseFunc[string](func() shepard.Option[string] {
		return shepard.Some[string]("vikings")
	})

	assert.True(t, shepard.Some[string]("barbarians").OrElse(vikings).Equal(shepard.Some[string]("barbarians")))
	assert.True(t, shepard.None[string]().OrElse(vikings).Equal(shepard.Some[string]("vikings")))
	assert.True(t, shepard.None[string]().OrElse(nobody).Equal(shepard.None[string]()))
}

func TestOption_Xor(t *testing.T) {
	x := shepard.Some[int](2)
	y := shepard.None[int]()
	assert.True(t, x.Xor(y).Equal(shepard.Some[int](2)))

	x2 := shepard.None[int]()
	y2 := shepard.Some[int](2)
	assert.True(t, x2.Xor(y2).Equal(shepard.Some[int](2)))

	x3 := shepard.Some[int](2)
	y3 := shepard.Some[int](2)
	assert.True(t, x3.Xor(y3).Equal(shepard.None[int]()))

	x4 := shepard.None[int]()
	y4 := shepard.None[int]()
	assert.True(t, x4.Xor(y4).Equal(shepard.None[int]()))
}

func TestOption_Insert(t *testing.T) {
	opt := shepard.None[int]()
	val := opt.Insert(1)
	assert.Equal(t, 1, *val)
	assert.Equal(t, 1, opt.Unwrap())

	val = opt.Insert(2)
	assert.Equal(t, 2, *val)

	*val = 3
	assert.Equal(t, 3, opt.Unwrap())
}

func TestOption_GetOrInsert(t *testing.T) {
	x := shepard.None[int]()

	y := x.GetOrInsert(5)
	assert.Equal(t, 5, *y)
	*y = 7
	assert.True(t, x.Equal(shepard.Some(7)))
}

func TestOption_GetOrInsertDefault(t *testing.T) {
	x := shepard.None[int]()

	y := x.GetOrInsertDefault()
	assert.Equal(t, 0, *y)
	*y = 7
	assert.True(t, x.Equal(shepard.Some(7)))

	x2 := shepard.None[testutils.TestType]()
	y3 := x2.GetOrInsertDefault()
	assert.Equal(t, "test", y3.Val)
	y3.Val = "bla"
	assert.True(t, x2.Equal(shepard.Some[testutils.TestType](testutils.TestType{Val: "bla"})))
}

func TestOption_GetOrInsertWith(t *testing.T) {
	x := shepard.None[int]()

	y := x.GetOrInsertWith(func() int { return 5 })
	assert.Equal(t, 5, *y)
	*y = 7
	assert.True(t, x.Equal(shepard.Some(7)))
}

func TestOption_Take(t *testing.T) {
	x := shepard.Some[int](2)
	y := x.Take()
	assert.True(t, x.Equal(shepard.None[int]()))
	assert.True(t, y.Equal(shepard.Some[int](2)))

	x2 := shepard.None[int]()
	y2 := x2.Take()
	assert.True(t, x2.Equal(shepard.None[int]()))
	assert.True(t, y2.Equal(shepard.None[int]()))
}

func TestOption_Replace(t *testing.T) {
	x := shepard.Some[int](2)
	old := x.Replace(5)
	assert.True(t, x.Equal(shepard.Some(5)))
	assert.True(t, old.Equal(shepard.Some(2)))

	x2 := shepard.None[int]()
	old2 := x2.Replace(3)
	assert.True(t, x2.Equal(shepard.Some(3)))
	assert.True(t, old2.Equal(shepard.None[int]()))
}

func BenchmarkOption_Unwrap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		shepard.Some[int](n).Unwrap()
	}
}

func BenchmarkOption_UnwrapCustomType(b *testing.B) {
	for n := 0; n < b.N; n++ {
		shepard.Some[testutils.TestType](testutils.TestType{}.Default()).Unwrap()
	}
}

func BenchmarkOption_GetOrInsertDefault(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bla := shepard.None[int]()
		bla.GetOrInsertDefault()
	}
}

func BenchmarkOption_GetOrInsertDefaultCustomType(b *testing.B) {
	for n := 0; n < b.N; n++ {
		opt := shepard.None[testutils.TestType]()
		opt.GetOrInsertDefault()
	}
}

func BenchmarkOption_GetOrInsertWith(b *testing.B) {
	for n := 0; n < b.N; n++ {
		opt := shepard.None[int]()
		opt.GetOrInsertWith(func() int { return n })
	}
}

func BenchmarkOption_GetOrInsertWithCustomType(b *testing.B) {
	for n := 0; n < b.N; n++ {
		opt := shepard.None[testutils.TestType]()
		opt.GetOrInsertWith(func() testutils.TestType { return testutils.TestType{Val: "test"} })
	}
}

func BenchmarkOption_Complex(b *testing.B) {
	opt := shepard.None[int]()
	for n := 0; n < b.N; n++ {
		opt.GetOrInsertWith(func() int { return n })
		opt = opt.Filter(func(v *int) bool { return *v == 100 })
		opt.GetOrInsertDefault()
		opt.Unwrap()
	}
}

// ============================================================================
// Comparison Benchmarks: Option[T] vs Native Go Patterns
// ============================================================================

// --- Some/Unwrap vs Pointer Dereference ---

func BenchmarkOption_SomeUnwrap_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		opt := shepard.Some[int](n)
		_ = opt.Unwrap()
	}
}

func BenchmarkNative_PointerDeref_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		val := n
		ptr := &val
		_ = *ptr
	}
}

func BenchmarkOption_SomeUnwrap_CustomType(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		opt := shepard.Some[testutils.TestType](testutils.TestType{Val: "test"})
		_ = opt.Unwrap()
	}
}

func BenchmarkNative_PointerDeref_CustomType(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		val := testutils.TestType{Val: "test"}
		ptr := &val
		_ = *ptr
	}
}

// --- UnwrapOr vs Nil Coalescing ---

func BenchmarkOption_UnwrapOr_Some_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		opt := shepard.Some[int](n)
		_ = opt.UnwrapOr(0)
	}
}

func BenchmarkOption_UnwrapOr_None_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		opt := shepard.None[int]()
		_ = opt.UnwrapOr(42)
	}
}

func BenchmarkNative_NilCoalesce_Some_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		val := n
		ptr := &val
		var result int
		if ptr != nil {
			result = *ptr
		} else {
			result = 0
		}
		_ = result
	}
}

func BenchmarkNative_NilCoalesce_None_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var ptr *int
		var result int
		if ptr != nil {
			result = *ptr
		} else {
			result = 42
		}
		_ = result
	}
}

// --- AndThen vs If Chain ---

func BenchmarkOption_AndThen_Int(b *testing.B) {
	b.ReportAllocs()
	square := shepard.OptionAndThenFunc[int](func(val int) shepard.Option[int] {
		return shepard.Some[int](val * val)
	})
	for n := 0; n < b.N; n++ {
		opt := shepard.Some[int](n)
		_ = opt.AndThen(square)
	}
}

func BenchmarkNative_IfChain_Int(b *testing.B) {
	b.ReportAllocs()
	square := func(val int) *int {
		result := val * val
		return &result
	}
	for n := 0; n < b.N; n++ {
		val := n
		ptr := &val
		var result *int
		if ptr != nil {
			result = square(*ptr)
		}
		_ = result
	}
}

// --- Filter vs If Check ---

func BenchmarkOption_Filter_Int(b *testing.B) {
	b.ReportAllocs()
	isEven := shepard.OptionFilterFunc[int](func(v *int) bool {
		return *v%2 == 0
	})
	for n := 0; n < b.N; n++ {
		opt := shepard.Some[int](n)
		_ = opt.Filter(isEven)
	}
}

func BenchmarkNative_IfFilter_Int(b *testing.B) {
	b.ReportAllocs()
	isEven := func(v int) bool {
		return v%2 == 0
	}
	for n := 0; n < b.N; n++ {
		val := n
		ptr := &val
		var result *int
		if ptr != nil && isEven(*ptr) {
			result = ptr
		}
		_ = result
	}
}

// --- Equal vs Direct Comparison ---

func BenchmarkOption_Equal_Int(b *testing.B) {
	b.ReportAllocs()
	opt1 := shepard.Some[int](42)
	opt2 := shepard.Some[int](42)
	for n := 0; n < b.N; n++ {
		_ = opt1.Equal(opt2)
	}
}

func BenchmarkNative_DirectEqual_Int(b *testing.B) {
	b.ReportAllocs()
	val1 := 42
	val2 := 42
	ptr1 := &val1
	ptr2 := &val2
	for n := 0; n < b.N; n++ {
		var equal bool
		if ptr1 != nil && ptr2 != nil {
			equal = *ptr1 == *ptr2
		} else {
			equal = ptr1 == nil && ptr2 == nil
		}
		_ = equal
	}
}

func BenchmarkOption_Equal_CustomType(b *testing.B) {
	b.ReportAllocs()
	opt1 := shepard.Some[testutils.TestType](testutils.TestType{Val: "test"})
	opt2 := shepard.Some[testutils.TestType](testutils.TestType{Val: "test"})
	for n := 0; n < b.N; n++ {
		_ = opt1.Equal(opt2)
	}
}

func BenchmarkNative_DirectEqual_CustomType(b *testing.B) {
	b.ReportAllocs()
	val1 := testutils.TestType{Val: "test"}
	val2 := testutils.TestType{Val: "test"}
	ptr1 := &val1
	ptr2 := &val2
	for n := 0; n < b.N; n++ {
		var equal bool
		if ptr1 != nil && ptr2 != nil {
			equal = ptr1.Val == ptr2.Val
		} else {
			equal = ptr1 == nil && ptr2 == nil
		}
		_ = equal
	}
}

// --- Chained Operations ---

func BenchmarkOption_ChainedOps_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		opt := shepard.Some[int](n)
		opt = opt.Filter(func(v *int) bool { return *v%2 == 0 })
		opt = opt.AndThen(func(v int) shepard.Option[int] { return shepard.Some[int](v * 2) })
		_ = opt.UnwrapOr(0)
	}
}

func BenchmarkNative_ChainedOps_Int(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		val := n
		ptr := &val
		// Filter: keep if even
		if ptr != nil && *ptr%2 != 0 {
			ptr = nil
		}
		// Map: double the value
		if ptr != nil {
			doubled := *ptr * 2
			ptr = &doubled
		}
		// UnwrapOr
		var result int
		if ptr != nil {
			result = *ptr
		} else {
			result = 0
		}
		_ = result
	}
}
