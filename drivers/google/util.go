package google

import (
	"net/http"

	"google.golang.org/api/googleapi"
)

// transientErrorCodes defines HTTP status codes that should trigger a retry
// This may be a non-exhaustive list that requires updates in the future.
var transientErrorCodes = map[int]bool{
	http.StatusInternalServerError: true, // 500
	http.StatusBadGateway:          true, // 502
	http.StatusServiceUnavailable:  true, // 503
	http.StatusGatewayTimeout:      true, // 504
	http.StatusTooManyRequests:     true, // 429
}

func isTransientError(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		_, isTransient := transientErrorCodes[gerr.Code]
		return isTransient
	}
	return false
}
