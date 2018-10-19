// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/rs/zerolog/log"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("region", instance.Region).
		Str("image", instance.Image).
		Str("size", instance.Size).
		Str("name", instance.Name).
		Logger()

	logger.Debug().
		Msg("deleting instance")

	_ = p.deleteFloatingIps(instance)

	err := servers.Delete(p.computeClient, instance.ID).ExtractErr()
	if err == nil {
		logger.Debug().
			Msg("instance deleted")
		return nil
	}

	logger.Error().
		Err(err).
		Msg("deleting instance failed, attempting to force")

	err = servers.ForceDelete(p.computeClient, instance.ID).ExtractErr()

	if err == nil {
		logger.Debug().
			Msg("instance deleted")
		return nil
	}

	logger.Error().
		Err(err).
		Msg("force-deleting instance failed")

	return err
}

func (p *provider) deleteFloatingIps(instance *autoscaler.Instance) error {
	floatingips.DisassociateInstance(p.computeClient, instance.ID, floatingips.DisassociateOpts{
		FloatingIP: instance.Address,
	})
	// Remove our allocated ip from the pool.
	allPages, err := floatingips.List(p.computeClient).AllPages()
	ips, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return err
	}
	for _, fip := range ips {
		if fip.InstanceID == instance.ID {
			floatingips.Delete(p.computeClient, fip.ID)
		}
	}

	return nil
}
