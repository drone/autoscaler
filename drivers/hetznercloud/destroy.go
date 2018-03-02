// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"strconv"

	"github.com/drone/autoscaler"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/rs/zerolog/log"
)

// Destroy destroyes the HetznerCloud instance.
func (p *Provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("type", instance.Size).
		Str("image", instance.Image).
		Str("datacenter", instance.Region).
		Str("name", instance.Name).
		Logger()

	id, err := strconv.Atoi(instance.ID)

	if err != nil {
		return err
	}

	logger.Debug().
		Msg("teardown server")

	// TODO: stop the agent via remote docker connection?

	// _, err = p.Provider.Execute(ctx, instance, teardownScript)
	// if err != nil {
	// 	// if we cannot gracefully shutdown the agent we should
	// 	// still continue and destroy the instance. I think.
	// 	logger.Error().
	// 		Err(err).
	// 		Msg("teardown failed")
	// }

	logger.Debug().
		Msg("deleting server")

	_, err = p.client().Server.Delete(ctx, &hcloud.Server{ID: id})

	if err != nil {
		logger.Error().
			Err(err).
			Msg("deleting instance failed")

		return err
	}

	logger.Debug().
		Msg("instance deleted")

	return nil
}

var teardownScript = `
set -x;
docker stop -t 3600 agent
docker ps -a
`
