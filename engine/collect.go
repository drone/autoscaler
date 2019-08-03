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

type collector struct {
	wg sync.WaitGroup

	servers  autoscaler.ServerStore
	provider autoscaler.Provider
	client   clientFunc
}

func (c *collector) Collect(ctx context.Context) error {
	logger := log.Ctx(ctx)

	servers, err := c.servers.ListState(ctx, autoscaler.StateShutdown)
	if err != nil {
		return err
	}

	for _, server := range servers {
		server.State = autoscaler.StateStopping
		err = c.servers.Update(ctx, server)
		if err != nil {
			logger.Error().
				Err(err).
				Str("server", server.Name).
				Str("state", "stopping").
				Msg("failed to update server state")
			return err
		}

		c.wg.Add(1)
		go func(server *autoscaler.Server) {
			c.collect(ctx, server)
			c.wg.Done()
		}(server)
	}
	return nil
}

func (c *collector) collect(ctx context.Context, server *autoscaler.Server) error {
	logger := log.Ctx(ctx)
	logger.Debug().
		Str("server", server.Name).
		Msg("destroying server")

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

	in := &autoscaler.Instance{
		ID:       server.ID,
		Provider: server.Provider,
		Name:     server.Name,
		Address:  server.Address,
		Region:   server.Region,
		Image:    server.Image,
		Size:     server.Size,
	}

	client, closer, err := c.client(server)
	defer closer.Close()
	if err != nil {
		return err
	}

	timeout := time.Hour * 60
	err = client.ContainerStop(ctx, "agent", &timeout)
	if err != nil {
		logger.Warn().Err(err).
			Str("server", server.Name).
			Msg("cannot stop the agent")
	}

	err = c.provider.Destroy(ctx, in)
	if err != nil {
		logger.Error().
			Str("server", server.Name).
			Msg("failed to destroy server")

		server.Error = err.Error()
		server.State = autoscaler.StateError
	} else {
		logger.Debug().
			Str("server", server.Name).
			Msg("destroyed server")

		server.Stopped = time.Now().Unix()
		server.State = autoscaler.StateStopped
	}

	return c.servers.Update(ctx, server)
}
