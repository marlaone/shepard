package num

import (
	"strconv"
	"unsafe"

	"github.com/marlaone/shepard"
)

// ParseFloat64 is a specialized fast path for parsing float64 from string.
// Avoids the type switch overhead of the generic ParseString.
func ParseFloat64(s string) shepard.Result[float64, error] {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return shepard.Err[float64, error](err)
	}
	return shepard.Ok[float64, error](v)
}

// ParseFloat64Bytes is a specialized fast path for parsing float64 from bytes.
// Combines zero-copy conversion with no type switch overhead.
// This is the optimal choice for hot paths parsing float64 from byte slices.
func ParseFloat64Bytes(b []byte) shepard.Result[float64, error] {
	// Zero-copy conversion: safe because strconv.ParseFloat
	// does not retain the string beyond the function call
	s := unsafe.String(unsafe.SliceData(b), len(b))
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return shepard.Err[float64, error](err)
	}
	return shepard.Ok[float64, error](v)
}

// ParseFloat32 is a specialized fast path for parsing float32 from string.
func ParseFloat32(s string) shepard.Result[float32, error] {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return shepard.Err[float32, error](err)
	}
	return shepard.Ok[float32, error](float32(v))
}

// ParseFloat32Bytes is a specialized fast path for parsing float32 from bytes.
func ParseFloat32Bytes(b []byte) shepard.Result[float32, error] {
	s := unsafe.String(unsafe.SliceData(b), len(b))
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return shepard.Err[float32, error](err)
	}
	return shepard.Ok[float32, error](float32(v))
}
