package shepard

type Default[T any] interface {
	Default() T
}

func GetDefault[T any]() T {
	var valType T
	defaulter, ok := any(valType).(Default[T])
	if ok {
		return defaulter.Default()
	}
	return valType
}
