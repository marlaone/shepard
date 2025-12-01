// Package extract provides type-safe extractors for HTTP request data
// inspired by axum's extractor pattern.
package extract

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/http"
)

// Extractor is the base interface for all extractors
type Extractor[T any] interface {
	Extract(ctx *http.Context, key string) shepard.Result[T, http.HttpError]
}

// --- Path Parameter Extractor ---

// Path extracts a path parameter and parses it to type T
type Path[T any] struct{}

// Extract extracts and parses a path parameter
func (p Path[T]) Extract(ctx *http.Context, key string) shepard.Result[T, http.HttpError] {
	paramOpt := ctx.Param(key)
	if paramOpt.IsNone() {
		return shepard.Err[T](
			http.BadRequest(fmt.Sprintf("missing path parameter: %s", key)),
		)
	}

	return parseValue[T](paramOpt.Unwrap(), fmt.Sprintf("path parameter %s", key))
}

// ExtractPath is a convenience function for extracting path parameters
func ExtractPath[T any](ctx *http.Context, key string) shepard.Result[T, http.HttpError] {
	var p Path[T]
	return p.Extract(ctx, key)
}

// OptionalPath extracts a path parameter as Option, returning None if not found
func OptionalPath[T any](ctx *http.Context, key string) shepard.Option[T] {
	result := ExtractPath[T](ctx, key)
	if result.IsErr() {
		return shepard.None[T]()
	}
	return shepard.Some(result.Unwrap())
}

// --- Query Parameter Extractor ---

// Query extracts a query parameter and parses it to type T
type Query[T any] struct{}

// Extract extracts and parses a query parameter
func (q Query[T]) Extract(ctx *http.Context, key string) shepard.Result[T, http.HttpError] {
	queryOpt := ctx.Query(key)
	if queryOpt.IsNone() {
		return shepard.Err[T](
			http.BadRequest(fmt.Sprintf("missing query parameter: %s", key)),
		)
	}

	return parseValue[T](queryOpt.Unwrap(), fmt.Sprintf("query parameter %s", key))
}

// ExtractQuery is a convenience function for extracting query parameters
func ExtractQuery[T any](ctx *http.Context, key string) shepard.Result[T, http.HttpError] {
	var q Query[T]
	return q.Extract(ctx, key)
}

// OptionalQuery extracts a query parameter as Option, returning None if not found
func OptionalQuery[T any](ctx *http.Context, key string) shepard.Option[T] {
	result := ExtractQuery[T](ctx, key)
	if result.IsErr() {
		return shepard.None[T]()
	}
	return shepard.Some(result.Unwrap())
}

// --- Header Extractor ---

// Header extracts a header value as a string
type Header struct{}

// Extract extracts a header value
func (h Header) Extract(ctx *http.Context, key string) shepard.Result[string, http.HttpError] {
	headerOpt := ctx.Header(key)
	if headerOpt.IsNone() {
		return shepard.Err[string](
			http.BadRequest(fmt.Sprintf("missing header: %s", key)),
		)
	}

	return shepard.Ok[string, http.HttpError](headerOpt.Unwrap())
}

// ExtractHeader is a convenience function for extracting headers
func ExtractHeader(ctx *http.Context, key string) shepard.Result[string, http.HttpError] {
	var h Header
	return h.Extract(ctx, key)
}

// OptionalHeader extracts a header as Option, returning None if not found
func OptionalHeader(ctx *http.Context, key string) shepard.Option[string] {
	result := ExtractHeader(ctx, key)
	if result.IsErr() {
		return shepard.None[string]()
	}
	return shepard.Some(result.Unwrap())
}

// --- JSON Body Extractor ---

// Json extracts the request body as JSON and deserializes it to type T
func Json[T any](ctx *http.Context) shepard.Result[T, http.HttpError] {
	result := http.ContextBodyJson[T](ctx)
	if result.IsErr() {
		return shepard.Err[T](
			http.BadRequest(fmt.Sprintf("failed to parse JSON body: %v", result.UnwrapErr())),
		)
	}
	return shepard.Ok[T, http.HttpError](result.Unwrap())
}

// --- Form Body Extractor ---

// Form extracts URL-encoded form data (application/x-www-form-urlencoded) from the request body
func Form(ctx *http.Context) shepard.Result[http.FormValues, http.HttpError] {
	result := http.ContextBodyForm(ctx)
	if result.IsErr() {
		return shepard.Err[http.FormValues](
			http.BadRequest(fmt.Sprintf("failed to parse form body: %v", result.UnwrapErr())),
		)
	}
	return shepard.Ok[http.FormValues, http.HttpError](result.Unwrap())
}

// --- Multipart Form Extractor ---

// Multipart extracts multipart form data (multipart/form-data) from the request body
func Multipart(ctx *http.Context) shepard.Result[http.MultipartForm, http.HttpError] {
	result := http.ContextBodyMultipart(ctx)
	if result.IsErr() {
		return shepard.Err[http.MultipartForm](
			http.BadRequest(fmt.Sprintf("failed to parse multipart form: %v", result.UnwrapErr())),
		)
	}
	return shepard.Ok[http.MultipartForm, http.HttpError](result.Unwrap())
}

