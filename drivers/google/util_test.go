package google

import (
	"net/http"
	"testing"
	"time"

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
			name:      "504 Gateway Timeout",
			err:       &googleapi.Error{Code: http.StatusGatewayTimeout},
			transient: true,
		},
		{
			// 429 is handled via asRateLimitError, not isTransientError
			name:      "429 Too Many Requests",
			err:       &googleapi.Error{Code: http.StatusTooManyRequests},
			transient: false,
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
			name:      "408 RequestTimeout",
			err:       &googleapi.Error{Code: http.StatusRequestTimeout},
			transient: true,
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

func TestAsRateLimitError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantNil     bool
		wantRetryAt time.Duration
	}{
		{
			name:    "nil error",
			err:     nil,
			wantNil: true,
		},
		{
			name:    "non-429 google error",
			err:     &googleapi.Error{Code: http.StatusServiceUnavailable},
			wantNil: true,
		},
		{
			name:        "429 with no Retry-After header",
			err:         &googleapi.Error{Code: http.StatusTooManyRequests, Header: http.Header{}},
			wantNil:     false,
			wantRetryAt: 10 * time.Second,
		},
		{
			name: "429 with integer Retry-After",
			err: &googleapi.Error{
				Code:   http.StatusTooManyRequests,
				Header: http.Header{"Retry-After": []string{"30"}},
			},
			wantNil:     false,
			wantRetryAt: 30 * time.Second,
		},
		{
			name: "429 with zero Retry-After",
			err: &googleapi.Error{
				Code:   http.StatusTooManyRequests,
				Header: http.Header{"Retry-After": []string{"0"}},
			},
			wantNil:     false,
			wantRetryAt: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := asRateLimitError(test.err)
			if test.wantNil {
				if result != nil {
					t.Errorf("asRateLimitError(%v) = %v, want nil", test.err, result)
				}
				return
			}
			if result == nil {
				t.Fatalf("asRateLimitError(%v) = nil, want non-nil", test.err)
			}
			if result.retryAfter != test.wantRetryAt {
				t.Errorf("retryAfter = %v, want %v", result.retryAfter, test.wantRetryAt)
			}
		})
	}
}

func TestParseRetryAfter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{
			name:  "empty string falls back to 10s",
			input: "",
			want:  10 * time.Second,
		},
		{
			name:  "integer seconds",
			input: "60",
			want:  60 * time.Second,
		},
		{
			name:  "zero seconds",
			input: "0",
			want:  0,
		},
		{
			name:  "unparseable falls back to 10s",
			input: "not-a-date",
			want:  10 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseRetryAfter(test.input)
			if got != test.want {
				t.Errorf("parseRetryAfter(%q) = %v, want %v", test.input, got, test.want)
			}
		})
	}
}
