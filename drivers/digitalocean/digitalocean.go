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
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/internal/scripts"
	"github.com/drone/autoscaler/drivers/internal/sshutil"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// Provider defines the DigitalOcean provider.
type Provider struct {
	config config.Config
}

// FromConfig loads the provider from the configuration.
func FromConfig(config config.Config) (autoscaler.Provider, error) {
	return &Provider{
		config: config,
	}, nil
}

// Create creates the DigitalOcean instance.
func (p *Provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	signer, err := sshutil.ParsePrivateKey(p.config.DigitalOcean.SSHKey)
	if err != nil {
		return nil, err
	}

	req := &godo.DropletCreateRequest{
		Name:   opts.Name,
		Region: p.config.DigitalOcean.Region,
		Size:   p.config.DigitalOcean.Size,
		IPv6:   p.config.DigitalOcean.IPv6,
		Tags:   p.config.DigitalOcean.Tags,
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: sshutil.Fingerprint(signer)},
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

	for {
		logger.Debug().
			Str("name", instance.Name).
			Msg("droplet network check")

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
			break
		}

		logger.Debug().
			Str("name", instance.Name).
			Msg("droplet network not available")

		time.Sleep(5 * time.Second)
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("droplet network ready")

	script, err := scripts.GenerateSetup(p.setupScriptOpts(instance))
	if err != nil {
		return instance, err
	}

	// ping the server in a loop until we can successfully
	// authenticate.
	for i := 0; i < 20; i++ {
		logger.Debug().
			Str("name", instance.Name).
			Str("ip", instance.Address).
			Str("port", "22").
			Str("user", "root").
			Msg("ping server")
		_, err = sshutil.Execute(instance.Address, "22", "root", "whoami", signer)
		if err == nil {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("install agent")

	out, err := sshutil.Execute(instance.Address, "22", "root", script, signer)
	if len(out) != 0 {
		// TODO bundle server logs in error message
	}
	if err != nil {
		logger.Error().
			Err(err).
			Str("name", instance.Name).
			Str("ip", instance.Address).
			Msg("install failed")
		return instance, err
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("install complete")

	return instance, nil
}

// Destroy destroyes the DigitalOcean instance.
func (p *Provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("region", instance.Region).
		Str("image", instance.Image).
		Str("size", instance.Size).
		Str("name", instance.Name).
		Logger()

	script, err := scripts.GenerateTeardown()
	if err != nil {
		return err
	}

	signer, err := sshutil.ParsePrivateKey(p.config.DigitalOcean.SSHKey)
	if err != nil {
		return err
	}

	logger.Debug().
		Msg("teardown droplet")

	_, err = sshutil.Execute(instance.Address, "22", "root", script, signer)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("teardown failed")
		return err
	}

	logger.Debug().
		Msg("deleting droplet")

	client := newClient(ctx, p.config.DigitalOcean.Token)
	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	_, err = client.Droplets.Delete(ctx, id)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("deleting droplet failed")
		return err
	}

	logger.Debug().
		Msg("droplet deleted")

	return nil
}

func (p *Provider) setupScriptOpts(instance *autoscaler.Instance) scripts.SetupOpts {
	opts := scripts.SetupOpts{}
	opts.Server.Host = p.config.Agent.Host
	opts.Server.Secret = p.config.Agent.Token
	opts.Agent.Image = p.config.Agent.Image
	opts.Agent.Capacity = p.config.Agent.Concurrency
	opts.Instance.Addr = instance.Address
	opts.Instance.Name = instance.Name
	opts.Cadvisor.Disable = false
	opts.Cadvisor.Secret = instance.Secret
	return opts
}

// helper function returns a new digitalocean client.
func newClient(ctx context.Context, token string) *godo.Client {
	return godo.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		)),
	)
}
