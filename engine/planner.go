// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sort"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	"github.com/drone/drone-go/drone"

	"github.com/dchest/uniuri"
)

// a planner is responsible for capacity planning. It will assess
// current build volume and plan the creation or termination of
// server resources accordingly.
type planner struct {
	os         string
	arch       string
	version    string
	kernel     string
	namePrefix string
	min        int           // min number of servers
	max        int           // max number of servers to allocate
	cap        int           // capacity per-server
	buffer     int           // buffer capacity to have warm and ready
	ttu        time.Duration // minimum server age
	tti        time.Duration // minimum server idle time
	labels     map[string]string

	client  drone.Client
	servers autoscaler.ServerStore
}

func (p *planner) Plan(ctx context.Context) error {
	// generate a unique identifier for the current
	// execution cycle for tracing and grouping logs.
	cycle := uniuri.New()

	log := logger.FromContext(ctx).WithField("id", cycle)
	log.Debugln("calculate unfinished jobs")

	pending, running, err := p.count(ctx)
	if err != nil {
		log.WithError(err).
			Errorln("cannot calculate unfinished jobs")
		return err
	}

	log.Debugln("calculate server capacity")

	capacity, servers, err := p.capacity(ctx)
	if err != nil {
		log.WithError(err).
			Errorln("cannot calculate server capacity")
		return err
	}

	log.
		WithField("min-pool", p.min).
		WithField("max-pool", p.max).
		WithField("server-buffer", p.buffer).
		WithField("server-capacity", capacity).
		WithField("server-count", servers).
		WithField("pending-builds", pending).
		WithField("running-builds", running).
		Debugln("check capacity")

	defer func() {
		log.Debugln("check capacity complete")
	}()

	ctx = logger.WithContext(ctx, log)

	// if MinIdle is being used, track busy servers
	if p.tti > 0 {
		_, err = p.updateBusy(ctx)
		if err != nil {
			log.WithError(err).
				Errorln("cannot check for busy servers")
			return err
		}
	}

	free := max(capacity-running-p.buffer, 0)
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

	log.Debugln("no capacity changes required")

	return nil
}


// helper function checks for busy running instances and updates idle timer
func (p *planner) updateBusy(ctx context.Context) (count int, err error) {
	logger := logger.FromContext(ctx)

	servers, err := p.servers.ListState(ctx, autoscaler.StateRunning)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot fetch server list")
		return count, err
	}

	// check for busy servers to update idle timers
	busy, err := p.listBusy(ctx)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot ascertain busy server list")
		return count, err
	}

	for _, server := range servers {
		if _, ok := busy[server.Name]; ok {
			err := p.servers.Busy(ctx, server)
			if err != nil {
				logger.WithError(err).
					WithField("server", server.Name).
					WithField("updated", server.Updated).
					Errorln("cannot update server as busy")
			}
			logger.WithField("server", server.Name).
				Debugln("updated busy server")
			count++
		}
	}
	logger.Debugf("found %d busy servers", count)
	return count, nil
}


// helper function allocates n new server instances.
func (p *planner) alloc(ctx context.Context, n int) error {
	logger := logger.FromContext(ctx)
	logger.Debugf("allocate %d servers", n)

	namePrefix := p.namePrefix
	if namePrefix == "" {
		namePrefix = "agent-"
	}
	for i := 0; i < n; i++ {
		server := &autoscaler.Server{
			Name:     p.namePrefix + uniuri.NewLen(8),
			State:    autoscaler.StatePending,
			Secret:   uniuri.New(),
			Capacity: p.cap,
		}

		err := p.servers.Create(ctx, server)
		if err != nil {
			logger.WithError(err).
				Errorln("cannot create server")
			return err
		}
	}
	return nil
}

// helper function marks instances for termination.
func (p *planner) mark(ctx context.Context, n int) error {
	logger := logger.FromContext(ctx)
	logger.Debugf("terminate %d servers", n)

	if n == 0 {
		return nil
	}

	servers, err := p.servers.ListState(ctx, autoscaler.StateRunning)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot fetch server list")
		return err
	}
	sort.Sort(sort.Reverse(byCreated(servers)))

	// Abort marking servers for termination if the total
	// number of running servers, minus the total number
	// of servers to terminate, falls below the minimum
	// number of servers (including the buffer).
	if len(servers)-n < p.min {
		logger.WithField("servers-to-terminate", n).
			WithField("servers-running", len(servers)).
			WithField("min-pool", p.min).
			Debugf("abort terminating instances to ensure minimum capacity met")
		return nil
	}

	busy, err := p.listBusy(ctx)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot ascertain busy server list")
		return err
	}

	var idle []*autoscaler.Server
	for _, server := range servers {
		// skip busy servers
		if _, ok := busy[server.Name]; ok {
			logger.WithField("server", server.Name).
				Debugln("server is busy")
			continue
		}

		// skip servers less than minage
		if time.Now().Before(time.Unix(server.Created, 0).Add(p.ttu)) {
			logger.
				WithField("server", server.Name).
				WithField("age", timeDiff(time.Now(), time.Unix(server.Created, 0))).
				WithField("min-age", p.ttu).
				Debugln("server min-age not reached")
			continue
		}

		// skip servers that have not reached a min idle time
		if time.Now().Before(time.Unix(server.LastBusy, 0).Add(p.tti)) {
			logger.
				WithField("server", server.Name).
				WithField("idle", timeDiff(time.Now(), time.Unix(server.LastBusy, 0))).
				WithField("min-idle", p.tti).
				Debugln("server min-idle not reached")
			continue
		}

		idle = append(idle, server)
		logger.WithField("server", server.Name).
			Debugln("server is idle")
	}

	// if there are no idle servers, there are no servers
	// to retire and we can exit.
	if len(idle) == 0 {
		logger.Debugln("no idle servers to shutdown")
	}

	if len(idle) > n {
		idle = idle[:n]
	}

	for _, server := range idle {
		server.State = autoscaler.StateShutdown
		err := p.servers.Update(ctx, server)
		if err != nil {
			logger.WithError(err).
				WithField("server", server.Name).
				WithField("state", "shutdown").
				Errorln("cannot update server state")
		}
	}

	return nil
}

// helper function returns the number of pending and
// running builds in the remote Drone installation.
func (p *planner) count(ctx context.Context) (pending, running int, err error) {
	stages, err := p.client.Queue()
	if err != nil {
		return pending, running, err
	}
	for _, stage := range stages {
		if p.match(stage) == false {
			continue
		}
		switch stage.Status {
		case drone.StatusPending:
			pending++
		case drone.StatusRunning:
			running++
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
	stages, err := p.client.Queue()
	if err != nil {
		return busy, err
	}
	for _, stage := range stages {
		if p.match(stage) == false {
			continue
		}
		if stage.Status == drone.StatusRunning {
			busy[stage.Machine] = struct{}{}
		}
	}
	return busy, nil
}

// helper function returns true if the os, arch, variant
// and kernel match the stage.
func (p *planner) match(stage *drone.Stage) bool {
	labelMatch := true

	if len(p.labels) > 0 || len(stage.Labels) > 0 {
		labelMatch = checkLabels(p.labels, stage.Labels)
	}

	return stage.OS == p.os &&
		stage.Arch == p.arch &&
		stage.Variant == p.version &&
		stage.Kernel == p.kernel &&
		labelMatch
}

func checkLabels(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}

func timeDiff(t time.Time, start time.Time) time.Duration {
	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}
	return d
}
