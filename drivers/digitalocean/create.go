// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package digitalocean

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/digitalocean/godo"
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

	req := &godo.DropletCreateRequest{
		Name:              opts.Name,
		Region:            p.region,
		Size:              p.size,
		Tags:              p.tags,
		IPv6:              false,
		PrivateNetworking: p.privateIP,
		UserData:          buf.String(),

		VPCUUID: 		   p.vpc,

		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: p.key},
		},
		Image: godo.DropletCreateImage{
			Slug: p.image,
		},
	}

	logger := logger.FromContext(ctx).
		WithField("region", req.Region).
		WithField("image", req.Image.Slug).
		WithField("size", req.Size).
		WithField("name", req.Name)

	logger.Debugln("instance create")

	client := newClient(ctx, p.token)
	droplet, _, err := client.Droplets.Create(ctx, req)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot create instance")
		return nil, err
	}

	if p.firewall != "" {
		_, err := client.Firewalls.AddDroplets(ctx, p.firewall, droplet.ID)
		if err != nil {
			logger.WithError(err).
				Errorln("cannot assign instance to firewall")
			return nil, err
		}
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderDigitalOcean,
		ID:       strconv.Itoa(droplet.ID),
		Name:     droplet.Name,
		Size:     req.Size,
		Region:   req.Region,
		Image:    req.Image.Slug,
	}

	logger.WithField("name", instance.Name).
		Infoln("instance created")

	// poll the digitalocean endpoint for server updates
	// and exit when a network address is allocated.
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.WithField("name", instance.Name).
				Debugln("cannot ascertain network")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute

			logger.WithField("name", instance.Name).
				Debugln("find instance network")

			droplet, _, err = client.Droplets.Get(ctx, droplet.ID)
			if err != nil {
				logger.WithError(err).
					Errorln("cannot find instance")
				return instance, err
			}

			for _, network := range droplet.Networks.V4 {
				if network.Type == "public" {
					instance.Address = network.IPAddress
				}
			}

			if instance.Address != "" {
				break poller
			}
		}
	}

	logger.
		WithField("name", instance.Name).
		WithField("ip", instance.Address).
		Debugln("instance network ready")

	return instance, nil
}
