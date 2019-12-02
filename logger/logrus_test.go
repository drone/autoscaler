// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogrus(t *testing.T) {
	logger := Logrus(
		logrus.NewEntry(
			logrus.StandardLogger(),
		),
	)
	if _, ok := logger.(*wrapLogrus); !ok {
		t.Errorf("Expect wrapped logrus")
	}
	if _, ok := logger.WithError(nil).(*wrapLogrus); !ok {
		t.Errorf("Expect WithError wraps logrus")
	}
	if _, ok := logger.WithField("foo", "bar").(*wrapLogrus); !ok {
		t.Errorf("Expect WithField logrus")
	}
}
