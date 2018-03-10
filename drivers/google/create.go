// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"fmt"

	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"

	"google.golang.org/api/compute/v1"
)

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	logger := log.Ctx(ctx).With().
		Str("zone", p.zone).
		Str("image", p.image).
		Str("size", p.size).
		Str("name", opts.Name).
		Logger()

	logger.Debug().
		Msg("instance insert")

	in := &compute.Instance{
		Name:        opts.Name,
		Description: "drone agent",
		MachineType: "", // TODO setup the machine type
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				Type:       "PERSISTENT",
				Mode:       "READ_WRITE",
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{Network: p.network},
		},
		Tags: &compute.Tags{
			Items: p.tags,
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Scopes: p.scopes,
				Email:  "default",
			},
		},
	}

	// TODO setup the disk

	op, err := p.service.Instances.Insert(p.project, p.zone, in).Context(ctx).Do()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("instance insert failed")
		return nil, err
	}

	logger.Debug().
		Msg("pending instance insert operation")

	err = p.waitZoneOperation(ctx, op.Name)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("instance insert operation failed")
		return nil, err
	}

	logger.Debug().
		Msg("instance insert operation complete")

	resp, err := p.service.Instances.Get(p.project, p.zone, opts.Name).Do()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot get instance details")
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderGoogle,
		ID:       fmt.Sprint(resp.Id),
		Name:     resp.Name,
		Image:    p.image,
		Region:   p.zone,
		Size:     p.size,
		Address:  resp.NetworkInterfaces[0].AccessConfigs[0].NatIP,
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("instance inserted")

	return instance, nil
}
