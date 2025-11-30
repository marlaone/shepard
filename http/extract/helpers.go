package extract

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/http"
)

// --- Must* functions ---
// These functions extract values or panic on error.
// Use with http.Recover() middleware for automatic error handling.

// MustPath extracts a path parameter or panics with HttpError
// Use with Recover middleware for automatic 400 Bad Request responses
func MustPath[T any](ctx *http.Context, key string) T {
	result := ExtractPath[T](ctx, key)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// MustQuery extracts a query parameter or panics with HttpError
// Use with Recover middleware for automatic 400 Bad Request responses
func MustQuery[T any](ctx *http.Context, key string) T {
	result := ExtractQuery[T](ctx, key)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// MustHeader extracts a header value or panics with HttpError
// Use with Recover middleware for automatic 400 Bad Request responses
func MustHeader(ctx *http.Context, key string) string {
	result := ExtractHeader(ctx, key)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// MustJson extracts and parses JSON body or panics with HttpError
// Use with Recover middleware for automatic 400 Bad Request responses
func MustJson[T any](ctx *http.Context) T {
	result := Json[T](ctx)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// MustForm extracts URL-encoded form data or panics with HttpError
// Use with Recover middleware for automatic 400 Bad Request responses
func MustForm(ctx *http.Context) http.FormValues {
	result := Form(ctx)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// MustMultipart extracts multipart form data or panics with HttpError
// Use with Recover middleware for automatic 400 Bad Request responses
func MustMultipart(ctx *http.Context) http.MultipartForm {
	result := Multipart(ctx)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// MustState extracts typed state or panics
// Use with Recover middleware for automatic 500 Internal Server Error responses
func MustState[T any](ctx *http.Context, key string) T {
	result := ExtractState[T](ctx, key)
	if result.IsErr() {
		panic(result.UnwrapErr())
	}
	return result.Unwrap()
}

// --- *Or functions ---
// These functions extract values or return a default value on error.

// PathOr extracts a path parameter or returns the default value
func PathOr[T any](ctx *http.Context, key string, defaultVal T) T {
	result := ExtractPath[T](ctx, key)
	if result.IsErr() {
		return defaultVal
	}
	return result.Unwrap()
}

// QueryOr extracts a query parameter or returns the default value
func QueryOr[T any](ctx *http.Context, key string, defaultVal T) T {
	result := ExtractQuery[T](ctx, key)
	if result.IsErr() {
		return defaultVal
	}
	return result.Unwrap()
}

// HeaderOr extracts a header value or returns the default value
func HeaderOr(ctx *http.Context, key string, defaultVal string) string {
	result := ExtractHeader(ctx, key)
	if result.IsErr() {
		return defaultVal
	}
	return result.Unwrap()
}

// StateOr extracts typed state or returns the default value
func StateOr[T any](ctx *http.Context, key string, defaultVal T) T {
	result := ExtractState[T](ctx, key)
	if result.IsErr() {
		return defaultVal
	}
	return result.Unwrap()
}

// --- Try* functions ---
// These functions return (value, bool) like Go map access pattern

// TryPath extracts a path parameter, returning (value, true) on success or (zero, false) on failure
func TryPath[T any](ctx *http.Context, key string) (T, bool) {
	result := ExtractPath[T](ctx, key)
	if result.IsErr() {
		var zero T
		return zero, false
	}
	return result.Unwrap(), true
}

// TryQuery extracts a query parameter, returning (value, true) on success or (zero, false) on failure
func TryQuery[T any](ctx *http.Context, key string) (T, bool) {
	result := ExtractQuery[T](ctx, key)
	if result.IsErr() {
		var zero T
		return zero, false
	}
	return result.Unwrap(), true
}

// TryHeader extracts a header value, returning (value, true) on success or ("", false) on failure
func TryHeader(ctx *http.Context, key string) (string, bool) {
	result := ExtractHeader(ctx, key)
	if result.IsErr() {
		return "", false
	}
	return result.Unwrap(), true
}

// TryJson extracts and parses JSON body, returning (value, true) on success or (zero, false) on failure
func TryJson[T any](ctx *http.Context) (T, bool) {
	result := Json[T](ctx)
	if result.IsErr() {
		var zero T
		return zero, false
	}
	return result.Unwrap(), true
}

// TryState extracts typed state, returning (value, true) on success or (zero, false) on failure
func TryState[T any](ctx *http.Context, key string) (T, bool) {
	result := ExtractState[T](ctx, key)
	if result.IsErr() {
		var zero T
		return zero, false
	}
	return result.Unwrap(), true
}

// --- Option-based functions ---
// These functions return shepard.Option for Rust-style handling

// PathOption extracts a path parameter as Option
func PathOption[T any](ctx *http.Context, key string) shepard.Option[T] {
	return OptionalPath[T](ctx, key)
}

// QueryOption extracts a query parameter as Option
func QueryOption[T any](ctx *http.Context, key string) shepard.Option[T] {
	return OptionalQuery[T](ctx, key)
}

// HeaderOption extracts a header value as Option
func HeaderOption(ctx *http.Context, key string) shepard.Option[string] {
	return OptionalHeader(ctx, key)
}
