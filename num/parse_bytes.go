package num

import (
	"unsafe"

	"github.com/marlaone/shepard"
)

// ParseBytes parses a byte slice as a number without allocating a string.
// Uses unsafe.String for zero-copy conversion - this is safe because strconv
// functions do not retain the string beyond the call.
func ParseBytes[T Number](b []byte) shepard.Result[T, error] {
	// Zero-copy conversion: safe because strconv.Parse* functions
	// do not retain the string beyond the function call
	s := unsafe.String(unsafe.SliceData(b), len(b))
	return ParseString[T](s)
}
