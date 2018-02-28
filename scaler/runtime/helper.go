// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"sync"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/drone-go/drone"
)

type scaler struct {
	mu sync.Mutex

	*allocator
	*collector
	*planner

	paused bool
}

// New returns a new autoscaler.
func New(
	client drone.Client,
	config config.Config,
	servers autoscaler.ServerStore,
	provider autoscaler.Provider,
) autoscaler.Scaler {
	return &scaler{
		allocator: &allocator{
			servers:  servers,
			provider: provider,
		},
		collector: &collector{
			servers:  servers,
			provider: provider,
		},
		planner: &planner{
			client:  client,
			servers: servers,
			ttu:     config.Pool.MinAge,
			min:     config.Pool.Min,
			max:     config.Pool.Max,
			cap:     config.Agent.Concurrency,
		},
	}
}

// Pause paueses the scaler.
func (s *scaler) Pause() {
	s.mu.Lock()
	s.paused = true
	s.mu.Unlock()
}

// Paused returns true if scaling is paused.
func (s *scaler) Paused() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paused
}

// Resume resumes the scaler.
func (s *scaler) Resume() {
	s.mu.Lock()
	s.paused = false
	s.mu.Unlock()
}

func (s *scaler) Scale(ctx context.Context) (err error) {
	if s.Paused() == false {
		err = s.planner.Plan(ctx)
	}
	return
}
