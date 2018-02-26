// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scaler

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/drone-go/drone"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Scaler represents the built-in auto-scaler.
type Scaler struct {
	mu sync.Mutex

	paused   bool
	Client   drone.Client
	Config   config.Config
	Servers  autoscaler.ServerStore
	Provider autoscaler.Provider
}

// Pause paueses the scaler.
func (s *Scaler) Pause() {
	s.mu.Lock()
	s.paused = true
	s.mu.Unlock()
}

// Paused returns true if scaling is paused.
func (s *Scaler) Paused() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paused
}

// Resume resumes the scaler.
func (s *Scaler) Resume() {
	s.mu.Lock()
	s.paused = false
	s.mu.Unlock()
}

// Scale execute the autoscaling algorithm.
func (s *Scaler) Scale(ctx context.Context) error {
	logger := log.Ctx(ctx)

	pending, running, err := s.count(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("Error fetching queue details")
		return err
	}
	capacity, servers, err := s.capacity(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("Error calculating server capacity")
		return err
	}

	ctx = logger.With().
		Int("min-pool", s.Config.Pool.Min).
		Int("max-pool", s.Config.Pool.Max).
		Int("capacity", capacity).
		Int("pending", pending).
		Int("running", running).
		Logger().WithContext(ctx)

	free := capacity - running
	need := requiredCapacity(pending, free, s.Config.Agent.Concurrency)
	if need > 0 {
		// do not increase the pool beyond the max size.
		if servers+need >= s.Config.Pool.Max {
			need = s.Config.Pool.Max - servers
		}
		if need < 0 {
			need = 0
		}
		return s.provision(ctx, need)
	}
	return s.collect(ctx)
}

// provision provisions n new servers.
func (s *Scaler) provision(ctx context.Context, n int) error {
	logger := log.Ctx(ctx)
	logger.Debug().
		Int("create", n).
		Msgf("create %d instances", n)

	var g errgroup.Group
	for i := 0; i < n; i++ {
		g.Go(func() error {
			return s.provisionOne(ctx)
		})
	}
	return g.Wait()
}

// provisionOne provisions a new server.
func (s *Scaler) provisionOne(ctx context.Context) error {
	opts := autoscaler.NewServerOpts("agent", s.Config.Agent.Concurrency)
	server, err := s.Provider.Create(ctx, opts)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).
			Msg("failed to provision instance")
		return err
	}
	return s.Servers.Create(ctx, server)
}

// collect garbage collects servers that can be released
// due to inactivity.
func (s *Scaler) collect(ctx context.Context) error {
	logger := log.Ctx(ctx)

	servers, err := s.Servers.List(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("failed to fetch server list")
		return err
	}

	// do not garbage collect any servers if the pool
	// is at or below the minimum limit.
	if len(servers) <= s.Config.Pool.Min {
		logger.Debug().
			Msg("minimum server capacity")
		return nil
	}

	// sort by created date. older servers are retired first.
	sort.Sort(sort.Reverse(autoscaler.ByCreated(servers)))

	busy, err := s.listBusy(ctx)
	if err != nil {
		logger.Error().Err(err).
			Msg("failed to fetch busy list")
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
		if time.Now().Before(time.Unix(server.Created, 0).Add(s.Config.Pool.MinAge)) {
			logger.Debug().
				Str("server", server.Name).
				TimeDiff("age", time.Now(), time.Unix(server.Created, 0)).
				Dur("min-age", s.Config.Pool.MinAge).
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
		return nil
	}

	// we need to make sure  the count of idle servers that
	// are retired is > than the min pool size.
	if keep := len(servers) - len(idle); keep < s.Config.Pool.Min {
		idle = idle[:len(servers)-s.Config.Pool.Min-keep]
	}

	var g errgroup.Group
	for _, server := range idle {
		var aserver = server
		g.Go(func() error {
			return s.collectOne(ctx, aserver)
		})
	}
	return g.Wait()
}

// collectOne garbage collects a server.
func (s *Scaler) collectOne(ctx context.Context, server *autoscaler.Server) error {
	logger := log.Ctx(ctx)
	logger.Debug().
		Str("server", server.Name).
		Msg("destroying server")

	err := s.Provider.Destroy(ctx, server)
	if err != nil {
		logger.Error().
			Str("server", server.Name).
			Msg("failed to destroy server")
		// TODO: flag the server as in an invalid state.
		return err
	}
	logger.Debug().
		Str("server", server.Name).
		Msg("destroyed server")
	return s.Servers.Delete(ctx, server)
}

// helper function returns a list of busy servers.
func (s *Scaler) listBusy(ctx context.Context) (map[string]struct{}, error) {
	busy := map[string]struct{}{}
	builds, err := s.Client.BuildQueue()
	if err != nil {
		return busy, err
	}
	for _, build := range builds {
		if build.Status != drone.StatusRunning {
			continue
		}
		build, err := s.Client.Build(build.Owner, build.Name, build.Number)
		if err != nil {
			return busy, err
		}
		for _, proc := range build.Procs {
			busy[proc.Machine] = struct{}{}
		}
	}
	return busy, nil
}

// helper function returns our current capacity.
func (s *Scaler) capacity(ctx context.Context) (capacity, count int, err error) {
	servers, err := s.Servers.List(ctx)
	if err != nil {
		return capacity, count, err
	}
	for _, server := range servers {
		count++
		capacity += server.Capacity
	}
	return
}

// helper function returns the number of pending and
// running builds in the remote Drone installation.
func (s *Scaler) count(ctx context.Context) (pending, running int, err error) {
	activity, err := s.Client.BuildQueue()
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

func requiredCapacity(pending, available, concurrency int) int {
	diff := pending - available
	if diff <= 0 {
		return 0
	}
	more := int(math.Ceil(float64(diff) / float64(concurrency)))
	return more
}

// // helper function returns true if the server is idle.
// func (s *Scaler) checkIdle(ctx context.Context, server *autoscaler.Server) bool {
// 	state := struct {
// 		Polling int `json:"polling_count"`
// 		Running int `json:"running_count"`
// 	}{}
// 	res, err := http.Get("http://" + server.Address + "/varz")
// 	if err != nil {
// 		return false
// 	}
// 	defer res.Body.Close()
// 	err = json.NewDecoder(res.Body).Decode(&state)
// 	if err != nil {
// 		return false
// 	}
// 	return state.Running > 0
// }
