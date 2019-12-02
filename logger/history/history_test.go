// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package history

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func TestLevels(t *testing.T) {
	hook := New()
	if diff := cmp.Diff(hook.Levels(), logrus.AllLevels); diff != "" {
		t.Errorf("Hook should return all levels")
		t.Log(diff)
	}
}

func TestConvertLevels(t *testing.T) {
	tests := []struct {
		before logrus.Level
		after  Level
	}{
		{logrus.PanicLevel, LevelError},
		{logrus.FatalLevel, LevelError},
		{logrus.ErrorLevel, LevelError},
		{logrus.WarnLevel, LevelWarn},
		{logrus.InfoLevel, LevelInfo},
		{logrus.DebugLevel, LevelDebug},
		{logrus.TraceLevel, LevelTrace},
	}
	if len(tests) != len(logrus.AllLevels) {
		t.Errorf("missing unit tests for all logrus levels")
	}
	for _, test := range tests {
		if got, want := convertLevel(test.before), test.after; got != want {
			t.Errorf("Want entry level %v, got %v", want, got)
		}
	}
}

func TestLimit(t *testing.T) {
	hook := NewLimit(4)
	hook.Fire(&logrus.Entry{})
	hook.Fire(&logrus.Entry{})
	hook.Fire(&logrus.Entry{})
	hook.Fire(&logrus.Entry{})
	hook.Fire(&logrus.Entry{})
	if got, want := len(hook.entries), 4; got != want {
		t.Errorf("Expect entries pruned to %d, got %d", want, got)
	}
}

func TestHistory(t *testing.T) {
	hook := New()

	now := time.Now()
	hook.Fire(&logrus.Entry{
		Level:   logrus.DebugLevel,
		Message: "foo",
		Data:    logrus.Fields{"foo": "bar"},
		Time:    now,
	})

	hook.Fire(&logrus.Entry{
		Level:   logrus.InfoLevel,
		Message: "bar",
		Data:    logrus.Fields{"baz": "qux"},
		Time:    now,
	})

	if len(hook.entries) != 2 {
		t.Errorf("Expected 2 hooks added to history")
	}

	entries := hook.Entries()
	if len(entries) != 2 {
		t.Errorf("Expected 2 hooks returned")
	}
	if entries[0] == hook.entries[0] {
		t.Errorf("Expect copy of entries, got a reference")
	}
	if entries[1] == hook.entries[1] {
		t.Errorf("Expect copy of entries, got a reference")
	}

	expect := []*Entry{
		{
			Level:   LevelDebug,
			Message: "foo",
			Data:    logrus.Fields{"foo": "bar"},
			Unix:    now.Unix(),
		},
		{
			Level:   LevelInfo,
			Message: "bar",
			Data:    logrus.Fields{"baz": "qux"},
			Unix:    now.Unix(),
		},
	}
	if diff := cmp.Diff(entries, expect); diff != "" {
		t.Errorf("Entries should return an exact copy of all entries")
		t.Log(diff)
	}
}

func TestFilter(t *testing.T) {
	hook := New()

	now := time.Now()
	hook.Fire(&logrus.Entry{
		Level:   logrus.DebugLevel,
		Message: "foo",
		Data:    logrus.Fields{"foo": "bar"},
		Time:    now,
	})

	hook.Fire(&logrus.Entry{
		Level:   logrus.InfoLevel,
		Message: "bar",
		Data:    logrus.Fields{"baz": "qux"},
		Time:    now,
	})

	expect := []*Entry{
		{
			Level:   LevelDebug,
			Message: "foo",
			Data:    logrus.Fields{"foo": "bar"},
			Unix:    now.Unix(),
		},
	}
	entries := hook.Filter(func(entry *Entry) bool {
		return entry.Data["foo"] == "bar"
	})
	if diff := cmp.Diff(entries, expect); diff != "" {
		t.Errorf("Entries should return an exact copy of all entries")
		t.Log(diff)
	}
}
