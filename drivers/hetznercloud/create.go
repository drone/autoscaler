// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
		Labels: p.labels,
	}

	for _, netName := range []string{
		p.network,
		p.private,
	} {
		if err := func() error {
			if netName == "" {
				return nil
			}
			for _, oldNet := range req.Networks {
				if oldNet.Name == netName {
					return nil
				}
			}
			net, _, err := p.client.Network.GetByName(ctx, netName)
			if err != nil {
				return err
			} else if net == nil {
				return errors.New(fmt.Sprintf("Network %s not found.", netName))
			}
			req.Networks = append(req.Networks, net)
			return nil
		}(); err != nil {
			return nil, err
		}
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

	logger.Debugln("instance create")

	resp, _, err := p.client.Server.Create(ctx, req)
	if err != nil {
		logger.WithError(err).
			Errorln("cannot create instance")
		return nil, err
	}

	logger.
		WithField("name", req.Name).
		WithField("full", req).
		Infoln("instance created")

	var ip string
	if p.private != "" {
		_, errC := p.client.Action.WatchOverallProgress(ctx, resp.NextActions)
		if err := <-errC; err != nil {
			logger.WithError(err).
				Errorln("Instance failed to start.")
			return nil, err
		}
		s, _, err := p.client.Server.GetByID(context.Background(), resp.Server.ID)
		if err != nil {
			logger.WithError(err).
				Errorln("Failed to retrieve created server.")
			return nil, err
		}
		// for some reason GetByID returns Networks without Name set, so we need to get the ID of the private network
		var id int
		for _, net := range req.Networks {
			if net.Name == p.private {
				id = net.ID
				break
			}
		}
		for _, net := range s.PrivateNet {
			if net.Network.ID == id {
				ip = net.IP.String()
				break
			}
		}
	} else {
		ip = resp.Server.PublicNet.IPv4.IP.String()
	}
	if ip == "" {
		return nil, errors.New("Instance address not set (Private network not found on instance?).")
	}

	return &autoscaler.Instance{
		Provider: autoscaler.ProviderHetznerCloud,
		ID:       strconv.Itoa(resp.Server.ID),
		Name:     resp.Server.Name,
		Address:  ip,
		Size:     req.ServerType.Name,
		Region:   datacenter,
		Image:    req.Image.Name,
	}, nil
}
