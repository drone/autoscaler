// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package store

import "strings"

// helper function returns true if the error message
// indicates the connection has been reset.
func isConnReset(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(),
		"connection reset by peer")
}
