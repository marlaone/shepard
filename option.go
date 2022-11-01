package shepard

type OptionUnwrapElseFunc[T comparable] func() T
type OptionOkOrElseFunc func() error
type OptionAndThenFunc[T comparable] func(val *T) *Option[T]
type OptionFilterFunc[T comparable] func(val *T) bool
type OptionOrElseFunc[T comparable] func() *Option[T]
type OptionGetOrInsertWithFunc[T comparable] func() T

type Option[T comparable] struct {
	some *T
}

func Some[T comparable](val T) *Option[T] {
	return &Option[T]{
		some: &val,
	}
}

func None[T comparable]() *Option[T] {
	return &Option[T]{
		some: nil,
	}
}

func (o *Option[T]) Unwrap() T {
	return *o.some
}

func (o *Option[T]) IsSome() bool {
	return o.some != nil
}

func (o *Option[T]) IsNone() bool {
	return o.some == nil
}

func (o *Option[T]) UnwrapOr(value T) T {
	if o.IsSome() {
		return o.Unwrap()
	}
	return value
}

func (o *Option[T]) UnwrapOrElse(elseFunc OptionUnwrapElseFunc[T]) T {
	if o.IsSome() {
		return o.Unwrap()
	}
	return elseFunc()
}

func (o *Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		var zero T
		return zero
	}
	return o.Unwrap()
}

func (o *Option[T]) OkOr(err error) *Result[T] {
	if o.IsSome() {
		return Ok(o.Unwrap())
	}
	return Err[T](err)
}

func (o *Option[T]) OkOrElse(op OptionOkOrElseFunc) *Result[T] {
	if o.IsSome() {
		return Ok(o.Unwrap())
	}
	return Err[T](op())
}

func (o *Option[T]) And(opt *Option[T]) *Option[T] {
	if o.IsSome() && opt.IsSome() {
		return opt
	}
	return None[T]()
}

func (o *Option[T]) AndThen(op OptionAndThenFunc[T]) *Option[T] {
	if o.IsSome() {
		return op(o.some)
	}
	return o
}

func (o *Option[T]) Filter(predicate OptionFilterFunc[T]) *Option[T] {
	if o.IsSome() && predicate(o.some) {
		return o
	}
	return None[T]()
}

func (o *Option[T]) Or(opt *Option[T]) *Option[T] {
	if o.IsSome() {
		return o
	}
	if opt.IsSome() {
		return opt
	}
	return None[T]()
}

func (o *Option[T]) OrElse(op OptionOrElseFunc[T]) *Option[T] {
	if o.IsSome() {
		return o
	}
	return op()
}

func (o *Option[T]) Xor(opt *Option[T]) *Option[T] {
	if o.IsSome() && opt.IsSome() {
		return None[T]()
	}
	if o.IsSome() {
		return o
	}
	if opt.IsSome() {
		return opt
	}
	return None[T]()
}

func (o *Option[T]) Insert(val T) *T {
	o.some = &val
	return o.some
}

func (o *Option[T]) GetOrInsert(val T) *T {
	if o.IsSome() {
		return o.some
	}
	o.some = &val
	return o.some
}

func (o *Option[T]) GetOrInsertDefault() *T {
	if o.IsSome() {
		return o.some
	}
	val := new(T)
	o.some = val
	return val
}

func (o *Option[T]) GetOrInsertWith(f OptionGetOrInsertWithFunc[T]) *T {
	if o.IsSome() {
		return o.some
	}
	val := f()
	o.some = &val
	return o.some
}

func (o *Option[T]) Take() T {
	val := *o.some
	o.some = nil
	return val
}

func (o *Option[T]) Replace(value T) *Option[T] {
	old := *o
	o.some = &value
	return &old
}

func (o *Option[T]) Equal(opt *Option[T]) bool {
	if o.IsSome() && opt.IsSome() {
		return o.Unwrap() == opt.Unwrap()
	}
	return o.IsNone() == opt.IsNone()
}
