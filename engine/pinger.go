// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
)

// this is a feature flag that can be used to enable
// experimental pinging and detection of zombie instances.
var enablePinger = false

func init() {
	enablePinger, _ = strconv.ParseBool(
		os.Getenv("DRONE_ENABLE_PINGER"),
	)
}

type pinger struct {
	wg sync.WaitGroup

	servers autoscaler.ServerStore
	client  clientFunc
}

func (p *pinger) Ping(ctx context.Context) error {
	if !enablePinger {
		return nil
	}

	servers, err := p.servers.ListState(ctx, autoscaler.StateRunning)
	if err != nil {
		return err
	}

	for _, server := range servers {
		p.wg.Add(1)
		go func(server *autoscaler.Server) {
			p.ping(ctx, server)
			p.wg.Done()
		}(server)
	}
	return nil
}

func (p *pinger) ping(ctx context.Context, server *autoscaler.Server) error {
	logger := logger.FromContext(ctx).
		WithField("ip", server.Address).
		WithField("name", server.Name)

	client, closer, err := p.client(server)
	if closer != nil {
		defer closer.Close()
	}
	if err != nil {
		logger.WithError(err).
			Errorln("cannot create docker client")
		return nil
	}

	// the system will attempt to ping the server a maximum of
	// five times, with a 1 minute timeout for each ping. If the
	// server cannot be reached, it will be placed in an error
	// state.

	const pingRetries = 5
	for i := 0; i < pingRetries; i++ {
		timeout, cancel := context.WithTimeout(ctx, time.Minute)
		_, err := client.Ping(timeout)
		cancel()
		if err == nil {
			logger.WithField("state", "healthy").
				WithField("server", server.Name).
				Debugln("server ping successful")
			return nil
		}
		if i < pingRetries-1 {
			logger.WithError(err).
				WithField("state", "maybe-unhealthy").
				WithField("server", server.Name).
				Debugln("server ping failed, will retry")
		} else {
			logger.WithError(err).
				WithField("state", "unhealthy").
				WithField("server", server.Name).
				Warnf("server failed to ping %d times, will try to mark as errored...", pingRetries)
		}
	}

	server, err = p.servers.Find(ctx, server.Name)
	if err != nil {
		// if the server no longer exists in the database
		// it is possible it was mutated by another goroutine.
		return err
	}

	if server.State != autoscaler.StateRunning {
		// if the server was mutated by another goroutine
		// we should exit without making any changes.
		return nil
	}

	server.Error = "Failed to ping the server"
	server.Stopped = time.Now().Unix()
	server.State = autoscaler.StateError
	err = p.servers.Update(ctx, server)
	if err != nil {
		logger.WithError(err).
			WithField("server", server.Name).
			WithField("state", "error").
			Errorln("failed to update server state")
		return err
	}

	return nil
}
