// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sort"
	"time"
	"os"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/limiter"
	"github.com/drone/drone-go/drone"

	"github.com/dchest/uniuri"
	"github.com/rs/zerolog/log"
)

// a planner is responsible for capacity planning. It will assess
// current build volume and plan the creation or termination of
// server resources accordingly.
type planner struct {
	min int           // min number of servers
	max int           // max number of servers to allocate
	cap int           // capacity per-server
	ttu time.Duration // minimum server age

	client  drone.Client
	servers autoscaler.ServerStore
}

func (p *planner) Plan(ctx context.Context) error {
	// generate a unique identifier for the current
	// execution cycle for tracing and grouping logs.
	cycle := uniuri.New()

	logger := log.Ctx(ctx).With().Str("id", cycle).Logger()

	pending, running, err := p.count(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("cannot fetch queue details")
		return err
	}

	capacity, servers, err := p.capacity(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("cannot calculate server capacity")
		return err
	}

	logger.Debug().
		Int("min-pool", p.min).
		Int("max-pool", p.max).
		Int("server-capacity", capacity).
		Int("server-count", servers).
		Int("pending-builds", pending).
		Int("running-builds", running).
		Msg("check capacity")

	defer func() {
		logger.Debug().
			Msg("check capacity complete")
	}()

	ctx = logger.WithContext(ctx)

	free := max(capacity-running, 0)
	diff := serverDiff(pending, free, p.cap)

	// if the server differential to handle the build volume
	// is positive, we can reduce server capacity.
	if diff < 0 {
		return p.mark(ctx,
			// we should adjust the desired capacity to ensure
			// we maintain the minimum required server count.
			serverFloor(servers, abs(diff), p.min),
		)
	}

	// if the server differential to handle the build volume
	// is positive, we need to allocate more server capacity.
	if diff > 0 {
		return p.alloc(ctx,
			// we should adjust the desired capacity to ensure
			// it does not exceed the max server count.
			serverCeil(servers, diff, p.max),
		)
	}

	logger.Debug().
		Msg("no capacity changes required")

	return nil
}

// helper function allocates n new server instances.
func (p *planner) alloc(ctx context.Context, n int) error {
	logger := log.Ctx(ctx)

	logger.Debug().
		Msgf("allocate %d servers", n)

	for i := 0; i < n; i++ {
		server := &autoscaler.Server{
			Name:     "agent-" + uniuri.NewLen(8),
			State:    autoscaler.StatePending,
			Secret:   uniuri.New(),
			Capacity: p.cap,
		}

		err := p.servers.Create(ctx, server)
		if limiter.IsError(err) {
			logger.Warn().Err(err).
				Msg("cannot create server")
			return err
		}
		if err != nil {
			logger.Error().Err(err).
				Msg("cannot create server")
			return err
		}
	}
	return nil
}

// helper funciton marks instances for termination.
func (p *planner) mark(ctx context.Context, n int) error {
	logger := log.Ctx(ctx)

	logger.Debug().
		Msgf("terminate %d servers", n)

	if n == 0 {
		return nil
	}

	servers, err := p.servers.ListState(ctx, autoscaler.StateRunning)
	if err != nil {
		logger.Error().Err(err).
			Msg("cannot fetch server list")
		return err
	}
	sort.Sort(sort.Reverse(byCreated(servers)))

	busy, err := p.listBusy(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("cannot ascertain busy server list")
		return err
	}

	var idle []*autoscaler.Server
	for _, server := range servers {
		// skip busy servers
		if _, ok := busy[server.Name]; ok {
			logger.Debug().
				Str("server", server.Name).
				Msg("server is busy")
			continue
		}

		// skip servers less than minage
		if time.Now().Before(time.Unix(server.Created, 0).Add(p.ttu)) {
			logger.Debug().
				Str("server", server.Name).
				TimeDiff("age", time.Now(), time.Unix(server.Created, 0)).
				Dur("min-age", p.ttu).
				Msg("server min-age not reached")
			continue
		}

		idle = append(idle, server)
		logger.Debug().
			Str("server", server.Name).
			Msg("server is idle")
	}

	// if there are no idle servers, there are no servers
	// to retire and we can exit.
	if len(idle) == 0 {
		logger.Debug().
			Msg("no idle servers to shutdown")
	}

	if len(idle) > n {
		idle = idle[:n]
	}

	for _, server := range idle {
		server.State = autoscaler.StateShutdown
		err := p.servers.Update(ctx, server)
		if err != nil {
			logger.Error().
				Err(err).
				Str("server", server.Name).
				Str("state", "shutdown").
				Msg("cannot update server state")
		}
	}

	return nil
}

// helper function returns the number of pending and
// running builds in the remote Drone installation.
func (p *planner) count(ctx context.Context) (pending, running int, err error) {
	if os.Getenv("ENABLE_MATRIX_CALC") == "true" {
		return p.countMatrix(ctx)
	}
	activity, err := p.client.BuildQueue()
	if err != nil {
		return pending, running, err
	}
	for _, activity := range activity {
		if activity.Status == drone.StatusPending {
			pending++
		} else {
			running++
		}
	}
	return
}

func (p *planner) countMatrix(ctx context.Context) (pending, running int, err error) {
	activity, err := p.client.BuildQueue()
	if err != nil {
		return pending, running, err
	}
	for _, activity := range activity {
		build, err := p.client.Build(activity.Owner, activity.Name, activity.Number);
		if err != nil {
			return pending, running, err
		}

		for _, process := range build.Procs {
			switch process.State {
				case drone.StatusPending:
					pending++
				case drone.StatusRunning:
					running++
			}
		}
	}
	return
}

// helper function returns our current capacity.
func (p *planner) capacity(ctx context.Context) (capacity, count int, err error) {
	servers, err := p.servers.List(ctx)
	if err != nil {
		return capacity, count, err
	}
	for _, server := range servers {
		switch server.State {
		case autoscaler.StateStopped:
			// ignore state
		case autoscaler.StateShutdown:
			// ignore state
		case autoscaler.StateStopping:
			// ignore state
		case autoscaler.StateError:
			// ignore state
		default:
			count++
			capacity += server.Capacity
		}
	}
	return
}

// helper function returns a list of busy servers.
func (p *planner) listBusy(ctx context.Context) (map[string]struct{}, error) {
	busy := map[string]struct{}{}
	builds, err := p.client.BuildQueue()
	if err != nil {
		return busy, err
	}
	for _, build := range builds {
		if build.Status != drone.StatusRunning {
			continue
		}
		build, err := p.client.Build(build.Owner, build.Name, build.Number)
		if err != nil {
			return busy, err
		}
		for _, proc := range build.Procs {
			busy[proc.Machine] = struct{}{}
		}
	}
	return busy, nil
}
