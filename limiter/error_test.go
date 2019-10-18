// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package limiter

import (
	"database/sql"
	"testing"
)

func TestIsError(t *testing.T) {
	var tests = []struct {
		err error
		res bool
	}{
		{nil, false},
		{errLicenseExpired, true},
		{errLicenseExpired, true},
		{sql.ErrNoRows, false},
	}
	for _, test := range tests {
		if got, want := IsError(test.err), test.res; got != want {
			t.Errorf("Want IsError %v, got %v", want, got)
		}
	}
}
