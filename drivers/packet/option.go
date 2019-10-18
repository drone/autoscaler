// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package packet

import (
	"io/ioutil"

	"github.com/drone/autoscaler/drivers/internal/userdata"
)

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithAPIKey returns an option to set the api key.
func WithAPIKey(apikey string) Option {
	return func(p *provider) {
		p.apikey = apikey
	}
}

// WithFacility returns an option to set the target facility.
func WithFacility(facility string) Option {
	return func(p *provider) {
		p.facility = facility
	}
}

// WithPlan returns an option to set the plan.
func WithPlan(plan string) Option {
	return func(p *provider) {
		p.plan = plan
	}
}

// WithOS returns an option to set the operating system.
func WithOS(os string) Option {
	return func(p *provider) {
		p.os = os
	}
}

// WithProject returns an option to set the project id.
func WithProject(project string) Option {
	return func(p *provider) {
		p.project = project
	}
}

// WithSSHKey returns an option to set the ssh key.
func WithSSHKey(sshkey string) Option {
	return func(p *provider) {
		p.sshkey = sshkey
	}
}

// WithHostname returns an option to set the hostname
func WithHostname(hostname string) Option {
	return func(p *provider) {
		if hostname != "" {
			p.hostname = hostname
		}
	}
}

// WithTags returns an option to set the image.
func WithTags(tags ...string) Option {
	return func(p *provider) {
		p.tags = tags
	}
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
