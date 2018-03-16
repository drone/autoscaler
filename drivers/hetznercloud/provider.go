// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// provider implement a Hetzner Cloud provider.
type provider struct {
	init sync.Once

	token      string
	datacenter string
	serverType string
	image      string
	userdata   *template.Template
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
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	return p
}
