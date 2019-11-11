// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("region", instance.Region).
		WithField("image", instance.Image).
		WithField("size", instance.Size).
		WithField("name", instance.Name)

	logger.Debugln("deleting instance")

	_ = p.deleteFloatingIps(instance)

	err := servers.Delete(p.computeClient, instance.ID).ExtractErr()
	if err == nil {
		logger.Debugln("instance deleted")
		return nil
	}

	logger.WithError(err).
		Errorln("deleting instance failed, attempting to force")

	err = servers.ForceDelete(p.computeClient, instance.ID).ExtractErr()

	if err == nil {
		logger.Debugln("instance deleted")
		return nil
	}

	logger.WithError(err).
		Errorln("force-deleting instance failed")

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
