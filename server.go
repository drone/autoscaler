// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package autoscaler

import (
	"context"
	"database/sql/driver"
	"errors"
)

// ServerState specifies the server state.
type ServerState string

// Value converts the value to a sql string.
func (s ServerState) Value() (driver.Value, error) {
	return string(s), nil
}

// ServerState type enumeration.
const (
	StatePending  = ServerState("pending")
	StateCreating = ServerState("creating")
	StateCreated  = ServerState("created")
	StateStaging  = ServerState("staging") // starting
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

	// Update the server record that it is busy.
	Busy(context.Context, *Server) error

	// Delete the server record from the store.
	Delete(context.Context, *Server) error

	// Purge old server records from the store.
	Purge(context.Context, int64) error
}

// Server stores the server details.
type Server struct {
	ID       string       `db:"server_id"       json:"id"`
	Provider ProviderType `db:"server_provider" json:"provider"`
	State    ServerState  `db:"server_state"    json:"state"`
	Name     string       `db:"server_name"     json:"name"`
	Image    string       `db:"server_image"    json:"image"`
	Region   string       `db:"server_region"   json:"region"`
	Size     string       `db:"server_size"     json:"size"`
	Platform string       `db:"server_platform" json:"platform"`
	Address  string       `db:"server_address"  json:"address"`
	Capacity int          `db:"server_capacity" json:"capacity"`
	Secret   string       `db:"server_secret"   json:"secret"`
	Error    string       `db:"server_error"    json:"error"`
	CAKey    []byte       `db:"server_ca_key"   json:"ca_key"`
	CACert   []byte       `db:"server_ca_cert"  json:"ca_cert"`
	TLSKey   []byte       `db:"server_tls_key"  json:"tls_key"`
	TLSCert  []byte       `db:"server_tls_cert" json:"tls_cert"`
	Created  int64        `db:"server_created"  json:"created"`
	Updated  int64        `db:"server_updated"  json:"updated"`
	Started  int64        `db:"server_started"  json:"started"`
	Stopped  int64        `db:"server_stopped"  json:"stopped"`
	LastBusy int64        `db:"server_lastbusy" json:"lastbusy"`
}
