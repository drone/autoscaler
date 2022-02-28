// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"io/ioutil"

	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithClient returns an option to set the Hetzner client.
func WithClient(client *hcloud.Client) Option {
	return func(p *provider) {
		p.client = client
	}
}

// WithDatacenter returns an option to set the datacenter.
func WithDatacenter(datacenter string) Option {
	return func(p *provider) {
		p.datacenter = datacenter
	}
}

// WithImage returns an option to set the image.
func WithImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithServerType returns an option to set the server type.
func WithServerType(serverType string) Option {
	return func(p *provider) {
		p.serverType = serverType
	}
}

// WithSSHKey returns an option to set the ssh key.
func WithSSHKey(key int) Option {
	return func(p *provider) {
		p.key = key
	}
}

// WithToken returns an option to set the auth token.
func WithToken(token string) Option {
	return WithClient(
		hcloud.NewClient(
			hcloud.WithToken(
				token,
			),
		),
	)
}

// WithUserData returns an option to set the cloud-init
// template from text.
func WithUserData(text string) Option {
	return func(p *provider) {
		if text != "" {
			p.userdata = userdata.Parse(text)
		}
	}
}

// WithUserDataFile returns an option to set the cloud-init
// template from file.
func WithUserDataFile(filepath string) Option {
	return func(p *provider) {
		if filepath != "" {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				panic(err)
			}
			p.userdata = userdata.Parse(string(b))
		}
	}
}

// WithNetwork creates the VM attached to a hetzner network
func WithNetwork(name string) Option {
	return func(p *provider) {
		p.network = name
	}
}

// WithPrivateIP uses a private IP, input is the network the IP belong to. Implies WithNetwork (you can also use both)
func WithPrivateIP(name string) Option {
	return func(p *provider) {
		p.private = name
	}
}
