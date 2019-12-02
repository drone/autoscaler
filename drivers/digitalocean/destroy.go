// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"strconv"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("region", instance.Region).
		WithField("image", instance.Image).
		WithField("size", instance.Size).
		WithField("name", instance.Name)

	client := newClient(ctx, p.token)
	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	_, res, err := client.Droplets.Get(ctx, id)
	if err != nil && res.StatusCode == 404 {
		logger.WithError(err).
			Warnln("droplet does not exist")
		return autoscaler.ErrInstanceNotFound
	} else if err != nil {
		logger.WithError(err).
			Errorln("cannot find droplet")
		return err
	}

	logger.Debugln("deleting droplet")

	_, err = client.Droplets.Delete(ctx, id)
	if err != nil {
		logger.WithError(err).
			Errorln("deleting droplet failed")
		return err
	}

	logger.Debugln("droplet deleted")

	return nil
}
