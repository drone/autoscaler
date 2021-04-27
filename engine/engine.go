// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/logger"
	"github.com/drone/autoscaler/metrics"
	"github.com/drone/drone-go/drone"
)

// defines the interval at which terminated instances are
// purged from the database.
const purge = time.Hour * 24

type engine struct {
	mu sync.Mutex

	allocator *allocator
	collector *collector
	installer *installer
	pinger    *pinger
	planner   *planner
	reaper    *reaper
	metrics   metrics.Collector

	interval time.Duration
	paused   bool
}

// New returns a new autoscale Engine.
func New(
	client drone.Client,
	config config.Config,
	servers autoscaler.ServerStore,
	provider autoscaler.Provider,
	metrics metrics.Collector,
) autoscaler.Engine {
	return &engine{
		paused:   false,
		interval: config.Interval,
		allocator: &allocator{
			servers:  servers,
			provider: provider,
			metrics:  metrics,
		},
		collector: &collector{
			timeout:  config.Timeout.Stop,
			servers:  servers,
			provider: provider,
			client:   newDockerClient,
		},
		installer: &installer{
			metrics:                  metrics,
			servers:                  servers,
			os:                       config.Agent.OS,
			arch:                     config.Agent.Arch,
			image:                    config.Agent.Image,
			secret:                   config.Agent.Token,
			envs:                     config.Agent.Environ,
			volumes:                  config.Agent.Volumes,
			labels:                   config.Agent.Labels,
			loggingDriver:            config.Agent.LoggingDriver,
			loggingOptions:           config.Agent.LoggingOptions,
			proto:                    config.Server.Proto,
			host:                     config.Server.Host,
			client:                   newDockerClient,
			runner:                   config.Runner,
			checkInterval:            config.Check.Interval,
			checkDeadline:            config.Check.Deadline,
			gcEnabled:                config.GC.Enabled,
			gcDebug:                  config.GC.Debug,
			gcImage:                  config.GC.Image,
			gcIgnore:                 config.GC.Images,
			gcInterval:               config.GC.Interval,
			gcCache:                  config.GC.Cache,
			gcLoggingDriver:          config.GC.LoggingDriver,
			gcLoggingOptions:         config.GC.LoggingOptions,
			watchtowerEnabled:        config.Watchtower.Enabled,
			watchtowerImage:          config.Watchtower.Image,
			watchtowerTimeout:        config.Watchtower.Timeout,
			watchtowerInterval:       config.Watchtower.Interval,
			watchtowerLoggingDriver:  config.Watchtower.LoggingDriver,
			watchtowerLoggingOptions: config.Watchtower.LoggingOptions,
		},
		pinger: &pinger{
			servers: servers,
			client:  newDockerClient,
		},
		planner: &planner{
			client:  client,
			servers: servers,
			os:      config.Agent.OS,
			arch:    config.Agent.Arch,
			version: config.Agent.Version,
			kernel:  config.Agent.Kernel,
			buffer:  config.CapacityBuffer,
			ttu:     config.Pool.MinAge,
			min:     config.Pool.Min,
			max:     config.Pool.Max,
			cap:     config.Agent.Concurrency,
			labels:  config.Agent.Labels,
		},
		reaper: &reaper{
			servers:  servers,
			provider: provider,
			interval: config.Reaper.Interval,
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
	e.reset(ctx)

	var wg sync.WaitGroup
	wg.Add(7)
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
	go func() {
		e.reap(ctx)
		wg.Done()
	}()
	go func() {
		e.ping(ctx)
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

// runs the ping process.
func (e *engine) ping(ctx context.Context) {
	const interval = time.Minute * 10
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			e.pinger.Ping(ctx)
		}
	}
}

// runs the purge process.
func (e *engine) purge(ctx context.Context) {
	const interval = time.Hour * 24
	const retain = time.Hour * 24 * -1

	logger := logger.FromContext(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			logger.WithField("ttl", retain.String()).
				Debugln("clear stopped servers from database")
			e.planner.servers.Purge(ctx, time.Now().Add(retain).Unix())
		}
	}
}

// runs the reaper process.
func (e *engine) reap(ctx context.Context) {
	// by default, the reaper is run hourly since in general this
	// should happen infrequently.
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(e.reaper.interval):
			e.reaper.Reap(ctx)
		}
	}
}

func (e *engine) reset(ctx context.Context) {
	// handle the situation where the autoscaler is stopped or
	// restarted during instance setup or teardown. If this happens
	// reset the instance state to resume.
	servers, _ := e.allocator.servers.List(ctx)
	for _, s := range servers {
		switch s.State {
		case autoscaler.StateStaging:
			log := logger.FromContext(ctx).
				WithField("instance", s.Name).
				WithField("address", s.Address).
				WithField("from-state", "staging").
				WithField("to-state", "created")
			log.Infoln("reset instance state")
			s.State = autoscaler.StateCreated
			if err := e.allocator.servers.Update(ctx, s); err != nil {
				log.WithError(err).
					Error("failed to reset instance state")
			}
		case autoscaler.StateStopping:
			log := logger.FromContext(ctx).
				WithField("instance", s.Name).
				WithField("address", s.Address).
				WithField("from-state", "stopping").
				WithField("to-state", "shutdown")
			log.Infoln("reset instance state")
			s.State = autoscaler.StateShutdown
			if err := e.allocator.servers.Update(ctx, s); err != nil {
				log.WithError(err).
					Errorln("failed to reset instance state")
			}
		}
	}
}
