// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

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

	// select random zone from the list
	zone := p.zones[rand.Intn(len(p.zones))]

	logger := logger.FromContext(ctx).
		WithField("zone", zone).
		WithField("image", p.image).
		WithField("size", p.size).
		WithField("name", opts.Name)

	logger.Debugln("instance insert")

	networkConfig := []*compute.AccessConfig{}

	if !p.privateIP {
		networkConfig = []*compute.AccessConfig{
			{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			},
		}
	}

    auto_restart := true
    on_host_maintenance := "MIGRATE"
	if p.provisioningModel == "SPOT" {
	    auto_restart = false
	    on_host_maintenance = "TERMINATE"
	}

	in := &compute.Instance{
		Name:           name,
		Zone:           fmt.Sprintf("projects/%s/zones/%s", p.project, zone),
		MinCpuPlatform: "Automatic",
		MachineType:    fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", p.project, zone, p.size),
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   p.userdataKey,
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
					DiskType:    fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", p.project, zone, p.diskType),
					DiskSizeGb:  p.diskSize,
				},
			},
		},
		CanIpForward: false,
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network:       p.network,
				Subnetwork:    p.subnetwork,
				AccessConfigs: networkConfig,
			},
		},
		Labels: p.labels,
		Scheduling: &compute.Scheduling{
			Preemptible:       false,
			ProvisioningModel: p.provisioningModel,
			OnHostMaintenance: on_host_maintenance,
			AutomaticRestart:  googleapi.Bool(auto_restart),
		},
		DeletionProtection: false,
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Scopes: p.scopes,
				Email:  p.serviceAccountEmail,
			},
		},
	}

	op, err := p.service.Instances.Insert(p.project, zone, in).Do()
	if err != nil {
		logger.WithError(err).
			Errorln("instance insert failed")
		return nil, err
	}

	logger.Debugln("pending instance insert operation")

	err = p.waitZoneOperation(ctx, op.Name, zone)
	if err != nil {
		logger.WithError(err).
			Errorln("instance insert operation failed")
		return nil, err
	}

	logger.Debugln("instance insert operation complete")

	resp, err := p.service.Instances.Get(p.project, zone, name).Do()
	if err != nil {
		logger.WithError(err).
			Errorln("cannot get instance details")
		return nil, err
	}

	address := resp.NetworkInterfaces[0].NetworkIP

	if !p.privateIP {
		address = resp.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}

	instance := &autoscaler.Instance{
		Provider:            autoscaler.ProviderGoogle,
		ID:                  name,
		Name:                opts.Name,
		Image:               p.image,
		Region:              zone,
		Size:                p.size,
		Address:             address,
		ServiceAccountEmail: p.serviceAccountEmail,
		Scopes:              p.scopes,
	}

	logger.
		WithField("name", instance.Name).
		WithField("ip", instance.Address).
		Debugln("instance inserted")

	return instance, nil
}
