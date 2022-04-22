// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package autoscaler

import (
	"context"
	"database/sql/driver"
	"errors"
)

// ProviderType specifies the hosting provider.
type ProviderType string

// Value converts the value to a sql string.
func (s ProviderType) Value() (driver.Value, error) {
	return string(s), nil
}

// Provider type enumeration.
const (
	ProviderAmazon       = ProviderType("amazon")
	ProviderAzure        = ProviderType("azure")
	ProviderBrightbox    = ProviderType("brightbox")
	ProviderDigitalOcean = ProviderType("digitalocean")
	ProviderGoogle       = ProviderType("google")
	ProviderHetznerCloud = ProviderType("hetznercloud")
	ProviderLinode       = ProviderType("linode")
	ProviderOpenStack    = ProviderType("openstack")
	ProviderPacket       = ProviderType("packet")
	ProviderScaleway     = ProviderType("scaleway")
	ProviderVultr        = ProviderType("vultr")
)

// ErrInstanceNotFound is returned when the requested
// instance does not exist in the cloud provider.
var ErrInstanceNotFound = errors.New("Not Found")

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
	Provider            ProviderType
	ID                  string
	Name                string
	Address             string
	Region              string
	Image               string
	Size                string
	ServiceAccountEmail string
	Scopes              []string
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
