// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import (
	"io/ioutil"

	"github.com/drone/autoscaler/drivers/internal/userdata"
)

// Option configures a Brightbox provider option.
type Option func(*provider) error

// WithApiURL returns an option to set the API endpoint URL
func WithApiURL(apiURL string) Option {
	return func(p *provider) error {
		p.apiURL = apiURL
		return nil
	}
}

// WithClientID returns an option to set the API client ID
func WithClientID(clientID string) Option {
	return func(p *provider) error {
		p.clientID = clientID
		return nil
	}
}

// WithClientSecret returns an option to set the API client secret
func WithClientSecret(clientSecret string) Option {
	return func(p *provider) error {
		p.clientSecret = clientSecret
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

// WithServerType returns an option to set the server type
func WithServerType(serverType string) Option {
	return func(p *provider) error {
		p.serverType = serverType
		return nil
	}
}

// WithServerGroups returns an option to set the server groups
func WithServerGroups(serverGroups []string) Option {
	return func(p *provider) error {
		p.serverGroups = serverGroups
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
