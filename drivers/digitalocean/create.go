// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/scripts"
	"github.com/drone/autoscaler/drivers/internal/sshutil"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

// Create creates the DigitalOcean instance.
func (p *Provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	req := &godo.DropletCreateRequest{
		Name:   opts.Name,
		Region: p.config.DigitalOcean.Region,
		Size:   p.config.DigitalOcean.Size,
		IPv6:   p.config.DigitalOcean.IPv6,
		Tags:   p.config.DigitalOcean.Tags,
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: sshutil.Fingerprint(p.signer)},
		},
		Image: godo.DropletCreateImage{
			Slug: p.config.DigitalOcean.Image,
		},
	}
	if req.Image.Slug == "" {
		req.Image.Slug = "docker-16-04"
	}
	if req.Size == "" {
		req.Size = "s-1vcpu-1gb"
	}
	if req.Region == "" {
		req.Region = "sfo1"
	}

	logger := log.Ctx(ctx).With().
		Str("region", req.Region).
		Str("image", req.Image.Slug).
		Str("size", req.Size).
		Str("name", req.Name).
		Logger()

	logger.Debug().
		Msg("droplet create")

	client := newClient(ctx, p.config.DigitalOcean.Token)
	droplet, _, err := client.Droplets.Create(ctx, req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("droplet create failed")
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderDigitalOcean,
		ID:       strconv.Itoa(droplet.ID),
		Name:     droplet.Name,
		Size:     req.Size,
		Region:   req.Region,
		Image:    req.Image.Slug,
		Secret:   uniuri.New(),
	}

	logger.Info().
		Str("name", instance.Name).
		Msg("droplet create success")

	// poll the digitalocean endpoint for server updates
	// and exit when a network address is allocated.
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.Debug().
				Str("name", instance.Name).
				Msg("droplet network deadline exceeded")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute

			logger.Debug().
				Str("name", instance.Name).
				Msg("check droplet network")

			droplet, _, err = client.Droplets.Get(ctx, droplet.ID)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("droplet details unavailable")
				return instance, err
			}

			for _, network := range droplet.Networks.V4 {
				if network.Type == "public" {
					instance.Address = network.IPAddress
				}
			}

			if instance.Address != "" {
				break poller
			}
		}
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("droplet network ready")

	// ping the server in a loop until we can successfully
	// authenticate.
	interval = time.Duration(0)
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
