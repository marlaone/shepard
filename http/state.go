package http

import (
	"fmt"
	"reflect"

	"github.com/marlaone/shepard"
)

// --- Typed State Functions ---
// These functions provide type-safe state storage using the type itself as a key,
// eliminating the need for string keys and type assertions.

// SetTypedState stores a value in context state using its type as the key
// Multiple values of the same type will overwrite each other
func SetTypedState[T any](ctx *Context, value T) {
	key := typeKey[T]()
	ctx.SetState(key, value)
}

// GetTypedState retrieves a value from context state by its type
// Returns None if not found or if type assertion fails
func GetTypedState[T any](ctx *Context) shepard.Option[T] {
	key := typeKey[T]()
	result := ctx.GetState(key)
	if result.IsNone() {
		return shepard.None[T]()
	}
	value, ok := result.Unwrap().(T)
	if !ok {
		return shepard.None[T]()
	}
	return shepard.Some(value)
}

// MustTypedState retrieves typed state or panics if not found
// Use with Recover middleware for automatic error handling
func MustTypedState[T any](ctx *Context) T {
	opt := GetTypedState[T](ctx)
	if opt.IsNone() {
		panic(fmt.Sprintf("typed state not found for type: %s", typeKey[T]()))
	}
	return opt.Unwrap()
}

// TypedStateOr retrieves typed state or returns the default value if not found
func TypedStateOr[T any](ctx *Context, defaultVal T) T {
	opt := GetTypedState[T](ctx)
	if opt.IsNone() {
		return defaultVal
	}
	return opt.Unwrap()
}

// HasTypedState returns true if typed state exists for the given type
func HasTypedState[T any](ctx *Context) bool {
	key := typeKey[T]()
	return ctx.GetState(key).IsSome()
}

// RemoveTypedState removes typed state for the given type
func RemoveTypedState[T any](ctx *Context) {
	key := typeKey[T]()
	ctx.state.Remove(key)
}

// typeKey generates a unique key for a type
func typeKey[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}

// --- WithState Middleware ---
// Creates middleware that injects typed state into every request context

// WithState creates middleware that injects a value into context state using its type as key
// This allows handlers to retrieve the value with GetTypedState[T](ctx)
func WithState[T any](value T) Middleware {
	return func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		SetTypedState(ctx, value)
		return next()
	}
}

// --- Named State (for multiple values of same type) ---
// When you need multiple values of the same type, use named state

// SetNamedState stores a value with both type and name as key
func SetNamedState[T any](ctx *Context, name string, value T) {
	key := namedTypeKey[T](name)
	ctx.SetState(key, value)
}

// GetNamedState retrieves a named value by type and name
func GetNamedState[T any](ctx *Context, name string) shepard.Option[T] {
	key := namedTypeKey[T](name)
	result := ctx.GetState(key)
	if result.IsNone() {
		return shepard.None[T]()
	}
	value, ok := result.Unwrap().(T)
	if !ok {
		return shepard.None[T]()
	}
	return shepard.Some(value)
}

// MustNamedState retrieves named state or panics if not found
func MustNamedState[T any](ctx *Context, name string) T {
	opt := GetNamedState[T](ctx, name)
	if opt.IsNone() {
		panic(fmt.Sprintf("named state '%s' not found for type: %s", name, typeKey[T]()))
	}
	return opt.Unwrap()
}

// NamedStateOr retrieves named state or returns the default value
func NamedStateOr[T any](ctx *Context, name string, defaultVal T) T {
	opt := GetNamedState[T](ctx, name)
	if opt.IsNone() {
		return defaultVal
	}
	return opt.Unwrap()
}

// WithNamedState creates middleware that injects a named value
func WithNamedState[T any](name string, value T) Middleware {
	return func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError] {
		SetNamedState(ctx, name, value)
		return next()
	}
}

// namedTypeKey generates a unique key for a named typed value
func namedTypeKey[T any](name string) string {
	return fmt.Sprintf("%s:%s", name, typeKey[T]())
}
