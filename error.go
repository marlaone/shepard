package shepard

type Error[T any] struct {
	val *T
}

func NewError[T any](val T) *Error[T] {
	return &Error[T]{
		val: &val,
	}
}

func (err *Error[T]) Value() T {
	return *err.val
}
