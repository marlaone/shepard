package num

import (
	"github.com/marlaone/shepard"
	"math"
)

// CheckedDiv divides two Number`s, checking for underflow, overflow and division by zero. If any of that happens, shepard.None is returned.
func CheckedDiv[T Number](num T, v T) shepard.Option[T] {
	if v == 0 {
		return shepard.None[T]()
	}
	c := num / v
	if math.IsInf(float64(c), 1) || math.IsInf(float64(c), 0) {
		return shepard.None[T]()
	}
	if (c < 0) == ((num < 0) != (v < 0)) {
		return shepard.Some(c)
	}
	return shepard.None[T]()
}
