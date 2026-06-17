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
// 429 Too Many Requests is intentionally excluded — it is handled separately
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

// retryAfterError is a 429 rate-limit error carrying the server-requested
// back-off duration parsed from the Retry-After response header.
type retryAfterError struct {
	error
	retryAfter time.Duration
}

// asRateLimitError converts a 429 googleapi.Error into a *retryAfterError.
// Returns nil for any other error.
func asRateLimitError(err error) *retryAfterError {
	gerr, ok := err.(*googleapi.Error)
	if !ok || gerr.Code != http.StatusTooManyRequests {
		return nil
	}
	return &retryAfterError{
		error:      err,
		retryAfter: parseRetryAfter(gerr.Header.Get("Retry-After")),
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

const maxTransientDelay = 5 * time.Second

// retryAfterDelayType is a retry-go DelayTypeFunc that returns the Retry-After
// duration for 429 rate-limit errors and otherwise applies exponential backoff
// with jitter, capped at maxTransientDelay.
func retryAfterDelayType(n uint, err error, config *retry.Config) time.Duration {
	if rl, ok := err.(*retryAfterError); ok {
		return rl.retryAfter
	}
	d := retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)(n, err, config)
	if d > maxTransientDelay {
		return maxTransientDelay
	}
	return d
}

// doWithRetry executes fn with standard retry logic for transient Google API errors.
// - Transient errors (5xx, timeouts) are retried up to 5 times with exponential backoff
// - Rate limit errors (429) are wrapped as unrecoverable to break out for external handling
// - Other errors fail immediately without retry
// The returned error may be a *retryAfterError for rate limits, allowing the caller to
// implement backoff before retrying at a higher level (e.g., outer polling loop).
func doWithRetry(ctx context.Context, operation string, fn func() error) error {
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
		retry.DelayType(retryAfterDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.FromContext(ctx).
				WithField("attempt", n+1).
				WithField("operation", operation).
				WithError(err).
				Debugln("retrying transient error")
		}),
	)
}
