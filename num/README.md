# num

Implements checked implementation for mul/div/add/sub to keep you save from overflows.

## Number

```go
type Number interface {
	constraints.Integer | constraints.Float
}
```

## Add

`num.CheckedAdd[T num.Number](n T, v T) shepard.Option[T]`

### Examples
```go
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
```

## Sub

`num.CheckedSub[T num.Number](n T, v T) shepard.Option[T]`

### Examples

```go
func TestCheckedSub(t *testing.T) {
	x := int32(2147483327)
	y := int32(2147483327)

	assert.True(t, num.CheckedSub(x, y).Equal(shepard.Some[int32](0)))

	x2 := int32(-2147483327)
	y2 := int32(2147483327)
	assert.True(t, num.CheckedSub(x2, y2).Equal(shepard.None[int32]()))

	x3 := int32(2)
	y3 := int32(2)
	assert.True(t, num.CheckedSub(x3, y3).Equal(shepard.Some[int32](0)))

	x4 := math.MaxFloat64
	y4 := -math.MaxFloat64
	assert.True(t, num.CheckedSub(x4, y4).Equal(shepard.None[float64]()))

	x5 := 2.2
	y5 := 2.2
	assert.True(t, num.CheckedSub(x5, y5).Equal(shepard.Some[float64](0)))
}
```

## Mul

`num.CheckedMul[T num.Number](n T, v T) shepard.Option[T]`

### Examples

```go
func TestCheckedMul(t *testing.T) {
	x := int32(2147483327)
	y := int32(2147483327)

	assert.True(t, num.CheckedMul(x, y).Equal(shepard.None[int32]()))

	x2 := int32(-2147483327)
	y2 := int32(-2147483327)
	assert.True(t, num.CheckedMul(x2, y2).Equal(shepard.None[int32]()))

	x3 := int32(2)
	y3 := int32(2)
	assert.True(t, num.CheckedMul(x3, y3).Equal(shepard.Some[int32](4)))

	x4 := math.MaxFloat64
	y4 := math.MaxFloat64
	assert.True(t, num.CheckedMul(x4, y4).Equal(shepard.None[float64]()))

	x5 := 2.2
	y5 := 2.2
	assert.True(t, num.CheckedMul(x5, y5).Equal(shepard.Some[float64](4.4)))
}
```

## Div

`num.CheckedDiv[T num.Number](n T, v T) shepard.Option[T]`

```go
func TestCheckedDiv(t *testing.T) {
	x1 := math.MaxFloat64
	y1 := 0.5
	assert.True(t, num.CheckedDiv(x1, y1).Equal(shepard.None[float64]()))

	x2 := 4
	y2 := 2
	assert.True(t, num.CheckedDiv(x2, y2).Equal(shepard.Some(2)))
}
```