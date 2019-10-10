// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scaleway

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (p *provider) Destroy(ctx context.Context, inst *autoscaler.Instance) error {
	p.init.Do(func() {
		p.setup(ctx)
	})

	logger := log.Ctx(ctx).With().
		Str("datacenter", inst.Region).
		Str("image", inst.Image).
		Str("size", inst.Size).
		Str("name", inst.Name).
		Logger()

	api := instance.NewAPI(p.client)

	srvReq := &instance.GetServerRequest{
		ServerID: inst.ID,
	}
	_, err := api.GetServer(srvReq, scw.WithContext(ctx))
	if err != nil {
		scwErr, ok := err.(*scw.ResponseError)
		if ok && scwErr.StatusCode == 404 {
			return autoscaler.ErrInstanceNotFound
		} else {
			logger.Error().
				Err(err).
				Msg("cannot get server")
			return err
		}
	}

	// Issue "terminate" action, instead of DeleteServer, as terminate
	// cleans up volumes and IP addresses attached, too
	req := &instance.ServerActionRequest{
		ServerID: inst.ID,
		Action:   instance.ServerActionTerminate,
	}

	logger.Debug().
		Msg("terminating server")

	_, err = api.ServerAction(req, scw.WithContext(ctx))

	if err != nil {
		logger.Error().
			Err(err).
			Msg("terminating server failed")
		return err
	}

	logger.Info().
		Msg("server terminated")

	return err
}
