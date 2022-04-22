// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/brightbox/gobrightbox"
)

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	p.init.Do(func() {
		p.setup(ctx)
	})

	buf := new(bytes.Buffer)
	err := p.userdata.Execute(buf, &opts)
	if err != nil {
		return nil, err
	}

	var userdata = new(string)
	*userdata = base64.StdEncoding.EncodeToString(buf.Bytes())

	in := &gobrightbox.ServerOptions{
		Image:        p.image,
		Name:         &opts.Name,
		ServerType:   p.serverType,
		UserData:     userdata,
		ServerGroups: p.serverGroups,
	}

	logger := logger.FromContext(ctx).
		WithField("name", opts.Name).
		WithField("type", p.serverType).
		WithField("image", p.image).
		WithField("groups", p.serverGroups)

	logger.Debugln("instance create")

	server, err := p.client.CreateServer(in)
	if err != nil {
		logger.WithError(err).
			Errorln("instance create failed")
		return nil, err
	}

	// wait for the server to become active
	interval := time.Second * 10
poller:
	for {
		select {
		case <-ctx.Done():
			logger.WithField("ID", server.ID).
				Debugln("instance create deadline exceeded")

			return nil, ctx.Err()
		case <-time.After(interval):
			server, err := p.client.Server(server.ID)
			if err != nil {
				logger.WithError(err).
					Errorln("cannot get instance details")
				continue
			}

			if server.Status == "active" {
				break poller
			} else if server.Status == "failed" {
				err = errors.New("brightbox: new server entered 'failed' state")
				logger.WithError(err).
					Errorln("instance create failed")
				return nil, err
			}
		}
	}

	logger.Infoln("instance create success")

	cip := &gobrightbox.CloudIPOptions{
		Name: &opts.Name,
	}

	logger.Debugln("map Cloud IP")

	cloudip, err := p.client.CreateCloudIP(cip)
	if err != nil {
		logger.WithError(err).
			Errorln("failed to provision CloudIP")
		return nil, err
	}

	err = p.client.MapCloudIPtoServer(cloudip.ID, server.ID)
	if err != nil {
		logger.WithError(err).
			Errorln("failed to map CloudIP")
		return nil, err
	}

	logger.Infoln("map Cloud IP success")

	return &autoscaler.Instance{
		Provider: autoscaler.ProviderBrightbox,
		Region:   p.region,
		ID:       server.ID,
		Name:     server.Name,
		Address:  cloudip.PublicIPv4,
		Image:    in.Image,
		Size:     in.ServerType,
	}, nil
}
