// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"strings"

	"github.com/drone/autoscaler/drivers/internal/scripts"
)

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithImage returns an option to set the image.
func WithImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithUserData returns an option to customize the cloud-init.
func WithUserData(userdata string) Option {
	return func(p *provider) {
		if userdata != "" {
			p.userdata = scripts.UserdataPrepare(userdata)
		}
	}
}

// WithRegion returns an option to set the target region.
func WithRegion(region string) Option {
	return func(p *provider) {
		p.region = region
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

// WithTags returns an option to set the image.
func WithTags(tags ...string) Option {
	return func(p *provider) {
		p.tags = tags
	}
}

// WithToken returns an option to set the auth token.
func WithToken(token string) Option {
	return func(p *provider) {
		p.token = token
	}
}

func withTags(tags string) Option {
	return func(p *provider) {
		p.tags = strings.Split(tags, ",")
	}
}
