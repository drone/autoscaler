// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"bytes"
	"context"
	"strconv"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/hetznercloud/hcloud-go/hcloud"
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

	req := hcloud.ServerCreateOpts{
		Name:     opts.Name,
		UserData: buf.String(),
		ServerType: &hcloud.ServerType{
			Name: p.serverType,
		},
		Image: &hcloud.Image{
			Name: p.image,
		},
		SSHKeys: []*hcloud.SSHKey{
			{
				ID: p.key,
			},
		},
	}

	datacenter := "unknown"

	if p.datacenter != "" {
		req.Datacenter = &hcloud.Datacenter{
			Name: p.datacenter,
		}

		datacenter = p.datacenter
	}

	logger := logger.FromContext(ctx).
		WithField("datacenter", datacenter).
		WithField("image", req.Image.Name).
		WithField("serverType", req.ServerType.Name).
		WithField("name", req.Name)

	if p.firewall != 0 {
		// Check provided FirewallID exists
		firewall, _, err := p.client.Firewall.GetByID(ctx, p.firewall)
		if err != nil {
			logger.WithError(err).
				Errorln("failed to retrieve firewall")
			return nil, err
		}
		if firewall == nil {
			logger.WithError(err).
				WithField("id", p.firewall).Errorf("cannot find firewall")
			return nil, err
		}
		_ = append(req.Firewalls, &hcloud.ServerCreateFirewall{
			Firewall: hcloud.Firewall{
				ID: firewall.ID,
			},
		})
	}
	logger.Debugln("instance create")

	resp, _, err := p.client.Server.Create(ctx, req)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot create instance")
		return nil, err
	}

	logger.
		WithField("name", req.Name).
		Infoln("instance created")

	return &autoscaler.Instance{
		Provider: autoscaler.ProviderHetznerCloud,
		ID:       strconv.Itoa(resp.Server.ID),
		Name:     resp.Server.Name,
		Address:  resp.Server.PublicNet.IPv4.IP.String(),
		Size:     req.ServerType.Name,
		Region:   datacenter,
		Image:    req.Image.Name,
	}, nil
}
