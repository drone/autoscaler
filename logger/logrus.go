// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package logger

import "github.com/sirupsen/logrus"

// Logrus returns a Logger that wraps a logrus.Entry.
func Logrus(entry *logrus.Entry) Logger {
	return &wrapLogrus{entry}
}

type wrapLogrus struct {
	*logrus.Entry
}

func (w *wrapLogrus) WithError(err error) Logger {
	return &wrapLogrus{w.Entry.WithError(err)}
}

func (w *wrapLogrus) WithField(key string, value interface{}) Logger {
	return &wrapLogrus{w.Entry.WithField(key, value)}
}
