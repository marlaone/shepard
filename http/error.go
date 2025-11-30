package http

import (
	"errors"
	"fmt"

	"github.com/marlaone/shepard"
)

// Common errors
var (
	ErrInvalidMethod = errors.New("invalid HTTP method")
)

// HttpError represents an HTTP-specific error with status code
type HttpError struct {
	StatusCode StatusCode
	Message    string
	Cause      shepard.Option[error]
}

// Error implements the error interface
func (e HttpError) Error() string {
	base := fmt.Sprintf("%d %s: %s", e.StatusCode.Int(), e.StatusCode.Text(), e.Message)
	if e.Cause.IsSome() {
		return fmt.Sprintf("%s (caused by: %v)", base, e.Cause.Unwrap())
	}
	return base
}

// Unwrap returns the underlying cause for use with errors.Unwrap
func (e HttpError) Unwrap() error {
	if e.Cause.IsSome() {
		return e.Cause.Unwrap()
	}
	return nil
}

// Is implements errors.Is for comparing HttpErrors by status code
func (e HttpError) Is(target error) bool {
	var httpErr HttpError
	if errors.As(target, &httpErr) {
		return e.StatusCode == httpErr.StatusCode
	}
	return false
}

// BadRequest creates a 400 Bad Request error
func BadRequest(msg string) HttpError {
	return HttpError{
		StatusCode: StatusBadRequest,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(msg string) HttpError {
	return HttpError{
		StatusCode: StatusUnauthorized,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// Forbidden creates a 403 Forbidden error
func Forbidden(msg string) HttpError {
	return HttpError{
		StatusCode: StatusForbidden,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// NotFound creates a 404 Not Found error
func NotFound(msg string) HttpError {
	return HttpError{
		StatusCode: StatusNotFound,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// MethodNotAllowed creates a 405 Method Not Allowed error
func MethodNotAllowed(msg string) HttpError {
	return HttpError{
		StatusCode: StatusMethodNotAllowed,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// Conflict creates a 409 Conflict error
func Conflict(msg string) HttpError {
	return HttpError{
		StatusCode: StatusConflict,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// UnprocessableEntity creates a 422 Unprocessable Entity error
func UnprocessableEntity(msg string) HttpError {
	return HttpError{
		StatusCode: StatusUnprocessableEntity,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// InternalError creates a 500 Internal Server Error with a cause
func InternalError(cause error) HttpError {
	return HttpError{
		StatusCode: StatusInternalServerError,
		Message:    "internal server error",
		Cause:      shepard.Some(cause),
	}
}

// ServiceUnavailable creates a 503 Service Unavailable error
func ServiceUnavailable(msg string) HttpError {
	return HttpError{
		StatusCode: StatusServiceUnavailable,
		Message:    msg,
		Cause:      shepard.None[error](),
	}
}

// Rejection represents an extractor rejection (used by FromRequest extractors)
type Rejection struct {
	StatusCode StatusCode
	Message    string
}

// Error implements the error interface
func (r Rejection) Error() string {
	return fmt.Sprintf("%d %s: %s", r.StatusCode.Int(), r.StatusCode.Text(), r.Message)
}

// IntoHttpError converts a Rejection into an HttpError
func (r Rejection) IntoHttpError() HttpError {
	return HttpError{
		StatusCode: r.StatusCode,
		Message:    r.Message,
		Cause:      shepard.None[error](),
	}
}

// --- Result-returning error helpers ---
// These functions return shepard.Result[Response, HttpError] for use directly in handlers,
// eliminating the need for verbose shepard.Err[Response, HttpError](...) wrapping.

// Err converts an HttpError to a Result[Response, HttpError] for returning from handlers
func (e HttpError) Err() shepard.Result[Response, HttpError] {
	return shepard.Err[Response, HttpError](e)
}

// BadRequestErr returns a 400 Bad Request error Result
func BadRequestErr(msg string) shepard.Result[Response, HttpError] {
	return BadRequest(msg).Err()
}

// UnauthorizedErr returns a 401 Unauthorized error Result
func UnauthorizedErr(msg string) shepard.Result[Response, HttpError] {
	return Unauthorized(msg).Err()
}

// ForbiddenErr returns a 403 Forbidden error Result
func ForbiddenErr(msg string) shepard.Result[Response, HttpError] {
	return Forbidden(msg).Err()
}

// NotFoundErr returns a 404 Not Found error Result
func NotFoundErr(msg string) shepard.Result[Response, HttpError] {
	return NotFound(msg).Err()
}

// MethodNotAllowedErr returns a 405 Method Not Allowed error Result
func MethodNotAllowedErr(msg string) shepard.Result[Response, HttpError] {
	return MethodNotAllowed(msg).Err()
}

// ConflictErr returns a 409 Conflict error Result
func ConflictErr(msg string) shepard.Result[Response, HttpError] {
	return Conflict(msg).Err()
}

// UnprocessableEntityErr returns a 422 Unprocessable Entity error Result
func UnprocessableEntityErr(msg string) shepard.Result[Response, HttpError] {
	return UnprocessableEntity(msg).Err()
}

// InternalErr returns a 500 Internal Server Error Result with a cause
func InternalErr(cause error) shepard.Result[Response, HttpError] {
	return InternalError(cause).Err()
}

// ServiceUnavailableErr returns a 503 Service Unavailable error Result
func ServiceUnavailableErr(msg string) shepard.Result[Response, HttpError] {
	return ServiceUnavailable(msg).Err()
}
