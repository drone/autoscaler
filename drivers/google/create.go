// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	"github.com/google/uuid"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// insertUnconfirmedError signals that Google may have accepted an Insert (so
// the name may be reserved and the VM may exist) but its outcome couldn't be
// confirmed: the Insert got no definitive answer, the operation poll gave up
// or ctx ended, or the follow-up Get failed. It is not an Insert rejection
// and must never be retried as one - re-inserting the same name just 409s.
//
// It is always returned with a partial *autoscaler.Instance (ID and Region).
// That lets the engine record the VM on the errored server, so the reaper
// destroys it instead of only deleting the DB row (empty server.ID) and
// orphaning a running VM.
type insertUnconfirmedError struct {
	Err error
}

// Error implements the error interface.
func (e *insertUnconfirmedError) Error() string {
	return fmt.Sprintf("instance insert may have been accepted but its outcome could not be confirmed: %s", e.Err)
}

// Unwrap lets errors.Is/errors.As see through to the underlying error.
func (e *insertUnconfirmedError) Unwrap() error { return e.Err }

// Create provisions an instance. createTimeout bounds the whole call: the
// insert, its operation poll and the follow-up get.
func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, p.createTimeout)
	defer cancel()

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

	// This UUID is used to make sure retried inserts are treated as the same call
	// in case an insert returns a transient error but is still actioned on google's
	// api side.
	requestID := uuid.New().String()

	// An expired ctx means no request was sent, so there is nothing to
	// clean up.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var op *compute.Operation
	err = doWithRetry(ctx, "instances.insert", p.maxTransientDelay, func() error {
		var err error
		op, err = p.service.Instances.Insert(p.project, zone, in).RequestId(requestID).Context(ctx).Do()
		return err
	})
	if err != nil {
		if !insertOutcomeKnown(err) {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.WithError(err).
					Warnln("slow response from provider: the instance insert request did not complete in time, leaving the instance for the reaper")
			} else {
				logger.WithError(err).
					Errorln("instance insert failed")
			}
			return p.unconfirmedInsert(opts, name, zone, err)
		}
		logger.WithError(err).
			Errorln("instance insert failed")
		return nil, err
	}

	logger.Debugln("pending instance insert operation")
	// TODO: may be worth moving these pollings to a separate loop in Allocate
	// that way polling is decoupled from the initial creation loop and can be more
	// robust between insert calls / and be safe during autoscaler restarts.
	err = p.waitZoneOperation(ctx, op.Name, zone)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.WithError(err).
				Warnln("slow response from provider: the instance insert did not finish in time, leaving the instance for the reaper")
		} else {
			logger.WithError(err).
				Errorln("instance insert operation failed")
		}
		return p.unconfirmedInsert(opts, name, zone, err)
	}

	logger.Debugln("instance insert operation complete")

	var resp *compute.Instance
	err = doWithRetry(ctx, "instances.get", p.maxTransientDelay, func() error {
		var err error
		resp, err = p.service.Instances.Get(p.project, zone, name).Context(ctx).Do()
		return err
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.WithError(err).
				Warnln("slow response from provider: the instance details were not available in time, leaving the instance for the reaper")
		} else {
			logger.WithError(err).
				Errorln("cannot get instance details")
		}
		return p.unconfirmedInsert(opts, name, zone, err)
	}
	if resp == nil {
		return nil, errors.New("instances.Get returned no details")
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

// insertOutcomeKnown reports whether err from an Insert is a definitive answer
// from Google (a rejection such as a stockout, quota, rate limit or conflict),
// so nothing was created. Anything else - a context error, a network error, or
// a transient status that survived every retry - leaves the outcome unknown:
// the request may have been accepted.
func insertOutcomeKnown(err error) bool {
	var gerr *googleapi.Error
	return errors.As(err, &gerr) && !isTransientError(gerr)
}

// unconfirmedInsert handles a failure after an Insert that may have been
// accepted. An *operationError (only reported once DONE) means nothing was
// created, so err is returned as-is with no instance. Anything else is an
// unknown outcome: return the instance's identity with an
// *insertUnconfirmedError.
func (p *provider) unconfirmedInsert(opts autoscaler.InstanceCreateOpts, name, zone string, err error) (*autoscaler.Instance, error) {
	var opErr *operationError
	if errors.As(err, &opErr) {
		return nil, err
	}
	return &autoscaler.Instance{
		Provider:            autoscaler.ProviderGoogle,
		ID:                  name,
		Name:                opts.Name,
		Image:               p.image,
		Region:              zone,
		Size:                p.size,
		ServiceAccountEmail: p.serviceAccountEmail,
		Scopes:              p.scopes,
	}, &insertUnconfirmedError{Err: err}
}
