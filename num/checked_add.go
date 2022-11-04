package num

import (
	"github.com/marlaone/shepard"
	"math"
)

// CheckedAdd adds two Number`s, checking for overflow. If overflow happens, shepard.None is returned.
func CheckedAdd[T Number](num T, v T) shepard.Option[T] {
	c := num + v
	if math.IsInf(float64(c), 1) || math.IsInf(float64(c), 0) {
		return shepard.None[T]()
	}
	if (c > num) == (v > 0) {
		return shepard.Some(c)
	}
	return shepard.None[T]()
}
