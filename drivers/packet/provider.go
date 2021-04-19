// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package packet

import (
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/packethost/packngo"
)

const consumerToken = "24e70949af5ecd17fe8e867b335fc88e7de8bd4ad617c0403d8769a376ddea72"

// provider implements a Packet.net provider.
type provider struct {
	init sync.Once

	apikey   string
	billing  string
	facility string
	os       string
	plan     string
	project  string
	sshkey   string
	hostname string
	tags     []string
	userdata *template.Template

	client *packngo.Client
}

// New returns a new Packet.net provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.facility == "" {
		p.facility = "ewr1"
	}
	if p.os == "" {
		p.os = "ubuntu_18_04"
	}
	if p.plan == "" {
		p.plan = "baremetal_0"
	}
	if p.billing == "" {
		p.billing = "hourly"
	}
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	if p.client == nil {
		p.client, _ = packngo.NewClientWithBaseURL(
			consumerToken, p.apikey, nil, "https://api.equinix.com/metal/v1/")
	}
	return p
}
