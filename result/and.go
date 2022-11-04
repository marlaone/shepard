package result

import "github.com/marlaone/shepard"

type AndThenFunc[T any, U any, E any] func(v T) shepard.Result[U, E]

// And returns res if the result is shepard.Ok, otherwise returns the shepard.Err value of self.
//
// Arguments passed to and are eagerly evaluated; if you are passing the result of a function call, it is recommended to use AndThen, which is lazily evaluated.
func And[T any, U any, E any](r shepard.Result[T, E], resb shepard.Result[U, E]) shepard.Result[U, E] {
	if r.IsOk() && resb.IsOk() {
		return resb
	}
	if r.IsErr() {
		return shepard.Err[U](r.Err().Unwrap())
	}
	return shepard.Err[U](resb.Err().Unwrap())
}

// AndThen calls op if the result is shepard.Ok, otherwise returns the shepard.Err value of self.
//
// This function can be used for control flow based on shepard.Result values.
func AndThen[T any, U any, E any](r shepard.Result[T, E], op AndThenFunc[T, U, E]) shepard.Result[U, E] {
	if r.IsErr() {
		return shepard.Err[U, E](r.Err().Unwrap())
	}
	return op(r.Unwrap())
}
