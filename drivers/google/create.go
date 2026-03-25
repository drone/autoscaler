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
	"time"

	"github.com/avast/retry-go"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	"github.com/google/uuid"

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
				StackType:     p.stackType,
				AccessConfigs: networkConfig,
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
				Email:  p.serviceAccountEmail,
			},
		},
	}

	// Cannot add this in the same way as v4 access configs since the instance creation
	// fails if any v6 access configs are specified for an instance with IPV4_ONLY stack type
	if p.stackType == "IPV4_IPV6" {
		in.NetworkInterfaces[0].Ipv6AccessConfigs = []*compute.AccessConfig{
			{
				Name:        "external-ipv6",
				Type:        "DIRECT_IPV6",
				NetworkTier: "PREMIUM",
			},
		}
	}

	// Generate a UUID for idempotent retries. This ensures that if the request
	// must be retried, the server will ignore duplicate operations.
	// This currently does not persist in the database so if an autoscaler dies
	// during a retry loop the database may show stuck in creating
	requestID := uuid.New().String()

	var op *compute.Operation
	err = retry.Do(
		func() error {
			var insertErr error
			op, insertErr = p.service.Instances.Insert(p.project, zone, in).RequestId(requestID).Context(ctx).Do()
			if insertErr != nil {
				// Return transient errors for retry, non-transient as unrecoverable
				if isTransientError(insertErr) {
					return insertErr
				}
				return retry.Unrecoverable(insertErr)
			}
			return nil
		},
		retry.Attempts(5),
		retry.MaxDelay(time.Second*5),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			logger.WithField("attempt", n+1).
				WithField("name", opts.Name).
				WithError(err).
				Debugln("retrying instance insert")
		}),
	)
	if err != nil {
		logger.WithError(err).
			Errorln("instance insert failed")
		return nil, err
	}

	logger.Debugln("pending instance insert operation")
	// TODO: may be worth moving these pollings to a separate loop in Allocate
	// that way polling is decoupled from the initial creation loop and can be more
	// robust between insert calls
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
