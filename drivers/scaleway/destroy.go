// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package scaleway

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (p *provider) Destroy(ctx context.Context, inst *autoscaler.Instance) error {
	p.init.Do(func() {
		p.setup(ctx)
	})

	logger := logger.FromContext(ctx).
		WithField("datacenter", inst.Region).
		WithField("image", inst.Image).
		WithField("size", inst.Size).
		WithField("name", inst.Name)

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
			logger.WithError(err).
				Errorln("cannot get server")
			return err
		}
	}

	// Issue "terminate" action, instead of DeleteServer, as terminate
	// cleans up volumes and IP addresses attached, too
	req := &instance.ServerActionRequest{
		ServerID: inst.ID,
		Action:   instance.ServerActionTerminate,
	}

	logger.Debugln("terminating server")

	_, err = api.ServerAction(req, scw.WithContext(ctx))

	if err != nil {
		logger.WithError(err).
			Errorln("terminating server failed")
		return err
	}

	logger.Infoln("server terminated")

	return err
}
