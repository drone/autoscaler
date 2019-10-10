// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scaleway

import (
	"io/ioutil"

	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Option configures a Scaleway provider option.
type Option func(*provider) error

// WithAccessKey returns an option to set the user access key
func WithAccessKey(accessKey string) Option {
	return func(p *provider) error {
		p.accessKey = accessKey
		return nil
	}
}

// WithSecretKey returns an option to set the user secret key
func WithSecretKey(secretKey string) Option {
	return func(p *provider) error {
		p.secretKey = secretKey
		return nil
	}
}

// WithOrganisationID returns an option to set the user organisation id
func WithOrganisationID(orgId string) Option {
	return func(p *provider) error {
		p.orgID = orgId
		return nil
	}
}

// WithImage returns an option to set the image.
func WithImage(image string) Option {
	return func(p *provider) error {
		p.image = image
		return nil
	}
}

// WithDynamicIP returns an option to enable a dynamic IP.
func WithDynamicIP(dynamicIP bool) Option {
	return func(p *provider) error {
		p.dynamicIP = dynamicIP
		return nil
	}
}

// WithTags returns an option to set the server tags.
func WithTags(tags ...string) Option {
	return func(p *provider) error {
		p.tags = tags
		return nil
	}
}

// WithZone returns an option to set the target zone.
func WithZone(name string) Option {
	return func(p *provider) error {
		if name == "" {
			return nil
		}
		zone, err := scw.ParseZone(name)
		if err != nil {
			return err
		}
		p.zone = zone
		return nil
	}
}

// WithSize returns an option to set the instance size.
func WithSize(size string) Option {
	return func(p *provider) error {
		p.size = size
		return nil
	}
}

// WithUserData returns an option to set the cloud-init
// template from text.
func WithUserData(text string) Option {
	return func(p *provider) error {
		if text != "" {
			p.userdata = userdata.Parse(text)
		}
		return nil
	}
}

// WithUserDataFile returns an option to set the cloud-init
// template from file.
func WithUserDataFile(filepath string) Option {
	return func(p *provider) error {
		if filepath != "" {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				return err
			}
			p.userdata = userdata.Parse(string(b))
		}
		return nil
	}
}
