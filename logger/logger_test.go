// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package logger

import "testing"

func TestWithError(t *testing.T) {
	d := &discard{}
	if d.WithError(nil) != d {
		t.Errorf("Expect WithError to return base logger")
	}
}

func TestWithField(t *testing.T) {
	d := &discard{}
	if d.WithField("hello", "world") != d {
		t.Errorf("Expect WithField to return base logger")
	}
}
