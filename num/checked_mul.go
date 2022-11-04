package num

import (
	"github.com/marlaone/shepard"
	"math"
)

// CheckedMul multiplies two Number`s, checking for underflow or overflow. If underflow or overflow happens, shepard.None is returned.
func CheckedMul[T Number](num T, v T) shepard.Option[T] {
	if num == 0 || v == 0 {
		return shepard.Some[T](num * v)
	}
	c := num * v
	if math.IsInf(float64(c), 1) || math.IsInf(float64(c), 0) {
		return shepard.None[T]()
	}
	if (c < 0) == ((num < 0) != (v < 0)) {
		if c/v == num {
			return shepard.Some[T](c)
		}
	}
	return shepard.None[T]()
}
