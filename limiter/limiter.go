// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package limiter

import (
	"context"

	"github.com/drone/autoscaler"
)

type limiter struct {
	autoscaler.ServerStore
	license *License
}

func (l *limiter) Create(ctx context.Context, server *autoscaler.Server) error {
	if l.license.Expired() {
		return errLicenseExpired
	}
	// if the license is unlimited we can continue.
	// no need to calcular server capacity.
	if l.license.Lim == 0 {
		return l.ServerStore.Create(ctx, server)
	}
	count, err := l.count(ctx)
	if err != nil {
		return err
	}
	if count >= l.license.Lim {
		return errServerLimitExceeded
	}
	return l.ServerStore.Create(ctx, server)
}

// count is a helper function that returns the count or
// pending or running server instances.
func (l *limiter) count(ctx context.Context) (int, error) {
	servers, err := l.ServerStore.List(ctx)
	if err != nil {
		return 0, err
	}
	var count int
	for _, server := range servers {
		switch server.State {
		case autoscaler.StatePending,
			autoscaler.StateRunning,
			autoscaler.StateCreated,
			autoscaler.StateCreating:
			count++
		}
	}
	return count, nil
}
