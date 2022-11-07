package num

import (
	"errors"
	"github.com/marlaone/shepard/iter"
)

// Sum sums the elements of an iterator.
// Takes each element, adds them together, and returns the result.
// An empty iterator returns the zero value of the type.
//
// Panics when calling Sum and a primitive integer type is being returned, this method will panic if the computation overflows.
func Sum[T Number](iterator iter.Iter[T]) T {
	var sum T

	for {
		v := iterator.Next()
		if v.IsNone() {
			break
		}
		checked := CheckedAdd[T](sum, v.Unwrap())
		if checked.IsNone() {
			panic(errors.New("computation overflow"))
		}
		sum = checked.Unwrap()
	}
	return sum
}
