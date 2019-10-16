// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package store

import "sync"

var _ sync.Locker = (*noopLocker)(nil)

type noopLocker struct{}

func (*noopLocker) Lock()   {}
func (*noopLocker) Unlock() {}
