// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import "strings"

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithImage returns an option to set the image.
func WithImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithRegion returns an option to set the target region.
func WithRegion(region string) Option {
	return func(p *provider) {
		p.region = region
	}
}

// WithSecurityGroup returns an option to set the instance size.
func WithSecurityGroup(id string) Option {
	return func(p *provider) {
		p.groups = []string{id}
	}
}

// WithSize returns an option to set the instance size.
func WithSize(size string) Option {
	return func(p *provider) {
		p.size = size
	}
}

// WithSSHKey returns an option to set the ssh key.
func WithSSHKey(key string) Option {
	return func(p *provider) {
		p.key = key
	}
}

// WithSubnet returns an option to set the subnet id.
func WithSubnet(id string) Option {
	return func(p *provider) {
		p.subnet = id
	}
}

// WithTags returns an option to set the image.
func WithTags(tags ...string) Option {
	return func(p *provider) {
		p.tags = tags
	}
}

func withTags(tags string) Option {
	return func(p *provider) {
		p.tags = strings.Split(tags, ",")
	}
}
