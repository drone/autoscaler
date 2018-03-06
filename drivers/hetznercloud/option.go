// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import "github.com/hetznercloud/hcloud-go/hcloud"

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
