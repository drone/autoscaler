// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/drone-go/drone"
)

// defines the interval at which terminated instances are
// purged from the database.
const purge = time.Hour * 24

type engine struct {
	mu sync.Mutex

	allocator *allocator
	collector *collector
	planner   *planner

	interval time.Duration
	paused   bool
}

// New returns a new autoscale Engine.
func New(
	client drone.Client,
	config config.Config,
	servers autoscaler.ServerStore,
	provider autoscaler.Provider,
) autoscaler.Engine {
	return &engine{
		paused:   false,
		interval: config.Interval,
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
func (e *engine) Pause() {
	e.mu.Lock()
	e.paused = true
	e.mu.Unlock()
}

// Paused returns true if scaling is paused.
func (e *engine) Paused() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.paused
}

// Resume resumes the scaler.
func (e *engine) Resume() {
	e.mu.Lock()
	e.paused = false
	e.mu.Unlock()
}

func (e *engine) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		e.allocate(ctx)
		wg.Done()
	}()
	go func() {
		e.collect(ctx)
		wg.Done()
	}()
	go func() {
		e.plan(ctx)
		wg.Done()
	}()
	go func() {
		e.purge(ctx)
		wg.Done()
	}()
	wg.Wait()
}

// runs the allocation process.
func (e *engine) allocate(ctx context.Context) {
	const interval = time.Second * 10
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			e.allocator.Allocate(ctx)
		}
	}
}

// runs the collection process.
func (e *engine) collect(ctx context.Context) {
	const interval = time.Second * 10
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			e.collector.Collect(ctx)
		}
	}
}

// runs the planning process.
func (e *engine) plan(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(e.interval):
			if !e.Paused() {
				e.planner.Plan(ctx)
			}
		}
	}
}

// runs the purge process.
func (e *engine) purge(ctx context.Context) {
	const interval = time.Hour * 24
	const retain = time.Hour * 24 * -1
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			e.planner.servers.Purge(ctx, time.Now().Add(retain).Unix())
		}
	}
}
