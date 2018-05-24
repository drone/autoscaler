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

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("region", instance.Region).
		Str("image", instance.Image).
		Str("size", instance.Size).
		Str("name", instance.Name).
		Logger()

	client := newClient(ctx, p.token)
	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	_, res, err := client.Droplets.Get(ctx, id)
	if err != nil && res.StatusCode == 404 {
		logger.Warn().
			Err(err).
			Msg("droplet does not exist")
		return autoscaler.ErrInstanceNotFound
	} else if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot find droplet")
		return err
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
