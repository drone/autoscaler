// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/scripts"

	"github.com/alecthomas/template"
	"github.com/digitalocean/godo"
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

	req := &godo.DropletCreateRequest{
		Name:     opts.Name,
		Region:   p.region,
		Size:     p.size,
		Tags:     p.tags,
		IPv6:     false,
		UserData: buf.String(),
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: p.key},
		},
		Image: godo.DropletCreateImage{
			Slug: p.image,
		},
	}

	logger := log.Ctx(ctx).With().
		Str("region", req.Region).
		Str("image", req.Image.Slug).
		Str("size", req.Size).
		Str("name", req.Name).
		Logger()

	logger.Debug().
		Msg("instance create")

	client := newClient(ctx, p.token)
	droplet, _, err := client.Droplets.Create(ctx, req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot create instance")
		return nil, err
	}

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderDigitalOcean,
		ID:       strconv.Itoa(droplet.ID),
		Name:     droplet.Name,
		Size:     req.Size,
		Region:   req.Region,
		Image:    req.Image.Slug,
	}

	logger.Info().
		Str("name", instance.Name).
		Msg("instance created")

	// poll the digitalocean endpoint for server updates
	// and exit when a network address is allocated.
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.Debug().
				Str("name", instance.Name).
				Msg("cannot ascertain network")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute

			logger.Debug().
				Str("name", instance.Name).
				Msg("find instance network")

			droplet, _, err = client.Droplets.Get(ctx, droplet.ID)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("cannot find instance")
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
		Msg("instance network ready")

	return instance, nil
}

var cloudInitT = template.Must(template.New("_").Funcs(scripts.UserdataFuncmap).Parse(`#cloud-config
write_files:
  - path: /etc/systemd/system/docker.service.d/override.conf
    content: |
      [Service]
      ExecStart=
      ExecStart=/usr/bin/dockerd
  - path: /etc/default/docker
    content: |
      DOCKER_OPTS=""
  - path: /etc/docker/daemon.json
    content: |
      {
        "dns": [ "8.8.8.8", "8.8.4.4" ],
        "hosts": [ "0.0.0.0:2376", "unix:///var/run/docker.sock" ],
        "tls": true,
        "tlsverify": true,
        "tlscacert": "/etc/docker/ca.pem",
        "tlscert": "/etc/docker/server-cert.pem",
        "tlskey": "/etc/docker/server-key.pem"
      }
  - path: /etc/docker/ca.pem
    encoding: b64
    content: {{ .CACert | base64 }}
  - path: /etc/docker/server-cert.pem
    encoding: b64
    content: {{ .TLSCert | base64 }}
  - path: /etc/docker/server-key.pem
    encoding: b64
    content: {{ .TLSKey | base64 }}

runcmd:
  - [ systemctl, daemon-reload ]
  - [ systemctl, restart, docker ]
`))
