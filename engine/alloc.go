// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/drone/autoscaler"

	"github.com/rs/zerolog/log"
)

type allocator struct {
	wg sync.WaitGroup

	servers  autoscaler.ServerStore
	provider autoscaler.Provider
}

func (a *allocator) Allocate(ctx context.Context) error {
	logger := log.Ctx(ctx)

	servers, err := a.servers.ListState(ctx, autoscaler.StatePending)
	if err != nil {
		return err
	}

	for _, server := range servers {
		server.State = autoscaler.StateStaging
		err = a.servers.Update(ctx, server)
		if err != nil {
			logger.Error().
				Err(err).
				Str("server", server.Name).
				Str("state", "staging").
				Msg("failed to update server state")
			return err
		}

		a.wg.Add(1)
		go func(server *autoscaler.Server) {
			a.allocate(ctx, server)
			a.wg.Done()
		}(server)
	}
	return nil
}

func (a *allocator) allocate(ctx context.Context, server *autoscaler.Server) error {
	logger := log.Ctx(ctx)
	defer func() {
		if err := recover(); err != nil {
			logger.Error().
				Err(err.(error)).
				Str("server", server.Name).
				Msg("unexpected panic")
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	opt := autoscaler.InstanceCreateOpts{Name: server.Name}
	instance, err := a.provider.Create(ctx, opt)
	if err != nil {
		log.Ctx(ctx).Error().
			Err(err).
			Str("server", server.Name).
			Msg("failed to provision server")

		server.State = autoscaler.StateError
	} else {
		logger.Debug().
			Str("server", server.Name).
			Msg("provisioned server")

		server.State = autoscaler.StateRunning
	}
	if lerr, ok := err.(*autoscaler.InstanceError); ok {
		server.Error = string(lerr.Logs)
	}
	if instance != nil {
		server.ID = instance.ID
		server.Address = instance.Address
		server.Image = instance.Image
		server.Provider = instance.Provider
		server.Region = instance.Region
		server.Secret = instance.Secret
		server.Size = instance.Size
		server.Started = time.Now().Unix()
	}
	return a.servers.Update(ctx, server)
}
