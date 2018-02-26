// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import (
	"context"
)

// A Scaler implements an algorithm to automatically scale up
// or scale down the available pool of servers.
type Scaler interface {
	Pause()
	Paused() bool
	Resume()

	Scale(context.Context) error
}
