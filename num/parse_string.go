package num

import (
	"fmt"
	"github.com/marlaone/shepard"
	"reflect"
	"strconv"
)

func ParseString[T Number](n string) shepard.Result[T, error] {
	var target T
	numType := reflect.TypeOf(target)
	switch numType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(n, 10, numType.Bits())
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(n, numType.Bits())
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(n, 10, numType.Bits())
		if err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Ok[T, error](T(v))
	}
	return shepard.Err[T, error](fmt.Errorf("unknown numeric type: %T", target))
}
