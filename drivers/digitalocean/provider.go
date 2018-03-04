// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/internal/sshutil"

	"github.com/digitalocean/godo"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

// Provider defines the DigitalOcean provider.
type Provider struct {
	signer ssh.Signer
	config config.Config
}

// FromConfig loads the provider from the configuration.
func FromConfig(config config.Config) (autoscaler.Provider, error) {
	signer, err := sshutil.ParsePrivateKey(config.DigitalOcean.SSHKey)
	if err != nil {
		return nil, err
	}
	return &Provider{
		signer: signer,
		config: config,
	}, nil
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
