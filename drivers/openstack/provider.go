// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"regexp"
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

// provider implements an OpenStack provider
type provider struct {
	init sync.Once

	key      string
	region   string
	image    string
	flavor   string
	network  string
	pool     string
	userdata *template.Template
	groups   []string
	metadata map[string]string

	computeClient *gophercloud.ServiceClient
	networkClient *gophercloud.ServiceClient
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

	if p.networkClient == nil {
		authOpts, err := openstack.AuthOptionsFromEnv()
		if err != nil {
			return nil, err
		}

		authClient, err := openstack.AuthenticatedClient(authOpts)
		if err != nil {
			return nil, err
		}

		p.networkClient, err = openstack.NewNetworkV2(authClient, gophercloud.EndpointOpts{
			Region: p.region,
		})
		if err != nil {
			return nil, err
		}
	}

	if p.image != "" && !isUUID(p.image) {
		uuid, err := images.IDFromName(p.computeClient, p.image)
		if err != nil {
			return nil, err
		}
		p.image = uuid
	}

	if p.flavor != "" && !isUUID(p.flavor) {
		uuid, err := flavors.IDFromName(p.computeClient, p.flavor)
		if err != nil {
			return nil, err
		}
		p.flavor = uuid
	}

	if p.network != "" && !isUUID(p.network) {
		uuid, err := networks.IDFromName(p.networkClient, p.network)
		if err != nil {
			return nil, err
		}
		p.network = uuid
	}

	return p, nil
}

func isUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
