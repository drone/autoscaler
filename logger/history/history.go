// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package history

// Package history implements a logrus hook that provides access
// to log recent log activity.
import (
	"sync"

	"github.com/sirupsen/logrus"
)

// default log entry limit.
const defaultLimit = 250

// Level is the log level.
type Level string

// log levels.
const (
	LevelError = Level("error")
	LevelWarn  = Level("warn")
	LevelInfo  = Level("info")
	LevelDebug = Level("debug")
	LevelTrace = Level("trace")
)

// Entry provides a log entry.
type Entry struct {
	Level   Level
	Message string
	Data    map[string]interface{}
	Unix    int64
}

// Hook is a logrus hook that track the log history.
type Hook struct {
	sync.RWMutex
	limit   int
	entries []*Entry
}

// New returns a new history hook.
func New() *Hook {
	return NewLimit(defaultLimit)
}

// NewLimit returns a new history hook with a custom
// history limit.
func NewLimit(limit int) *Hook {
	return &Hook{limit: limit}
}

// Fire receives the log entry.
func (h *Hook) Fire(e *logrus.Entry) error {
	h.Lock()
	if len(h.entries) >= h.limit {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, &Entry{
		Level:   convertLevel(e.Level),
		Data:    convertFields(e.Data),
		Unix:    e.Time.Unix(),
		Message: e.Message,
	})
	h.Unlock()
	return nil
}

// Levels returns the supported log levels.
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Entries returns a list of all entries.
func (h *Hook) Entries() []*Entry {
	h.RLock()
	defer h.RUnlock()
	entries := make([]*Entry, len(h.entries))
	for i, entry := range h.entries {
		entries[i] = copyEntry(entry)
	}
	return entries
}

// Filter returns a list of all entries for which the filter
// function returns true.
func (h *Hook) Filter(filter func(*Entry) bool) []*Entry {
	h.RLock()
	defer h.RUnlock()
	var entries []*Entry
	for _, entry := range h.entries {
		if filter(entry) {
			entries = append(entries, copyEntry(entry))
		}
	}
	return entries
}

// helper funtion copies an entry for threadsafe access.
func copyEntry(src *Entry) *Entry {
	dst := new(Entry)
	*dst = *src
	dst.Data = map[string]interface{}{}
	for k, v := range src.Data {
		dst.Data[k] = v
	}
	return dst
}

// helper function converts a logrus.Level to the local type.
func convertLevel(level logrus.Level) Level {
	switch level {
	case logrus.PanicLevel:
		return LevelError
	case logrus.FatalLevel:
		return LevelError
	case logrus.ErrorLevel:
		return LevelError
	case logrus.WarnLevel:
		return LevelWarn
	case logrus.DebugLevel:
		return LevelDebug
	case logrus.TraceLevel:
		return LevelTrace
	default:
		return LevelInfo
	}
}

// helper fucntion copies logrus.Fields to a basic map.
func convertFields(src logrus.Fields) map[string]interface{} {
	dst := map[string]interface{}{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
