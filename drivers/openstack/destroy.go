// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"context"
	"fmt"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("region", instance.Region).
		WithField("image", instance.Image).
		WithField("flavor", instance.Size).
		WithField("name", instance.Name)

	logger.Debugln("deleting instance")

	err := p.deleteFloatingIps(instance)
	if err != nil {
		logger.WithError(err).
			Debugln("failed to delete floating ips")

		return err
	}

	err = servers.Delete(p.computeClient, instance.ID).ExtractErr()
	if err == nil {
		logger.Debugln("instance deleted")
		return nil
	}

	if err.Error() == "Resource not found" {
		logger.WithError(err).
			Debugln("instance does not exist")
		return autoscaler.ErrInstanceNotFound
	}

	logger.WithError(err).
		Errorln("attempting to force delete")

	err = servers.ForceDelete(p.computeClient, instance.ID).ExtractErr()
	if err == nil {
		logger.Debugln("instance deleted")
		return nil
	}

	if err.Error() == "Resource not found" {
		logger.WithError(err).
			Debugln("instance does not exist")
		return autoscaler.ErrInstanceNotFound
	}

	logger.WithError(err).
		Errorln("force-deleting instance failed")

	return err
}

func (p *provider) deleteFloatingIps(instance *autoscaler.Instance) error {
	return floatingips.List(p.computeClient).EachPage(func(page pagination.Page) (bool, error) {
		ips, err := floatingips.ExtractFloatingIPs(page)
		if err != nil {
			return false, err
		}

		for _, ip := range ips {
			if ip.InstanceID == instance.ID {
				if err := floatingips.DisassociateInstance(p.computeClient, instance.ID, floatingips.DisassociateOpts{
					FloatingIP: ip.IP,
				}).ExtractErr(); err != nil {
					return false, fmt.Errorf("failed to disassociate floating ip: %s", err)
				}

				if err := floatingips.Delete(p.computeClient, ip.ID).ExtractErr(); err != nil {
					return false, fmt.Errorf("failed to delete floating ip: %s", err)
				}
			}
		}

		return true, nil
	})
}
