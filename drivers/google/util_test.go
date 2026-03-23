package google

import (
	"net/http"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestIsTransientError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		transient bool
	}{
		{
			name:      "nil error",
			err:       nil,
			transient: false,
		},
		{
			name:      "503 Service Unavailable",
			err:       &googleapi.Error{Code: http.StatusServiceUnavailable},
			transient: true,
		},
		{
			name:      "429 Too Many Requests",
			err:       &googleapi.Error{Code: http.StatusTooManyRequests},
			transient: true,
		},
		{
			name:      "502 Bad Gateway",
			err:       &googleapi.Error{Code: http.StatusBadGateway},
			transient: true,
		},
		{
			name:      "500 Internal Server Error",
			err:       &googleapi.Error{Code: http.StatusInternalServerError},
			transient: true,
		},
		{
			name:      "404 Not Found",
			err:       &googleapi.Error{Code: http.StatusNotFound},
			transient: false,
		},
		{
			name:      "400 Bad Request",
			err:       &googleapi.Error{Code: http.StatusBadRequest},
			transient: false,
		},
		{
			name:      "401 Unauthorized",
			err:       &googleapi.Error{Code: http.StatusUnauthorized},
			transient: false,
		},
		{
			name:      "403 Forbidden",
			err:       &googleapi.Error{Code: http.StatusForbidden},
			transient: false,
		},
		{
			name:      "409 Conflict",
			err:       &googleapi.Error{Code: http.StatusConflict},
			transient: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isTransientError(test.err)
			if result != test.transient {
				t.Errorf("isTransientError(%v) = %v, want %v", test.err, result, test.transient)
			}
		})
	}
}
