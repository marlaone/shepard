package http

import (
	"errors"
	"testing"

	"github.com/marlaone/shepard"
	"github.com/stretchr/testify/assert"
)

func TestHttpError_Error(t *testing.T) {
	err := HttpError{
		StatusCode: StatusBadRequest,
		Message:    "invalid input",
		Cause:      shepard.None[error](),
	}

	assert.Equal(t, "400 Bad Request: invalid input", err.Error())
}

func TestHttpError_ErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := HttpError{
		StatusCode: StatusInternalServerError,
		Message:    "something went wrong",
		Cause:      shepard.Some(cause),
	}

	errStr := err.Error()
	assert.Contains(t, errStr, "500 Internal Server Error")
	assert.Contains(t, errStr, "something went wrong")
	assert.Contains(t, errStr, "underlying error")
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("missing field")

	assert.Equal(t, StatusBadRequest, err.StatusCode)
	assert.Equal(t, "missing field", err.Message)
	assert.True(t, err.Cause.IsNone())
}

func TestNotFound(t *testing.T) {
	err := NotFound("user not found")

	assert.Equal(t, StatusNotFound, err.StatusCode)
	assert.Equal(t, "user not found", err.Message)
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("invalid token")

	assert.Equal(t, StatusUnauthorized, err.StatusCode)
	assert.Equal(t, "invalid token", err.Message)
}

func TestForbidden(t *testing.T) {
	err := Forbidden("access denied")

	assert.Equal(t, StatusForbidden, err.StatusCode)
	assert.Equal(t, "access denied", err.Message)
}

func TestInternalError(t *testing.T) {
	cause := errors.New("db connection failed")
	err := InternalError(cause)

	assert.Equal(t, StatusInternalServerError, err.StatusCode)
	assert.Equal(t, "internal server error", err.Message)
	assert.True(t, err.Cause.IsSome())
	assert.Equal(t, cause, err.Cause.Unwrap())
}

func TestHttpError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := HttpError{
		StatusCode: StatusBadRequest,
		Message:    "bad",
		Cause:      shepard.Some(cause),
	}

	// Should work with errors.Unwrap
	unwrapped := errors.Unwrap(err)
	assert.Equal(t, cause, unwrapped)
}

func TestHttpError_Is(t *testing.T) {
	err1 := BadRequest("test")
	err2 := BadRequest("different message")
	err3 := NotFound("test")

	// Same status code should match
	assert.True(t, errors.Is(err1, err2))

	// Different status codes should not match
	assert.False(t, errors.Is(err1, err3))
}

func TestRejection(t *testing.T) {
	rejection := Rejection{
		StatusCode: StatusUnprocessableEntity,
		Message:    "validation failed",
	}

	assert.Equal(t, "422 Unprocessable Entity: validation failed", rejection.Error())
}

// --- Tests for Result-returning error helpers ---

func TestHttpError_Err(t *testing.T) {
	httpErr := BadRequest("invalid")
	result := httpErr.Err()

	assert.True(t, result.IsErr())
	assert.Equal(t, StatusBadRequest, result.UnwrapErr().StatusCode)
}

func TestBadRequestErr(t *testing.T) {
	result := BadRequestErr("missing field")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusBadRequest, err.StatusCode)
	assert.Equal(t, "missing field", err.Message)
}

func TestUnauthorizedErr(t *testing.T) {
	result := UnauthorizedErr("invalid token")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusUnauthorized, err.StatusCode)
	assert.Equal(t, "invalid token", err.Message)
}

func TestForbiddenErr(t *testing.T) {
	result := ForbiddenErr("access denied")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusForbidden, err.StatusCode)
	assert.Equal(t, "access denied", err.Message)
}

func TestNotFoundErr(t *testing.T) {
	result := NotFoundErr("user not found")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusNotFound, err.StatusCode)
	assert.Equal(t, "user not found", err.Message)
}

func TestMethodNotAllowedErr(t *testing.T) {
	result := MethodNotAllowedErr("POST not allowed")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusMethodNotAllowed, err.StatusCode)
}

func TestConflictErr(t *testing.T) {
	result := ConflictErr("resource conflict")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusConflict, err.StatusCode)
}

func TestUnprocessableEntityErr(t *testing.T) {
	result := UnprocessableEntityErr("validation failed")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusUnprocessableEntity, err.StatusCode)
}

func TestInternalErr(t *testing.T) {
	cause := errors.New("db error")
	result := InternalErr(cause)

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusInternalServerError, err.StatusCode)
	assert.True(t, err.Cause.IsSome())
	assert.Equal(t, cause, err.Cause.Unwrap())
}

func TestServiceUnavailableErr(t *testing.T) {
	result := ServiceUnavailableErr("service down")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, StatusServiceUnavailable, err.StatusCode)
}
