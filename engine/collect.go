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

type collector struct {
	wg sync.WaitGroup

	servers  autoscaler.ServerStore
	provider autoscaler.Provider
	client   clientFunc
}

func (c *collector) Collect(ctx context.Context) error {
	logger := logger.FromContext(ctx)

	servers, err := c.servers.ListState(ctx, autoscaler.StateShutdown)
	if err != nil {
		return err
	}

	for _, server := range servers {
		server.State = autoscaler.StateStopping
		err = c.servers.Update(ctx, server)
		if err != nil {
			logger.WithError(err).
				WithField("server", server.Name).
				WithField("state", "stopping").
				Errorln("failed to update server state")
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
	logger := logger.FromContext(ctx)
	logger.WithField("server", server.Name).
		Debugln("destroying server")

	defer func() {
		if err := recover(); err != nil {
			logger.WithError(err.(error)).
				WithField("server", server.Name).
				Errorln("unexpected panic")
		}
	}()

	ctx, cancel := context.WithTimeout(ctx,
		time.Hour+time.Minute+time.Minute, // two minute buffer
	)
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
	if closer != nil {
		defer closer.Close()
	}
	if err != nil {
		return err
	}

	timeout := time.Hour
	err = client.ContainerStop(ctx, "agent", &timeout)
	if err != nil {
		logger.WithError(err).
			WithField("server", server.Name).
			Errorln("cannot stop the agent")
	}

	err = c.provider.Destroy(ctx, in)
	if err == autoscaler.ErrInstanceNotFound {
		logger.
			WithField("state", "error").
			WithField("server", server.Name).
			Infoln("server no longer exists. nothing to destroy")

		server.Stopped = time.Now().Unix()
		server.State = autoscaler.StateStopped
	} else if err != nil {
		logger.WithField("server", server.Name).
			Errorln("failed to destroy server")

		server.Error = err.Error()
		server.State = autoscaler.StateError
	} else {
		logger.WithField("server", server.Name).
			Debugln("destroyed server")

		server.Stopped = time.Now().Unix()
		server.State = autoscaler.StateStopped
	}

	err = c.servers.Update(ctx, server)
	if err != nil {
		logger.WithError(err).
			WithField("server", server.Name).
			Errorln("failed to update server state")
		return err
	}

	return nil
}
