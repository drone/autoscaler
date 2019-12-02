// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package logger defines interfaces that logger drivers
// implement to log messages.
package logger

// A Logger represents an active logging object that generates
// lines of output to an io.Writer.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Traceln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	WithError(error) Logger
	WithField(string, interface{}) Logger
}

// Default returns the default logger.
var Default = Discard()

// Discard returns a no-op logger
func Discard() Logger {
	return &discard{}
}

type discard struct{}

func (*discard) Debug(args ...interface{})                 {}
func (*discard) Debugf(format string, args ...interface{}) {}
func (*discard) Debugln(args ...interface{})               {}
func (*discard) Error(args ...interface{})                 {}
func (*discard) Errorf(format string, args ...interface{}) {}
func (*discard) Errorln(args ...interface{})               {}
func (*discard) Info(args ...interface{})                  {}
func (*discard) Infof(format string, args ...interface{})  {}
func (*discard) Infoln(args ...interface{})                {}
func (*discard) Trace(args ...interface{})                 {}
func (*discard) Tracef(format string, args ...interface{}) {}
func (*discard) Traceln(args ...interface{})               {}
func (*discard) Warn(args ...interface{})                  {}
func (*discard) Warnf(format string, args ...interface{})  {}
func (*discard) Warnln(args ...interface{})                {}
func (d *discard) WithError(error) Logger                  { return d }
func (d *discard) WithField(string, interface{}) Logger    { return d }
