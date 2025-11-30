package num

import (
	"fmt"
	"strconv"

	"github.com/marlaone/shepard"
)

func ParseString[T Number](n string) shepard.Result[T, error] {
	var target T
	// Use type switch instead of reflection for better performance
	switch any(target).(type) {
	case int:
		v, err := strconv.ParseInt(n, 10, 0)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case int8:
		v, err := strconv.ParseInt(n, 10, 8)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case int16:
		v, err := strconv.ParseInt(n, 10, 16)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case int32:
		v, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case int64:
		v, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case uint:
		v, err := strconv.ParseUint(n, 10, 0)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case uint8:
		v, err := strconv.ParseUint(n, 10, 8)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case uint16:
		v, err := strconv.ParseUint(n, 10, 16)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case uint32:
		v, err := strconv.ParseUint(n, 10, 32)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case uint64:
		v, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case float32:
		v, err := strconv.ParseFloat(n, 32)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case float64:
		v, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	}
	return shepard.Err[T, error](fmt.Errorf("unknown numeric type: %T", target))
}
