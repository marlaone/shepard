package shepard_json

import (
	"encoding/json"
	"github.com/marlaone/shepard"
)

type Result[T any, E any] shepard.Result[T, E]

func (r *Result[T, E]) MarshalJSON() ([]byte, error) {
	res := r.IntoResult()
	if res.IsOk() {
		return json.Marshal(res.Unwrap())
	}
	return nil, nil
}

func (r *Result[T, E]) UnmarshalJSON(b []byte) error {
	if string(b) == "{}" {
		var zero T
		*r = Result[T, E](*shepard.Ok[T, E](zero))
		return nil
	}
	var s T
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*r = Result[T, E](*shepard.Ok[T, E](s))
	return nil
}

func (r *Result[T, E]) IntoResult() *shepard.Result[T, E] {
	if r == nil {
		var zero T
		return shepard.Ok[T, E](zero)
	}
	res := shepard.Result[T, E](*r)
	return &res
}

func ParseResult[T, E any](res *shepard.Result[T, E]) *Result[T, E] {
	r := Result[T, E](*res)
	return &r
}
