// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package store

import (
	"database/sql"
	"errors"
	"testing"
)

func TestConnectionReset(t *testing.T) {
	if isConnReset(nil) {
		t.Errorf("Expect nil error returns false")
	}
	if isConnReset(sql.ErrNoRows) {
		t.Errorf("Expect ErrNoRows returns false")
	}
	if !isConnReset(errors.New("read: connection reset by peer")) {
		t.Errorf("Expect connection reset by peer return true")
	}
}
