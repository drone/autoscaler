// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"strconv"

	"github.com/drone/autoscaler"

	"github.com/rs/zerolog/log"
)

// Destroy destroyes the DigitalOcean instance.
func (p *Provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("region", instance.Region).
		Str("image", instance.Image).
		Str("size", instance.Size).
		Str("name", instance.Name).
		Logger()

	client := newClient(ctx, p.config.DigitalOcean.Token)
	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	logger.Debug().
		Msg("teardown droplet")

	_, err = p.Provider.Execute(ctx, instance, teardownScript)
	if err != nil {
		// if we cannot gracefully shutdown the agent we should
		// still continue and destroy the droplet. I think.
		logger.Error().
			Err(err).
			Msg("teardown failed")

		// TODO(bradrydzewski) we should snapshot the error logs
	}

	logger.Debug().
		Msg("deleting droplet")

	_, err = client.Droplets.Delete(ctx, id)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("deleting droplet failed")
		return err
	}

	logger.Debug().
		Msg("droplet deleted")

	return nil
}

var teardownScript = `
set -x;

docker stop -t 3600 agent
docker ps -a
`
