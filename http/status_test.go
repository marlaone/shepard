package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCode_Int(t *testing.T) {
	tests := []struct {
		status   StatusCode
		expected int
	}{
		{StatusOK, 200},
		{StatusCreated, 201},
		{StatusNoContent, 204},
		{StatusBadRequest, 400},
		{StatusUnauthorized, 401},
		{StatusForbidden, 403},
		{StatusNotFound, 404},
		{StatusMethodNotAllowed, 405},
		{StatusInternalServerError, 500},
		{StatusBadGateway, 502},
		{StatusServiceUnavailable, 503},
	}

	for _, tt := range tests {
		t.Run(tt.status.Text(), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.Int())
		})
	}
}

func TestStatusCode_Text(t *testing.T) {
	tests := []struct {
		status   StatusCode
		expected string
	}{
		{StatusOK, "OK"},
		{StatusCreated, "Created"},
		{StatusNoContent, "No Content"},
		{StatusBadRequest, "Bad Request"},
		{StatusUnauthorized, "Unauthorized"},
		{StatusForbidden, "Forbidden"},
		{StatusNotFound, "Not Found"},
		{StatusMethodNotAllowed, "Method Not Allowed"},
		{StatusInternalServerError, "Internal Server Error"},
		{StatusBadGateway, "Bad Gateway"},
		{StatusServiceUnavailable, "Service Unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.Text())
		})
	}
}

func TestStatusCode_IsSuccess(t *testing.T) {
	success := []StatusCode{StatusOK, StatusCreated, StatusAccepted, StatusNoContent}
	notSuccess := []StatusCode{StatusBadRequest, StatusNotFound, StatusInternalServerError}

	for _, s := range success {
		t.Run(s.Text()+"_success", func(t *testing.T) {
			assert.True(t, s.IsSuccess())
		})
	}

	for _, s := range notSuccess {
		t.Run(s.Text()+"_not_success", func(t *testing.T) {
			assert.False(t, s.IsSuccess())
		})
	}
}

func TestStatusCode_IsClientError(t *testing.T) {
	clientErrors := []StatusCode{StatusBadRequest, StatusUnauthorized, StatusForbidden, StatusNotFound}
	notClientErrors := []StatusCode{StatusOK, StatusInternalServerError}

	for _, s := range clientErrors {
		t.Run(s.Text()+"_client_error", func(t *testing.T) {
			assert.True(t, s.IsClientError())
		})
	}

	for _, s := range notClientErrors {
		t.Run(s.Text()+"_not_client_error", func(t *testing.T) {
			assert.False(t, s.IsClientError())
		})
	}
}

func TestStatusCode_IsServerError(t *testing.T) {
	serverErrors := []StatusCode{StatusInternalServerError, StatusBadGateway, StatusServiceUnavailable}
	notServerErrors := []StatusCode{StatusOK, StatusBadRequest, StatusNotFound}

	for _, s := range serverErrors {
		t.Run(s.Text()+"_server_error", func(t *testing.T) {
			assert.True(t, s.IsServerError())
		})
	}

	for _, s := range notServerErrors {
		t.Run(s.Text()+"_not_server_error", func(t *testing.T) {
			assert.False(t, s.IsServerError())
		})
	}
}

func TestStatusCode_IsRedirect(t *testing.T) {
	redirects := []StatusCode{StatusMovedPermanently, StatusFound, StatusSeeOther, StatusTemporaryRedirect, StatusPermanentRedirect}
	notRedirects := []StatusCode{StatusOK, StatusBadRequest, StatusInternalServerError}

	for _, s := range redirects {
		t.Run(s.Text()+"_redirect", func(t *testing.T) {
			assert.True(t, s.IsRedirect())
		})
	}

	for _, s := range notRedirects {
		t.Run(s.Text()+"_not_redirect", func(t *testing.T) {
			assert.False(t, s.IsRedirect())
		})
	}
}
