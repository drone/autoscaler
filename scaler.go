// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import (
	"context"
)

// A implements an algorithm to automatically scale up
// or scale down the available pool of servers.
type Scaler interface {
	// Allocate creates instances.
	Allocate(context.Context) error

	// Collect removes instances marked for shutdown.
	Collect(context.Context) error

	// Pause pauses the scaling routine.
	Pause()

	// Paused returns true if the scaling routine is paused.
	Paused() bool

	// Resume resumes the scaling routing.
	Resume()

	// Plan(context.Context) error
	// Process(context.Context) error

	Scale(context.Context) error
}
