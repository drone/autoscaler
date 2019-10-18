// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// provider implements an OpenStack provider
type provider struct {
	init sync.Once

	key      string
	region   string
	image    string
	flavor   string
	subnet   string
	pool     string
	userdata *template.Template
	groups   []string
	metadata map[string]string

	computeClient *gophercloud.ServiceClient
}

// New returns a new OpenStack provider.
func New(opts ...Option) (autoscaler.Provider, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}

	if p.userdata == nil {
		p.userdata = userdata.T
	}

	if p.computeClient == nil {
		authOpts, err := openstack.AuthOptionsFromEnv()
		if err != nil {
			return nil, err
		}

		authClient, err := openstack.AuthenticatedClient(authOpts)
		if err != nil {
			return nil, err
		}

		p.computeClient, err = openstack.NewComputeV2(authClient, gophercloud.EndpointOpts{
			Region: p.region,
		})
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
