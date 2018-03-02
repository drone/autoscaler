// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/scripts"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/rs/zerolog/log"
)

// Create creates the HetznerCloud instance.
func (p *Provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	logger := log.Ctx(ctx).With().
		Str("type", p.config.HetznerCloud.ServerType).
		Str("image", p.config.HetznerCloud.Image).
		Str("datacenter", p.config.HetznerCloud.Datacenter).
		Str("name", opts.Name).
		Logger()

	cloudinitTemplate := DefaultCloudConfig

	if p.config.CloudInit.Path != "" {
		if _, err := os.Stat(p.config.CloudInit.Path); os.IsNotExist(err) {
			logger.Error().
				Err(err).
				Msg("cloud-init doesn't exist")

			return nil, err
		}

		cloudinitRead, err := ioutil.ReadFile(p.config.CloudInit.Path)

		if err != nil {
			logger.Error().
				Err(err).
				Msg("cloud-init reading failed")

			return nil, err
		}

		cloudinitTemplate = string(cloudinitRead)
	}

	cloudinitGenerated, err := scripts.GenerateCloudInit(cloudinitTemplate, p.setupCloudInit(opts.Name))

	if err != nil {
		logger.Error().
			Err(err).
			Msg("cloud-init generate failed")

		return nil, err
	}

	req := hcloud.ServerCreateOpts{
		ServerType: &hcloud.ServerType{
			Name: p.config.HetznerCloud.ServerType,
		},
		Image: &hcloud.Image{
			Name: p.config.HetznerCloud.Image,
		},
		Datacenter: &hcloud.Datacenter{
			Name: p.config.HetznerCloud.Datacenter,
		},
		Name:     opts.Name,
		UserData: cloudinitGenerated,
	}

	logger.Debug().
		Msg("server create")

	srv, _, err := p.client().Server.Create(ctx, req)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("server create failed")

		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderHetznerCloud,
		ID:       strconv.Itoa(srv.Server.ID),
		Name:     srv.Server.Name,
		Address:  srv.Server.PublicNet.IPv4.IP.String(),
		Size:     req.ServerType.Name,
		Region:   req.Datacenter.Name,
		Image:    req.Image.Name,
		Secret:   p.secret,
	}

	logger.Info().
		Str("name", instance.Name).
		Msg("instance create success")

	// TODO: talk to the remote docker connection?

	// 	interval := time.Duration(0)

	// pinger:
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			logger.Debug().
	// 				Str("name", instance.Name).
	// 				Str("ip", instance.Address).
	// 				Str("port", "22").
	// 				Str("user", "root").
	// 				Msg("ping deadline exceeded")

	// 			return instance, ctx.Err()
	// 		case <-time.After(interval):
	// 			interval = time.Minute
	// 			logger.Debug().
	// 				Str("name", instance.Name).
	// 				Str("ip", instance.Address).
	// 				Str("port", "22").
	// 				Str("user", "root").
	// 				Msg("pinging server")

	// 			err = p.Provider.Ping(ctx, instance)

	// 			if err == nil {
	// 				break pinger
	// 			}
	// 		}
	// 	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("install complete")

	return instance, nil
}
