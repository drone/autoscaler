// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithToken returns an option to set the auth token.
func WithToken(token string) Option {
	return func(p *provider) {
		p.token = token
	}
}

// WithFingerprint returns an option to set the ssh key.
func WithFingerprint(fingerprint string) Option {
	return func(p *provider) {
		p.key = fingerprint
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

// WithImage returns an option to set the image.
func WithImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithTags returns an option to set the image.
func WithTags(tags ...string) Option {
	return func(p *provider) {
		p.tags = tags
	}
}
