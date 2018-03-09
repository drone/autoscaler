// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import "context"

// ProviderType specifies the hosting provider.
type ProviderType string

// Provider type enumeration.
const (
	ProviderAmazon       = ProviderType("amazon")
	ProviderAzure        = ProviderType("azure")
	ProviderDigitalOcean = ProviderType("digitalocean")
	ProviderGoogle       = ProviderType("google")
	ProviderHetznerCloud = ProviderType("hetznercloud")
	ProviderLinode       = ProviderType("linode")
	ProviderOpenStack    = ProviderType("openstack")
	ProviderPacket       = ProviderType("packet")
	ProviderScaleway     = ProviderType("scaleway")
	ProviderVultr        = ProviderType("vultr")
)

// A Provider represents a hosting provider, such as
// Digital Ocean and is responsible for server management.
type Provider interface {
	// Create creates a new server.
	Create(context.Context, InstanceCreateOpts) (*Instance, error)
	// Destroy destroys an existing server.
	Destroy(context.Context, *Instance) error
}

// An Instance represents a server instance
// (e.g Digital Ocean Droplet).
type Instance struct {
	Provider ProviderType
	ID       string
	Name     string
	Address  string
	Region   string
	Image    string
	Size     string
}

// InstanceCreateOpts define soptional instructions for
// creating server instances.
type InstanceCreateOpts struct {
	Name    string
	CAKey   []byte
	CACert  []byte
	TLSKey  []byte
	TLSCert []byte
}

// InstanceError snapshots an error creating an instance
// with server logs.
type InstanceError struct {
	Err  error
	Logs []byte
}

// Error implements the error interface.
func (e *InstanceError) Error() string {
	return e.Err.Error()
}
