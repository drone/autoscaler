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
			logger.WithField("error", err).
				WithField("server", server.Name).
				Errorln("unexpected panic")
		}
	}()

	// if the server was never created there is nothing
	// to terminate, so we can just set the agent state
	// to term
	if server.ID == "" {
		logger.WithField("server", server.Name).
			Debugln("server never provisioned. nothing to stop")

		server.Stopped = time.Now().Unix()
		server.State = autoscaler.StateStopped

		err := c.servers.Update(ctx, server)
		if err != nil {
			logger.WithError(err).
				WithField("server", server.Name).
				Errorln("cannot update server state")
		} else {
			logger.WithField("server", server.Name).
				Debugln("updated server state to stopped")
		}
		return err
	}

	// first we need to gracefully shutdown the runner so
	// that in-progress pipelines can complete. They will
	// have up to 60 minutes to complete before being
	// force-killed.
	if server.Address != "" {
		client, closer, err := c.client(server)
		if closer != nil {
			defer closer.Close()
		}
		if err != nil {
			return err
		}

		logger.WithField("server", server.Name).
			Debugln("stopping the agent")

		ctxStop, cancel := context.WithTimeout(ctx, time.Hour)
		defer cancel()

		// 1 minute offset between docker stop timeout and
		// the context timeout.
		timeout := time.Hour - time.Minute
		err = client.ContainerStop(ctxStop, "agent", &timeout)
		if err != nil {
			logger.WithError(err).
				WithField("server", server.Name).
				Errorln("cannot stop the agent")
		} else {
			logger.WithField("server", server.Name).
				Debugln("stopped the agent")
		}
	}

	// next we need to terminate the remote instance (e.g. in aws).
	// It is possible the server was terminated out-of-band in which
	// case there is nothing to terminate.

	in := &autoscaler.Instance{
		ID:       server.ID,
		Provider: server.Provider,
		Name:     server.Name,
		Address:  server.Address,
		Region:   server.Region,
		Image:    server.Image,
		Size:     server.Size,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	err := c.provider.Destroy(ctx, in)
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
