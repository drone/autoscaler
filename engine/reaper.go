// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
)

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
	enabled  bool
	interval time.Duration
}

func (r *reaper) Reap(ctx context.Context) error {
	if !r.enabled {
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
	logger := logger.FromContext(ctx)
	logger.
		WithField("state", "error").
		WithField("server", server.Name).
		Debugln("inspecting failed server")

	// if the server ID is an empty string it indicates
	// the server was never provisioned, but still has an
	// entry in the database. In this case, we can simply
	// delete the database entry
	if server.ID == "" {
		logger.
			WithField("state", "error").
			WithField("server", server.Name).
			Infoln("server never provisioned. nothing to destroy")
	} else {
		logger.
			WithField("state", "error").
			WithField("server", server.Name).
			Infoln("destroy provisioned server")

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
			logger.
				WithField("state", "error").
				WithField("server", server.Name).
				Infoln("server no longer exists. nothing to destroy")

			// this accounts for the fact that the server can be
			// manually terminated outside of the autoscaler. In
			// this case the reaper continues and updates the
			// server state to stopped (below)

		} else if err != nil {
			logger.WithError(err).
				WithField("state", "error").
				WithField("server", server.Name).
				Errorln("cannot destroy server")
			return err
		}
	}

	server.Stopped = time.Now().Unix()
	server.State = autoscaler.StateStopped
	err := r.servers.Update(ctx, server)
	if err != nil {
		logger.WithError(err).
			WithField("server", server.Name).
			WithField("state", "stopped").
			Errorln("failed to update server state")
		return err
	}

	return nil
}
