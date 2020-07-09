// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"io/ioutil"
	"net/http"

	"github.com/drone/autoscaler/drivers/internal/userdata"

	"google.golang.org/api/compute/v1"
)

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithClient returns an option to set the default http
// Client used with the Google Compute provider.
func WithClient(client *http.Client) Option {
	return func(p *provider) {
		service, err := compute.New(client)
		if err != nil {
			panic(err)
		}
		p.service = service
	}
}

// WithDiskSize returns an option to set the instance disk
// size in gigabytes.
func WithDiskSize(diskSize int64) Option {
	return func(p *provider) {
		p.diskSize = diskSize
	}
}

// WithDiskType returns an option to set the instance disk type.
func WithDiskType(diskType string) Option {
	return func(p *provider) {
		p.diskType = diskType
	}
}

// WithLabels returns an option to set the metadata labels.
func WithLabels(labels map[string]string) Option {
	return func(p *provider) {
		p.labels = labels
	}
}

// WithMachineImage returns an option to set the image.
func WithMachineImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithMachineType returns an option to set the instance type.
func WithMachineType(size string) Option {
	return func(p *provider) {
		p.size = size
	}
}

// WithNetwork returns an option to set the network.
func WithNetwork(network string) Option {
	return func(p *provider) {
		p.network = network
	}
}

// WithSubNetwork returns an option to set the subnetwork.
func WithSubnetwork(subnetwork string) Option {
	return func(p *provider) {
		p.subnetwork = subnetwork
	}
}

// WithPrivateIP returns an option to set the private IP address.
func WithPrivateIP(private bool) Option {
	return func(p *provider) {
		p.privateIP = private
	}
}

// WithProject returns an option to set the project.
func WithProject(project string) Option {
	return func(p *provider) {
		p.project = project
	}
}

// WithTags returns an option to set the resource tags.
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

// WithZone returns an option to set the target zone.
func WithZone(zone string) Option {
	return func(p *provider) {
		p.zone = zone
	}
}

// WithScopes returns an option to set the scopes.
func WithScopes(scopes string) Option {
	return func(p *provider) {
		p.scopes = scopes
	}
}