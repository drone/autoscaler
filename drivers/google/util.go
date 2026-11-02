package google

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/drone/autoscaler/logger"
	"google.golang.org/api/googleapi"
)

// transientErrorCodes defines HTTP status codes that should trigger a retry.
// Rate limits are intentionally excluded — they are handled separately
// via asRateLimitError so that the Retry-After header can be honoured.
var transientErrorCodes = map[int]bool{
	http.StatusRequestTimeout:      true, // 408
	http.StatusInternalServerError: true, // 500
	http.StatusBadGateway:          true, // 502
	http.StatusServiceUnavailable:  true, // 503
	http.StatusGatewayTimeout:      true, // 504
}

// isTransientError reports whether err should be retried strictly based on
// HTTP status codes; not on any underlying network issues.
func isTransientError(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		_, isTransient := transientErrorCodes[gerr.Code]
		return isTransient
	}
	return false
}

// retryAfterError is a rate-limit error carrying the server-requested
// back-off duration parsed from the Retry-After response header (or a
// default when the response didn't send one).
type retryAfterError struct {
	Err        error
	retryAfter time.Duration
}

// Error implements the error interface.
func (e *retryAfterError) Error() string { return e.Err.Error() }

// Unwrap lets errors.Is/errors.As see through to the underlying error.
func (e *retryAfterError) Unwrap() error { return e.Err }

// asRateLimitError converts a rate-limited googleapi.Error into a
// *retryAfterError. Returns nil for any other error.
func asRateLimitError(err error) *retryAfterError {
	gerr, ok := err.(*googleapi.Error)
	if !ok || !isRateLimitError(gerr) {
		return nil
	}
	return &retryAfterError{
		Err:        err,
		retryAfter: parseRetryAfter(gerr.Header.Get("Retry-After")),
	}
}

// isRateLimitError reports whether gerr represents a rate limit, in either
// of the two shapes Compute uses:
//   - 429 Too Many Requests, usually with a Retry-After header.
//   - 403 Forbidden with an ErrorItem Reason of "rateLimitExceeded", which is
//     how Compute reports an exhausted per-project-per-region quota (e.g.
//     "Write requests per minute per region"), without a Retry-After header.
//     Other 403s (permission denied) have a different Reason and are not
//     retryable.
func isRateLimitError(gerr *googleapi.Error) bool {
	if gerr.Code == http.StatusTooManyRequests {
		return true
	}
	if gerr.Code == http.StatusForbidden {
		for _, item := range gerr.Errors {
			if item.Reason == "rateLimitExceeded" {
				return true
			}
		}
	}
	return false
}

// sleepOrDone waits out d, returning ctx.Err() if ctx ends first.
func sleepOrDone(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// parseRetryAfter parses a Retry-After header value (integer seconds or HTTP-date).
// Falls back to 10s when the header is absent or unparseable.
func parseRetryAfter(s string) time.Duration {
	if s == "" {
		return 10 * time.Second
	}
	if secs, err := strconv.Atoi(s); err == nil {
		if secs <= 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(s); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
		return 0
	}
	return 10 * time.Second
}

// retryAfterDelayType returns a retry-go DelayTypeFunc that returns the
// Retry-After duration for rate-limit errors and otherwise applies
// exponential backoff with jitter, capped at maxDelay. It is a constructor
// because retry.DelayTypeFunc's signature has no room for a config value.
func retryAfterDelayType(maxDelay time.Duration) retry.DelayTypeFunc {
	return func(n uint, err error, config *retry.Config) time.Duration {
		if rl, ok := err.(*retryAfterError); ok {
			return rl.retryAfter
		}
		d := retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)(n, err, config)
		if d > maxDelay {
			return maxDelay
		}
		return d
	}
}

// doWithRetry executes fn with standard retry logic for transient Google API errors.
// - Transient errors (5xx, timeouts) are retried up to 5 times with exponential backoff
// - Rate limit errors (see isRateLimitError) are wrapped as unrecoverable to break out for external handling
// - Other errors fail immediately without retry
// Backoff between transient retries is capped at maxDelay.
// The returned error may be a *retryAfterError for rate limits, allowing the caller to
// implement backoff before retrying at a higher level (e.g., outer polling loop).
func doWithRetry(ctx context.Context, operation string, maxDelay time.Duration, fn func() error) error {
	return retry.Do(
		func() error {
			err := fn()
			if err == nil {
				return nil
			}
			if rl := asRateLimitError(err); rl != nil {
				return retry.Unrecoverable(rl)
			}
			if isTransientError(err) {
				return err
			}
			return retry.Unrecoverable(err)
		},
		retry.Attempts(5),
		retry.Context(ctx),
		retry.LastErrorOnly(true),
		retry.DelayType(retryAfterDelayType(maxDelay)),
		retry.OnRetry(func(n uint, err error) {
			logger.FromContext(ctx).
				WithField("attempt", n+1).
				WithField("operation", operation).
				WithError(err).
				Debugln("retrying transient error")
		}),
	)
}