// MultipartWithLimit extracts multipart form data with a custom memory limit
func MultipartWithLimit(maxMemory int64) func(*http.Context) shepard.Result[http.MultipartForm, http.HttpError] {
	return func(ctx *http.Context) shepard.Result[http.MultipartForm, http.HttpError] {
		result := http.ContextBodyMultipartWithLimit(ctx, maxMemory)
		if result.IsErr() {
			return shepard.Err[http.MultipartForm](
				http.BadRequest(fmt.Sprintf("failed to parse multipart form: %v", result.UnwrapErr())),
			)
		}
		return shepard.Ok[http.MultipartForm, http.HttpError](result.Unwrap())
	}
}

// --- State Extractor ---

// State extracts a value from the context state
type State[T any] struct{}

// Extract extracts a state value and type-asserts it to T
func (s State[T]) Extract(ctx *http.Context, key string) shepard.Result[T, http.HttpError] {
	stateOpt := ctx.GetState(key)
	if stateOpt.IsNone() {
		return shepard.Err[T](
			http.InternalError(fmt.Errorf("state not found: %s", key)),
		)
	}

	value, ok := stateOpt.Unwrap().(T)
	if !ok {
		return shepard.Err[T](
			http.InternalError(fmt.Errorf("state type mismatch for key: %s", key)),
		)
	}

	return shepard.Ok[T, http.HttpError](value)
}

// ExtractState is a convenience function for extracting state
func ExtractState[T any](ctx *http.Context, key string) shepard.Result[T, http.HttpError] {
	var s State[T]
	return s.Extract(ctx, key)
}

// --- Method Extractor ---

// ExtractMethod returns the HTTP method from the request
func ExtractMethod(ctx *http.Context) http.Method {
	return ctx.Method()
}

// --- Custom Extractor Registry ---

// CustomExtractorFn is a function type for custom extractors
type CustomExtractorFn func(ctx *http.Context) shepard.Result[any, http.HttpError]

var (
	customExtractors = make(map[string]CustomExtractorFn)
	extractorsMu     sync.RWMutex
)

// RegisterExtractor registers a custom extractor function
func RegisterExtractor(name string, fn CustomExtractorFn) {
	extractorsMu.Lock()
	defer extractorsMu.Unlock()
	customExtractors[name] = fn
}

// ExtractCustom extracts a value using a registered custom extractor
func ExtractCustom[T any](ctx *http.Context, name string) shepard.Result[T, http.HttpError] {
	extractorsMu.RLock()
	fn, ok := customExtractors[name]
	extractorsMu.RUnlock()

	if !ok {
		return shepard.Err[T](
			http.InternalError(fmt.Errorf("custom extractor not found: %s", name)),
		)
	}

	result := fn(ctx)
	if result.IsErr() {
		return shepard.Err[T](result.UnwrapErr())
	}

	value, ok := result.Unwrap().(T)
	if !ok {
		return shepard.Err[T](
			http.InternalError(fmt.Errorf("custom extractor type mismatch: %s", name)),
		)
	}

	return shepard.Ok[T, http.HttpError](value)
}

// --- Value Parsing ---

// parseValue parses a string value into type T
func parseValue[T any](value string, context string) shepard.Result[T, http.HttpError] {
	var result T
	var err error

	// Use type switch on the zero value to determine the target type
	switch any(result).(type) {
	case string:
		result = any(value).(T)
		return shepard.Ok[T, http.HttpError](result)

	case int:
		var i int
		i, err = strconv.Atoi(value)
		if err == nil {
			result = any(i).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case int64:
		var i int64
		i, err = strconv.ParseInt(value, 10, 64)
		if err == nil {
			result = any(i).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case int32:
		var i int64
		i, err = strconv.ParseInt(value, 10, 32)
		if err == nil {
			result = any(int32(i)).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case uint:
		var u uint64
		u, err = strconv.ParseUint(value, 10, 64)
		if err == nil {
			result = any(uint(u)).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case uint64:
		var u uint64
		u, err = strconv.ParseUint(value, 10, 64)
		if err == nil {
			result = any(u).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case float64:
		var f float64
		f, err = strconv.ParseFloat(value, 64)
		if err == nil {
			result = any(f).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case float32:
		var f float64
		f, err = strconv.ParseFloat(value, 32)
		if err == nil {
			result = any(float32(f)).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	case bool:
		var b bool
		lowerVal := strings.ToLower(value)
		switch lowerVal {
		case "true", "1", "yes", "on":
			b = true
		case "false", "0", "no", "off":
			b = false
		default:
			err = fmt.Errorf("invalid boolean value: %s", value)
		}
		if err == nil {
			result = any(b).(T)
			return shepard.Ok[T, http.HttpError](result)
		}

	default:
		return shepard.Err[T](
			http.BadRequest(fmt.Sprintf("unsupported type for %s", context)),
		)
	}

	return shepard.Err[T](
		http.BadRequest(fmt.Sprintf("failed to parse %s: %v", context, err)),
	)
}
