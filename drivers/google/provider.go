// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"sync"

	"github.com/drone/autoscaler"
	"google.golang.org/api/compute/v1"
)

// provider implements a DigitalOcean provider.
type provider struct {
	init sync.Once

	key   string
	proj  string
	zone  string
	token string
	size  string
	image string
	tags  []string

	service *compute.Service
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.zone == "" {
		p.zone = "us-central1-a"
	}
	if p.size == "" {
		p.size = "n1-standard-1"
	}
	if p.image == "" {
		p.image = "ubuntu-1510-wily-v20151114"
	}
	return p
}
