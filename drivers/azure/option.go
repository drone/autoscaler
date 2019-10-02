// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
)


type Option func(*provider)

func WithSubscriptionId(n string) Option {
	return func(p *provider) {
		p.subscriptionId = n
	}
}

func WithResourceGroup(n string) Option {
	return func(p *provider) {
		p.resourceGroup = n
	}
}

func WithImage(n string) Option {
	return func(p *provider) {
		p.image = n
	}
}

func WithLocation(n string) Option {
	return func(p *provider) {
		p.location = n
	}
}

func WithAdminUsername(n string) Option {
	return func(p *provider) {
		p.adminUsername = n
	}
}

func WithAdminPassword(n string) Option {
	return func(p *provider) {
		p.adminPassword = n
	}
}

func WithVMName(n string) Option {
	return func(p *provider) {
		p.vmName = n
	}
}

func WithSSHKey(n string) Option {
	return func(p *provider) {
		p.sshKey = n
	}
}

func WithVMSize(n string) Option {
	return func(p *provider) {
		p.vmSize = compute.VirtualMachineSizeTypes(n)
	}
}

func WithVNet(n string) Option {
	return func(p *provider) {
		p.vnet = n
	}
}

func WithNSG(n string) Option {
	return func(p *provider) {
		p.nsg = n
	}
}

func WithSubnet(n string) Option {
	return func(p *provider) {
		p.subnet = n
	}
}

func WithVolumeSize(n int32) Option {
	return func(p *provider) {
		p.volumeSize = n
	}
}

func WithImageOffer(n string) Option {
	return func(p *provider) {
		p.imageOffer = n
	}
}

func WithImagePublisher(n string) Option {
	return func(p *provider) {
		p.imagePublisher = n
	}
}

func WithImageSKU(n string) Option {
	return func(p *provider) {
		p.imageSKU = n
	}
}


func WithImageVersion(n string) Option {
	return func(p *provider) {
		p.imageVersion = n
	}
}
