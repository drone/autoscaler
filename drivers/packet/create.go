// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"context"
	"time"

	"github.com/drone/autoscaler"
	"github.com/packethost/packngo"
	"github.com/rs/zerolog/log"
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

	logger := log.Ctx(ctx).With().
		Str("project", p.project).
		Str("facility", p.facility).
		Str("billing", p.billing).
		Str("plan", p.plan).
		Str("os", p.os).
		Str("hostname", p.hostname).
		Logger()

	cr := &packngo.DeviceCreateRequest{
		HostName:     p.hostname,
		Facility:     p.facility,
		Plan:         p.plan,
		OS:           p.os,
		ProjectID:    p.project,
		BillingCycle: p.billing,
		UserData:     buf.String(),
	}

	logger.Debug().
		Msg("provisioning new agent machine")

	d, _, err := p.client.Devices.Create(cr)
	if err != nil {
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderPacket,
		ID:       d.ID,
		Name:     opts.Name,
		Image:    d.OS.Slug,
		Region:   d.Facility.Code,
		Size:     d.Plan.Slug,
	}

	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.Debug().
				Str("name", instance.Name).
				Msg("device still inactive after given tiemout")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute
			d, _, err := p.client.Devices.Get(d.ID)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("getting device info")
				return nil, err
			}
			if d.State == "active" {
				for _, ip := range d.Network {
					if ip.Public && ip.AddressFamily == 4 {
						instance.Address = ip.Address
						break poller
					}
				}
			}
		}
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("device inserted")

	return instance, nil
}
