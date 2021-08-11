// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"bytes"
	"context"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	"strconv"

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
		logger.Infoln("firewall specified")

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

		logger.WithField("firewall", firewall.Name).Debugln("firewall found")
		req.Firewalls = append(req.Firewalls, &hcloud.ServerCreateFirewall{
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

	err = p.WatchActionProgress(ctx, resp.Action)
	if err != nil {
		return nil, err
	}

	logger.
		WithField("name", req.Name).
		Infoln("instance created")

	address := resp.Server.PublicNet.IPv4.IP.String()
	if p.network != 0 {
		logger.Infoln("network specified")

		network, _, err := p.client.Network.GetByID(ctx, p.network)
		if err != nil {
			logger.WithError(err).
				Errorln("failed to retrieve network")
			return nil, err
		}
		if network == nil {
			logger.WithError(err).
				WithField("id", p.network).Errorf("cannot find network")
			return nil, err
		}
		logger.WithField("network", network.Name).Debugln("network found")

		action, _, err := p.client.Server.AttachToNetwork(ctx, resp.Server, hcloud.ServerAttachToNetworkOpts{
			Network: network,
		})
		if err != nil {
			return nil, err
		}

		logger.Infoln("Waiting for IP to be assigned")
		err = p.WatchActionProgress(ctx, action)
		if err != nil {
			return nil, err
		}
		server, _, err := p.client.Server.GetByID(ctx, resp.Server.ID)
		if err != nil {
			logger.WithError(err).
				Errorln("cannot find instance")
			return nil, err
		}
		address = server.PrivateNet[0].IP.String()
	}

	logger.
		WithField("name", req.Name).
		Infoln("instance created in network")

	return &autoscaler.Instance{
		Provider: autoscaler.ProviderHetznerCloud,
		ID:       strconv.Itoa(resp.Server.ID),
		Name:     resp.Server.Name,
		Address:  address,
		Size:     req.ServerType.Name,
		Region:   datacenter,
		Image:    req.Image.Name,
	}, nil
}

func (p *provider) WatchActionProgress(ctx context.Context, action *hcloud.Action) error {
	progressCh, errCh := p.client.Action.WatchProgress(ctx, action)
outer:
	for {
		select {
		case err := <-errCh:
			return err
		case p := <-progressCh:
			if int64(p) == 100 {
				break outer
			}
		}
	}
	return nil
}
