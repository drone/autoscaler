// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import (
	"context"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("id", instance.ID).
		WithField("name", instance.Name).
		WithField("ip", instance.Address)

	logger.Debugln("terminate instance")

	server, err := p.client.Server(instance.ID)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot retrive instance details")
		return err
	}

	if len(server.CloudIPs) > 0 {

		logger.Debugln("unmap Cloud IP")

		cipID := server.CloudIPs[0].ID

		err := p.client.UnMapCloudIP(cipID)
		if err != nil {
			logger.WithError(err).
				Errorln("failed to unmap Cloud IP")
			return err
		}

		// wait for CIP to become unmapped
		interval := time.Second * 2
	poller:
		for {
			select {
			case <-ctx.Done():
				logger.WithField("ID", server.ID).
					Debugln("unmap Cloud IP deadline exceeded")

				return err
			case <-time.After(interval):
				cip, err := p.client.CloudIP(cipID)
				if err != nil {
					logger.WithError(err).
						Errorln("cannot retrive Cloud IP details")
					continue
				}

				if cip.Status == "unmapped" {
					break poller
				}
			}
		}

		err = p.client.DestroyCloudIP(cipID)
		if err != nil {
			logger.WithError(err).
				Errorln("failed to destroy Cloud IP")
			return err
		}
	}

	logger.Infoln("unmap Cloud IP success")

	err = p.client.DestroyServer(instance.ID)
	if err != nil {
		logger.WithError(err).
			Errorln("terminate instance failed")
		return err
	}

	logger.Infoln("terminate instance success")

	return nil
}
