// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/internal/base"
	"github.com/drone/autoscaler/drivers/internal/scripts"
	"github.com/drone/autoscaler/drivers/internal/sshutil"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/crypto/ssh"
)

// Provider defines the HetznerCloud provider.
type Provider struct {
	autoscaler.Provider

	signer ssh.Signer
	config config.Config
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
	}, nil
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

// helper function returns a new HetznerCloud client.
func newClient(ctx context.Context, token string) *hcloud.Client {
	return hcloud.NewClient(
		hcloud.WithToken(
			token,
		),
	)
}
