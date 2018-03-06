// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"sync"

	"github.com/drone/autoscaler"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// provider implement a Hetzner Cloud provider.
type provider struct {
	init sync.Once

	token      string
	datacenter string
	serverType string
	image      string
	key        int

	client *hcloud.Client
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.datacenter == "" {
		p.datacenter = "nbg1-dc3"
	}
	if p.serverType == "" {
		p.serverType = "cx11"
	}
	if p.image == "" {
		p.image = "ubuntu-16.04"
	}
	return p
}
