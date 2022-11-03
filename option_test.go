package shepard_test

import (
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
	assert.True(t, x.OkOr("bar").Equal(shepard.Ok[string, string]("foo")))

	x2 := shepard.None[string]()
	assert.True(t, x2.OkOr("bar").Equal(shepard.Err[string, string]("bar")))
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
	square := shepard.OptionAndThenFunc[int](func(val *int) shepard.Option[int] {
		*val = *val * *val
		return shepard.Some[int](*val)
	})

	fail := shepard.OptionAndThenFunc[int](func(val *int) shepard.Option[int] {
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
