// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/scripts"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/rs/zerolog/log"
)

// Create creates the HetznerCloud instance.
func (p *Provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	req := hcloud.ServerCreateOpts{
		Name: opts.Name,
		ServerType: &hcloud.ServerType{
			Name: p.config.HetznerCloud.ServerType,
		},
		Image: &hcloud.Image{
			Name: p.config.HetznerCloud.Image,
		},
		Datacenter: &hcloud.Datacenter{
			Name: p.config.HetznerCloud.Datacenter,
		},
		SSHKeys: []*hcloud.SSHKey{
			&hcloud.SSHKey{
				ID: p.config.HetznerCloud.SSHKeyID,
			},
		},
	}
	if req.ServerType.Name == "" {
		req.ServerType.Name = "cx11"
	}
	if req.Image.Name == "" {
		req.Image.Name = "ubuntu-16.04"
	}
	if req.Datacenter.Name == "" {
		req.Datacenter.Name = "nbg1-dc3"
	}

	logger := log.Ctx(ctx).With().
		Str("datacenter", req.Datacenter.Name).
		Str("image", req.Image.Name).
		Str("serverType", req.ServerType.Name).
		Str("name", req.Name).
		Logger()

	logger.Debug().
		Msg("instance create")

	client := newClient(ctx, p.config.HetznerCloud.Token)
	resp, _, err := client.Server.Create(ctx, req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("instance create failed")
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderHetznerCloud,
		ID:       strconv.Itoa(resp.Server.ID),
		Name:     resp.Server.Name,
		Address:  resp.Server.PublicNet.IPv4.IP.String(),
		Size:     req.ServerType.Name,
		Region:   req.Datacenter.Name,
		Image:    req.Image.Name,
		Secret:   uniuri.New(),
	}

	logger.Info().
		Str("name", instance.Name).
		Msg("instance create success")

	// ping the server in a loop until we can successfully
	// authenticate.
	interval := time.Duration(0)
pinger:
	for {
		select {
		case <-ctx.Done():
			logger.Debug().
				Str("name", instance.Name).
				Str("ip", instance.Address).
				Str("port", "22").
				Str("user", "root").
				Msg("ping deadline exceeded")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute
			logger.Debug().
				Str("name", instance.Name).
				Str("ip", instance.Address).
				Str("port", "22").
				Str("user", "root").
				Msg("ping server")

			err = p.Provider.Ping(ctx, instance)
			if err == nil {
				break pinger
			}
		}
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("install agent")

	script, err := scripts.GenerateSetup(p.setupScriptOpts(instance))
	if err != nil {
		return instance, err
	}

	logs, err := p.Provider.Execute(ctx, instance, script)
	if err != nil {
		logger.Error().
			Err(err).
			Str("name", instance.Name).
			Str("ip", instance.Address).
			Msg("install failed")
		return instance, &autoscaler.InstanceError{Err: err, Logs: logs}
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("install complete")

	return instance, nil
}
