// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scaler

import (
	"context"
	"time"

	"github.com/drone/autoscaler"
)

// Start executes the synchronizer in a loop.
func Start(ctx context.Context, scaler autoscaler.Scaler, duration time.Duration) error {
	for {
		select {
		case <-time.After(duration):
			if !scaler.Paused() {
				scaler.Scale(ctx)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
