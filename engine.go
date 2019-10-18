// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package autoscaler

import "context"

// An Engine is responsible for running the scaling
// alogirthm to provision and shutdown instances according
// to build volume.
type Engine interface {
	// Start starts the Engine. The context can be used
	// to cancel a running engine.
	Start(context.Context)
	// Pause pauses the Engine.
	Pause()
	// Paused returns true if th Engine is paused.
	Paused() bool
	// Resume resumes the Engine if paused.
	Resume()
}
