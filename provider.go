// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import "context"

// ProviderType specifies the hosting provider.
type ProviderType int

// Provider type enumeration.
const (
	ProviderUnknown ProviderType = iota
	ProviderAmazon
	ProviderAzure
	ProviderDigitalOcean
	ProviateGoogle
)

// A Provider represents a hosting provider, such as Digital Ocean
// and is responsible for server management.
type Provider interface {
	// Create creates a new server.
	Create(context.Context, *ServerOpts) (*Server, error)

	// Destroy destroys an existing server.
	Destroy(context.Context, *Server) error
}
