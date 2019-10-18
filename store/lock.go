// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package store

import "sync"

// NewLocker returns a new database mutex. If the driver
// is mysql or postgres a noop is returned.
func NewLocker(driver string) sync.Locker {
	switch driver {
	case "sqlite3":
		return new(sync.Mutex)
	default:
		return new(noopLocker)
	}
}

type noopLocker struct{}

func (*noopLocker) Lock()   {}
func (*noopLocker) Unlock() {}
