package shepard_json

import (
	"encoding/json"
	"github.com/marlaone/shepard"
)

type Option[T any] shepard.Option[T]

func (o *Option[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "{}" {
		opt := shepard.None[T]()
		*o = Option[T](opt)
		return nil
	}
	var s T
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	opt := o.IntoOption()
	opt.Replace(s)
	*o = Option[T](opt)
	return nil
}

func (o *Option[T]) MarshalJSON() ([]byte, error) {
	opt := o.IntoOption()
	if opt.IsSome() {
		return json.Marshal(opt.Unwrap())
	}
	return nil, nil
}

func (o *Option[T]) IntoOption() shepard.Option[T] {
	if o == nil {
		return shepard.None[T]()
	}
	opt := shepard.Option[T](*o)
	return opt
}

func ParseOption[T any](opt shepard.Option[T]) *Option[T] {
	o := Option[T](opt)
	return &o
}
