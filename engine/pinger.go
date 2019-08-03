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
	logger := log.Ctx(ctx).With().
		Str("ip", server.Address).
		Str("name", server.Name).
		Logger()

	client, closer, err := p.client(server)
	defer closer.Close()
	if err != nil {
		logger.Error().Err(err).
			Msg("cannot create docker client")
		return nil
	}

	// the system will attempt to ping the server a maximum of
	// five times, with a 1 minute timeout for each ping. If the
	// server cannot be reached, it will be placed in an error
	// state.

	for i := 0; i < 5; i++ {
		timeout, cancel := context.WithTimeout(ctx, time.Minute)
		_, err := client.Ping(timeout)
		cancel()
		if err == nil {
			logger.Debug().
				Str("state", "healthy").
				Msg("server ping successful")
			return nil
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
	return p.servers.Update(ctx, server)
}
