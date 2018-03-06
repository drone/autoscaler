// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"strconv"
	"text/template"

	"github.com/drone/autoscaler"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/rs/zerolog/log"
)

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	p.init.Do(func() {
		p.setup(ctx)
	})

	buf := new(bytes.Buffer)
	err := cloudInitT.Execute(buf, &opts)
	if err != nil {
		return nil, err
	}

	req := hcloud.ServerCreateOpts{
		Name:     opts.Name,
		UserData: buf.String(),
		ServerType: &hcloud.ServerType{
			Name: p.serverType,
		},
		Image: &hcloud.Image{
			Name: p.image,
		},
		Datacenter: &hcloud.Datacenter{
			Name: p.datacenter,
		},
		SSHKeys: []*hcloud.SSHKey{
			&hcloud.SSHKey{
				ID: p.key,
			},
		},
	}

	logger := log.Ctx(ctx).With().
		Str("datacenter", req.Datacenter.Name).
		Str("image", req.Image.Name).
		Str("serverType", req.ServerType.Name).
		Str("name", req.Name).
		Logger()

	logger.Debug().
		Msg("instance create")

	resp, _, err := p.client.Server.Create(ctx, req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot create instance")
		return nil, err
	}

	logger.Info().
		Str("name", req.Name).
		Msg("instance created")

	return &autoscaler.Instance{
		Provider: autoscaler.ProviderHetznerCloud,
		ID:       strconv.Itoa(resp.Server.ID),
		Name:     resp.Server.Name,
		Address:  resp.Server.PublicNet.IPv4.IP.String(),
		Size:     req.ServerType.Name,
		Region:   req.Datacenter.Name,
		Image:    req.Image.Name,
	}, nil
}

var cloudInitT = template.Must(template.New("_").Funcs(funcmap).Parse(`#cloud-config

apt_reboot_if_required: false
package_update: false
package_upgrade: false

apt:
  sources:
    docker.list:
      source: deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable
      keyid: 0EBFCD88

packages:
  - docker-ce

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

var funcmap = map[string]interface{}{
	"base64": func(src []byte) string {
		return base64.StdEncoding.EncodeToString(src)
	},
}
