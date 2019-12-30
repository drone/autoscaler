// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"strconv"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("region", instance.Region).
		WithField("image", instance.Image).
		WithField("size", instance.Size).
		WithField("name", instance.Name)

	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	logger.Debugln("deleting instance")

	_, err = p.client.Server.Delete(ctx, &hcloud.Server{ID: id})

	if err != nil {
		if err.Error() == "hcloud: server responded with status code 404" {
			logger.WithError(err).
				Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}

		logger.WithError(err).
			Errorln("deleting instance failed")
		return err
	}

	logger.Debugln("instance deleted")

	return nil
}
