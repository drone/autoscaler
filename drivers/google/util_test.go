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

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name        string
		err         *googleapi.Error
		rateLimited bool
	}{
		{
			name:        "429 Too Many Requests",
			err:         &googleapi.Error{Code: http.StatusTooManyRequests},
			rateLimited: true,
		},
		{
			name: "403 with rateLimitExceeded reason",
			err: &googleapi.Error{
				Code: http.StatusForbidden,
				Errors: []googleapi.ErrorItem{
					{Reason: "rateLimitExceeded"},
				},
			},
			rateLimited: true,
		},
		{
			name:        "403 with no Errors items",
			err:         &googleapi.Error{Code: http.StatusForbidden},
			rateLimited: false,
		},
		{
			// A genuine permission-denied 403 must not be treated as
			// retryable - only the specific rateLimitExceeded reason is.
			name: "403 with an unrelated reason",
			err: &googleapi.Error{
				Code: http.StatusForbidden,
				Errors: []googleapi.ErrorItem{
					{Reason: "insufficientPermissions"},
				},
			},
			rateLimited: false,
		},
		{
			name: "rateLimitExceeded reason on a non-403/429 code",
			err: &googleapi.Error{
				Code: http.StatusBadRequest,
				Errors: []googleapi.ErrorItem{
					{Reason: "rateLimitExceeded"},
				},
			},
			rateLimited: false,
		},
		{
			name:        "unrelated error code",
			err:         &googleapi.Error{Code: http.StatusInternalServerError},
			rateLimited: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isRateLimitError(test.err)
			if result != test.rateLimited {
				t.Errorf("isRateLimitError(%v) = %v, want %v", test.err, result, test.rateLimited)
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
		{
			// Matches the shape Compute actually returns for an exhausted
			// per-project-per-region write quota - a 403, not a 429, and
			// never with a Retry-After header, hence the fallback default.
			name: "403 write-quota exceeded, no Retry-After header",
			err: &googleapi.Error{
				Code:    http.StatusForbidden,
				Header:  http.Header{},
				Message: "Quota exceeded for quota metric 'Write requests' and limit 'Write requests per minute per region' of service 'compute.googleapis.com'",
				Errors: []googleapi.ErrorItem{
					{
						Reason:  "rateLimitExceeded",
						Message: "Quota exceeded for quota metric 'Write requests' and limit 'Write requests per minute per region' of service 'compute.googleapis.com'",
					},
				},
			},
			wantNil:     false,
			wantRetryAt: 10 * time.Second,
		},
		{
			name: "403 with an unrelated reason is not rate-limited",
			err: &googleapi.Error{
				Code: http.StatusForbidden,
				Errors: []googleapi.ErrorItem{
					{Reason: "insufficientPermissions"},
				},
			},
			wantNil: true,
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
