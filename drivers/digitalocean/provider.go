// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"sync"

	"github.com/drone/autoscaler"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// provider defines the DigitalOcean provider.
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
	p := &provider{
		region: "sfo1",
		size:   "s-1vcpu-1gb",
		image:  "docker-16-04",
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func newClient(ctx context.Context, token string) *godo.Client {
	return godo.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		)),
	)
}
