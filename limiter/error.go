// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package limiter

import "errors"

var (
	// Indicates the system has reached the limit on the
	// number servers that it can provision under the
	// current license.
	errServerLimitExceeded = errors.New("Server limit exceeded")

	// Indicates the license is expried. No new servers are
	// provisioned until the license is renewed.
	errLicenseExpired = errors.New("License expired")
)

// IsError returns true if the error is a Limit error.
func IsError(err error) bool {
	switch err {
	case errServerLimitExceeded, errLicenseExpired:
		return true
	default:
		return false
	}
}
