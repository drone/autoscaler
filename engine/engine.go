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

	"github.com/rs/zerolog/log"
)

// defines the interval at which terminated instances are
// purged from the database.
const purge = time.Hour * 24

type engine struct {
	mu sync.Mutex

	allocator *allocator
	collector *collector
	installer *installer
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
			client:   newDockerClient,
		},
		installer: &installer{
			servers: servers,
			image:   config.Agent.Image,
			secret:  config.Agent.Token,
			server:  config.Agent.Host,
			keepaliveTime: config.Agent.KeepaliveTime,
			keepaliveTimeout: config.Agent.KeepaliveTimeout,
			client:  newDockerClient,
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
	wg.Add(5)
	go func() {
		e.allocate(ctx)
		wg.Done()
	}()
	go func() {
		e.install(ctx)
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

// runs the installation process.
func (e *engine) install(ctx context.Context) {
	// Any Staging server is not considered installed so revert its state to run the installer again.
	// This happens when the autoscaler is stopped after the server is created, but before the installation is complete.
	stagings, err := e.allocator.servers.ListState(ctx, autoscaler.StateStaging)
	if err != nil {
		log.Warn().Err(err)
	} else {
		for _, s := range stagings {
			s.State = autoscaler.StateCreated
			err = e.allocator.servers.Update(ctx, s)
		}
	}

	const interval = time.Second * 10
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			e.installer.Install(ctx)
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
	// if e.planner.min+e.planner.max == 0 {
	// 	log.Ctx(ctx).Warn().
	// 		Msg("autoscaling disabled: pool min and max are zero")
	// 	return
	// }
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

	logger := log.Ctx(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			logger.Debug().
				Str("ttl", retain.String()).
				Msg("clear stopped servers from database")
			e.planner.servers.Purge(ctx, time.Now().Add(retain).Unix())
		}
	}
}
