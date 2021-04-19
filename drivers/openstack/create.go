// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"bytes"
	"context"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
)

// Create creates an OpenStack instance
func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	p.init.Do(func() {
		_ = p.setup(ctx)
	})

	buf := new(bytes.Buffer)
	err := p.userdata.Execute(buf, &opts)
	if err != nil {
		return nil, err
	}

	logger := logger.FromContext(ctx).
		WithField("region", p.region).
		WithField("image", p.image).
		WithField("flavor", p.flavor).
		WithField("network", p.network).
		WithField("pool", p.pool).
		WithField("name", opts.Name)

	logger.Debugln("instance create")

	nets := make([]servers.Network, 0)

	if p.network != "" {
		network, err := networks.Get(p.networkClient, p.network).Extract()
		if err != nil {
			logger.WithError(err).
				Debugln("failed to find network")
			return nil, err
		}

		nets = append(nets, servers.Network{
			UUID: network.ID,
		})
	}

	serverCreateOpts := servers.CreateOpts{
		Name:           opts.Name,
		ImageRef:       p.image,
		FlavorRef:      p.flavor,
		Networks:       nets,
		UserData:       buf.Bytes(),
		ServiceClient:  p.computeClient,
		Metadata:       p.metadata,
		SecurityGroups: p.groups,
	}
	createOpts := keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           p.key,
	}
	server, err := servers.Create(p.computeClient, createOpts).Extract()
	if err != nil {
		logger.WithError(err).
			Debugln("failed to create server")
		return nil, err
	}

	err = servers.WaitForStatus(p.computeClient, server.ID, "ACTIVE", 300)
	if err != nil {
		logger.WithError(err).
			Debugln("failed waiting for server")
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderOpenStack,
		ID:       server.ID,
		Name:     server.Name,
		Region:   p.region,
		Image:    p.image,
		Size:     p.flavor,
	}

	if p.network != "" {
		network, err := networks.Get(p.networkClient, p.network).Extract()
		if err != nil {
			logger.WithError(err).
				Debugln("failed to find network")
			return nil, err
		}

		if err := servers.ListAddresses(p.computeClient, server.ID).EachPage(func(page pagination.Page) (bool, error) {
			result, err := servers.ExtractAddresses(page)
			if err != nil {
				return false, err
			}

			for name, addresses := range result {
				if name == network.Name {
					for _, address := range addresses {
						instance.Address = address.Address
						return true, nil
					}
				}

			}

			return false, nil
		}); err != nil {
			logger.WithError(err).
				Debugln("failed to fetch address")
			return nil, err
		}
	}

	if p.pool != "" {
		ip, err := floatingips.Create(p.computeClient, floatingips.CreateOpts{
			Pool: p.pool,
		}).Extract()
		if err != nil {
			logger.WithError(err).
				Debugln("failed to create floating ip")
			return nil, err
		}

		if err := floatingips.AssociateInstance(p.computeClient, server.ID, floatingips.AssociateOpts{
			FloatingIP: ip.IP,
		}).ExtractErr(); err != nil {
			logger.WithError(err).
				Debugln("failed to associate floating ip")
			return nil, err
		}

		instance.Address = ip.IP
	}

	logger.
		WithField("name", instance.Name).
		WithField("ip", instance.Address).
		Debugln("instance network ready")

	return instance, nil
}
