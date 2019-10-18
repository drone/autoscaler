// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"strconv"

	"github.com/drone/autoscaler"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/rs/zerolog/log"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("region", instance.Region).
		Str("image", instance.Image).
		Str("size", instance.Size).
		Str("name", instance.Name).
		Logger()

	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	logger.Debug().
		Msg("deleting instance")

	_, err = p.client.Server.Delete(ctx, &hcloud.Server{ID: id})
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
