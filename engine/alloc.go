// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/engine/certs"

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
		server.State = autoscaler.StateCreating
		err = a.servers.Update(ctx, server)
		if err != nil {
			logger.Error().
				Err(err).
				Str("server", server.Name).
				Str("state", "creating").
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

	ca, err := certs.GenerateCA()
	if err != nil {
		return err
	}

	cert, err := certs.GenerateCert(server.Name, ca)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	opts := autoscaler.InstanceCreateOpts{
		Name:    server.Name,
		CAKey:   ca.Key,
		CACert:  ca.Cert,
		TLSKey:  cert.Key,
		TLSCert: cert.Cert,
	}

	instance, err := a.provider.Create(ctx, opts)
	if err != nil {
		log.Ctx(ctx).Error().
			Err(err).
			Str("server", server.Name).
			Msg("failed to provision server")

		server.Error = err.Error()
		server.State = autoscaler.StateError
	} else {
		logger.Debug().
			Str("server", server.Name).
			Msg("provisioned server")

		server.State = autoscaler.StateCreated
	}
	if instance != nil {
		server.ID = instance.ID
		server.Address = instance.Address
		server.Image = instance.Image
		server.Provider = instance.Provider
		server.Region = instance.Region
		server.Size = instance.Size
		server.CACert = opts.CACert
		server.CAKey = opts.CAKey
		server.TLSCert = opts.TLSCert
		server.TLSKey = opts.TLSKey
		server.Started = time.Now().Unix()
	}

	err = a.servers.Update(ctx, server)
	if err != nil {
		logger.Error().
			Err(err).
			Str("server", server.Name).
			Msg("failed to update server state")
		return err
	}

	return nil
}
