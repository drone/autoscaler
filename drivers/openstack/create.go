// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"bytes"
	"context"

	"github.com/drone/autoscaler"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/rs/zerolog/log"
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
	// Make a floating ip to attach.
	ip, err := floatingips.Create(p.computeClient, floatingips.CreateOpts{
		Pool: p.pool,
	}).Extract()
	if err != nil {
		return nil, err
	}

	serverCreateOpts := servers.CreateOpts{
		Name:           opts.Name,
		ImageName:      p.image,
		FlavorName:     p.flavor,
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
		floatingips.Delete(p.computeClient, ip.ID)
		return nil, err
	}
	logger := log.Ctx(ctx).With().
		Str("region", p.region).
		Str("image", p.image).
		Str("sizes", p.flavor).
		Str("name", opts.Name).
		Logger()

	err = servers.WaitForStatus(p.computeClient, server.ID, "ACTIVE", 300)
	if err != nil {
		return nil, err
	}
	floatingips.AssociateInstance(p.computeClient, server.ID, floatingips.AssociateOpts{
		FloatingIP: ip.IP,
	})

	logger.Debug().
		Msg("instance create")

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderOpenStack,
		ID:       server.ID,
		Name:     server.Name,
		Region:   p.region,
		Address:  ip.IP,
		Image:    p.image,
		Size:     p.flavor,
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("instance network ready")

	return instance, nil
}
