// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"sync"

	"github.com/digitalocean/godo"
	"github.com/drone/autoscaler"
	"golang.org/x/oauth2"
)

// provider implements a DigitalOcean provider.
type provider struct {
	init sync.Once

	key    string
	region string
	token  string
	size   string
	image  string
	tags   []string
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.region == "" {
		p.region = "sfo1"
	}
	if p.size == "" {
		p.size = "s-1vcpu-1gb"
	}
	if p.image == "" {
		p.image = "docker-16-04"
	}
	return p
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
