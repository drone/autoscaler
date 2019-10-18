// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
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

	name := strings.ToLower(opts.Name)

	logger := log.Ctx(ctx).With().
		Str("zone", p.zone).
		Str("image", p.image).
		Str("size", p.size).
		Str("name", opts.Name).
		Logger()

	logger.Debug().
		Msg("instance insert")

	in := &compute.Instance{
		Name:           name,
		Zone:           fmt.Sprintf("projects/%s/zones/%s", p.project, p.zone),
		MinCpuPlatform: "Automatic",
		MachineType:    fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", p.project, p.zone, p.size),
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "user-data",
					Value: googleapi.String(buf.String()),
				},
			},
		},
		Tags: &compute.Tags{
			Items: p.tags,
		},
		Disks: []*compute.AttachedDisk{
			{
				Type:       "PERSISTENT",
				Boot:       true,
				Mode:       "READ_WRITE",
				AutoDelete: true,
				DeviceName: name,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s", p.image),
					DiskType:    fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", p.project, p.zone, p.diskType),
					DiskSizeGb:  p.diskSize,
				},
			},
		},
		CanIpForward: false,
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: p.network,
				AccessConfigs: []*compute.AccessConfig{
					{
						Name: "External NAT",
						Type: "ONE_TO_ONE_NAT",
					},
				},
			},
		},
		Labels: p.labels,
		Scheduling: &compute.Scheduling{
			Preemptible:       false,
			OnHostMaintenance: "MIGRATE",
			AutomaticRestart:  googleapi.Bool(true),
		},
		DeletionProtection: false,
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Scopes: p.scopes,
				Email:  "default",
			},
		},
	}

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

	resp, err := p.service.Instances.Get(p.project, p.zone, name).Do()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot get instance details")
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderGoogle,
		ID:       name,
		Name:     opts.Name,
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
