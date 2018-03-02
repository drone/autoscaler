// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/internal/base"
	"github.com/drone/autoscaler/drivers/internal/scripts"
	"github.com/drone/autoscaler/drivers/internal/sshutil"

	"github.com/dchest/uniuri"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/crypto/ssh"
)

// Provider defines the HetznerCloud provider.
type Provider struct {
	autoscaler.Provider

	signer ssh.Signer
	config config.Config
	secret string
}

// FromConfig loads the provider from the configuration.
func FromConfig(config config.Config) (autoscaler.Provider, error) {
	signer, err := sshutil.ParsePrivateKey(config.HetznerCloud.SSHKey)

	if err != nil {
		return nil, err
	}

	return &Provider{
		Provider: base.Provider("root", "22", signer),
		signer:   signer,
		config:   config,
		secret:   uniuri.New(),
	}, nil
}

func (p *Provider) setupCloudInit(name string) scripts.TemplateOpts {
	opts := scripts.TemplateOpts{}

	opts.Instance.Name = name
	opts.Server.Host = p.config.Agent.Host
	opts.Server.Secret = p.config.Agent.Token
	opts.Agent.Image = p.config.Agent.Image
	opts.Agent.Capacity = p.config.Agent.Concurrency
	opts.Cadvisor.Disable = p.config.Cadvisor.Disable
	opts.Cadvisor.Secret = p.secret

	return opts
}

func (p *Provider) client() *hcloud.Client {
	return hcloud.NewClient(
		hcloud.WithToken(
			p.config.HetznerCloud.Token,
		),
	)
}

const DefaultCloudConfig = `
#cloud-config

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

runcmd:
  - systemctl daemon-reload
  - systemctl enable docker
  - systemctl start --no-block docker
  - systemctl enable cadvisor
  - systemctl start --no-block cadvisor
  - systemctl enable agent
  - systemctl start --no-block agent

write_files:
{{ if not .Cadvisor.Disable }}
  - path: /etc/htaccess
    content: |
      admin:{SHA}{{ sha .Cadvisor.Secret }}
  - path: /etc/systemd/system/cadvisor.service
    content: |
      [Unit]
      Description=Cadvisor
      Requires=docker.service
      After=docker.service
      [Service]
      Restart=always
      RestartSec=10s
      ExecStop=/bin/sh -c 'docker ps | grep %p 1> /dev/null && docker stop -t 3600 %p || true'
      ExecStart=/usr/bin/docker run --rm \
        --name %p \
        --publish 8080:8080 \
        --volume /:/rootfs:ro \
        --volume /var/run:/var/run:rw \
        --volume /sys:/sys:ro \
        --volume /var/lib/docker/:/var/lib/docker:ro \
        --volume /dev/disk/:/dev/disk:ro \
        --volume /etc/htaccess:/etc/htaccess \
        google/cadvisor:latest \
        --http_auth_realm localhost \
        --http_auth_file /etc/htaccess
      [Install]
      WantedBy=multi-user.target
{{ end }}
  - path: /etc/systemd/system/agent.service
    content: |
      [Unit]
      Description=Agent
      Requires=docker.service
      After=docker.service
      [Service]
      Restart=always
      RestartSec=10s
      ExecStop=/bin/sh -c 'docker ps | grep %p 1> /dev/null && docker stop -t 3600 %p || true'
      ExecStart=/usr/bin/docker run --rm \
        --name %p \
        --env DRONE_SECRET={{ .Server.Secret }} \
        --env DRONE_SERVER={{ .Server.Host }} \
        --env DRONE_MAX_PROCS={{ .Agent.Capacity }} \
        --env DRONE_HOSTNAME={{ .Instance.Name }} \
        --volume /var/run/docker.sock:/var/run/docker.sock \
        {{ .Agent.Image }}
      [Install]
      WantedBy=multi-user.target
  - path: /etc/docker/daemon.json
    content: |
      {
        "storage-driver": "aufs"
      }
`
