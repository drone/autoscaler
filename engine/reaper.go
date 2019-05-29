// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/drone/autoscaler"

	"github.com/rs/zerolog/log"
)

// this is a feature flag that can be used to enable
// experimental reaping of errored instances.
var enableReaper = false

func init() {
	enableReaper, _ = strconv.ParseBool(
		os.Getenv("DRONE_ENABLE_REAPER"),
	)
}

//
// The reaper looks for and removes errored instances. The
// collector, on the other hand, is responsible for garbage
// collecting running instances that are no longer required.
//
// Note that I am open to using a more descriptive name if
// anyone has a better suggestion.
//

type reaper struct {
	wg sync.WaitGroup

	servers  autoscaler.ServerStore
	provider autoscaler.Provider
}

func (r *reaper) Reap(ctx context.Context) error {
	if !enableReaper {
		return nil
	}

	servers, err := r.servers.ListState(ctx, autoscaler.StateError)
	if err != nil {
		return err
	}

	for _, server := range servers {
		r.wg.Add(1)
		go func(server *autoscaler.Server) {
			r.reap(ctx, server)
			r.wg.Done()
		}(server)
	}
	return nil
}

func (r *reaper) reap(ctx context.Context, server *autoscaler.Server) error {
	logger := log.Ctx(ctx)
	logger.Debug().
		Str("state", "error").
		Str("server", server.Name).
		Msg("inspecting failed server")

	// if the server ID is an empty string it indicates
	// the server was never provisioned, but still has an
	// entry in the database. In this case, we can simply
	// delete the database entry
	if server.ID == "" {
		logger.Info().
			Str("state", "error").
			Str("server", server.Name).
			Msg("server never provisioned. nothing to destroy")
	} else {
		logger.Info().
			Str("state", "error").
			Str("server", server.Name).
			Msg("destroy provisioned server")

		in := &autoscaler.Instance{
			ID:       server.ID,
			Provider: server.Provider,
			Name:     server.Name,
			Address:  server.Address,
			Region:   server.Region,
			Image:    server.Image,
			Size:     server.Size,
		}

		err := r.provider.Destroy(ctx, in)
		// TODO implement ErrInstanceNotFound in Google driver
		// TODO implement ErrInstanceNotFound in Hetzner driver
		// TODO implement ErrInstanceNotFound in Packet driver
		if err == autoscaler.ErrInstanceNotFound {
			logger.Info().
				Str("state", "error").
				Str("server", server.Name).
				Msg("server no longer exists. nothing to destroy")

			// this accounts for the fact that the server can be
			// manually terminated outside of the autoscaler. In
			// this case the reaper continues and updates the
			// server state to stopped (below)

		} else if err != nil {
			logger.Error().Err(err).
				Str("state", "error").
				Str("server", server.Name).
				Msg("cannot destroy server")
			return err
		}
	}

	server.Stopped = time.Now().Unix()
	server.State = autoscaler.StateStopped
	return r.servers.Update(ctx, server)
}
