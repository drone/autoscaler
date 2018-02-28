// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import (
	"context"
	"errors"
)

// ServerState specifies the server state.
type ServerState string

// ServerState type enumeration.
const (
	StatePending  = ServerState("pending")
	StateStaging  = ServerState("staging")
	StateRunning  = ServerState("running")
	StateShutdown = ServerState("shutdown")
	StateStopping = ServerState("stopping")
	StateStopped  = ServerState("stopped")
	StateError    = ServerState("error")
)

// ErrServerNotFound is returned when the requested server
// does not exist in the store.
var ErrServerNotFound = errors.New("Not Found")

// A ServerStore persists server information.
type ServerStore interface {
	// Find a server by unique name.
	Find(context.Context, string) (*Server, error)

	// List returns all registered servers
	List(context.Context) ([]*Server, error)

	// ListState returns all servers with the given state.
	ListState(context.Context, ServerState) ([]*Server, error)

	// Create the server record in the store.
	Create(context.Context, *Server) error

	// Update the server record in the store.
	Update(context.Context, *Server) error

	// Delete the server record from the store.
	Delete(context.Context, *Server) error
}

// Server stores the server details.
type Server struct {
	Provider ProviderType `json:"provider"`
	State    ServerState  `json:"state"`
	UID      string       `json:"uid"`
	Name     string       `json:"name"`
	Image    string       `json:"image"`
	Region   string       `json:"region"`
	Size     string       `json:"size"`
	Address  string       `json:"address"`
	Secret   string       `json:"secret"`
	Capacity int          `json:"capacity"`
	Active   bool         `json:"active"`
	Healthy  bool         `json:"healthy"`
	Created  int64        `json:"created"`
	Updated  int64        `json:"updated"`
	Logs     string       `json:"-"`
}

// ByCreated sorts the server list by created date.
type ByCreated []*Server

func (a ByCreated) Len() int           { return len(a) }
func (a ByCreated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreated) Less(i, j int) bool { return a[i].Created < a[j].Created }
