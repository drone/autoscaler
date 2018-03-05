// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/drone/autoscaler"

	"github.com/jmoiron/sqlx"
)

// NewServerStore returns a new server store.
func NewServerStore(db *sqlx.DB) autoscaler.ServerStore {
	return &serverStore{db}
}

type serverStore struct {
	*sqlx.DB
}

func (db *serverStore) Find(ctx context.Context, name string) (*autoscaler.Server, error) {
	dest := &autoscaler.Server{Name: name}
	stmt, args, err := db.BindNamed(serverFindStmt, dest)
	if err != nil {
		return nil, err
	}
	err = db.GetContext(ctx, dest, stmt, args...)
	return dest, err
}

func (db *serverStore) List(ctx context.Context) ([]*autoscaler.Server, error) {
	dest := []*autoscaler.Server{}
	err := db.SelectContext(ctx, &dest, serverListStmt)
	return dest, err
}

func (db *serverStore) ListState(ctx context.Context, state autoscaler.ServerState) ([]*autoscaler.Server, error) {
	dest := []*autoscaler.Server{}
	stmt, args, err := db.BindNamed(serverListStateStmt, map[string]interface{}{"server_state": state})
	if err != nil {
		return nil, err
	}
	err = db.SelectContext(ctx, &dest, stmt, args...)
	if err == sql.ErrNoRows {
		return dest, nil
	}
	return dest, err
}

func (db *serverStore) Create(ctx context.Context, server *autoscaler.Server) error {
	server.Created = time.Now().Unix()
	server.Updated = time.Now().Unix()
	stmt, args, err := db.BindNamed(serverInsertStmt, server)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, stmt, args...)
	return err
}

func (db *serverStore) Update(ctx context.Context, server *autoscaler.Server) error {
	// before := server.Updated
	server.Updated = time.Now().Unix()
	stmt, args, err := db.BindNamed(serverUpdateStmt, server)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, stmt, args...)
	return err
}

func (db *serverStore) Delete(ctx context.Context, server *autoscaler.Server) error {
	stmt, args, err := db.BindNamed(serverDeleteStmt, server)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, stmt, args...)
	return err
}

func (db *serverStore) Purge(ctx context.Context, before int64) error {
	stmt, args, err := db.BindNamed(serverPurgeStmt, &autoscaler.Server{Stopped: before})
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, stmt, args...)
	return err
}

const serverFindStmt = `
SELECT 
 server_name
,server_id
,server_provider
,server_state
,server_image
,server_region
,server_size
,server_platform
,server_address
,server_capacity
,server_secret
,server_error
,server_ca_key
,server_ca_cert
,server_tls_key
,server_tls_cert
,server_created
,server_updated
,server_started
,server_stopped
FROM servers 
WHERE server_name=:server_name
`

const serverListStmt = `
SELECT 
 server_name
,server_id
,server_provider
,server_state
,server_image
,server_region
,server_size
,server_platform
,server_address
,server_capacity
,server_secret
,server_error
,server_ca_key
,server_ca_cert
,server_tls_key
,server_tls_cert
,server_created
,server_updated
,server_started
,server_stopped
FROM servers 
ORDER BY server_created ASC
`

const serverListStateStmt = `
SELECT 
 server_name
,server_id
,server_provider
,server_state
,server_image
,server_region
,server_size
,server_platform
,server_address
,server_capacity
,server_secret
,server_error
,server_ca_key
,server_ca_cert
,server_tls_key
,server_tls_cert
,server_created
,server_updated
,server_started
,server_stopped
FROM servers 
WHERE server_state=:server_state
ORDER BY server_created ASC
`

const serverInsertStmt = `
INSERT INTO servers (
 server_name
,server_id
,server_provider
,server_state
,server_image
,server_region
,server_size
,server_platform
,server_address
,server_capacity
,server_secret
,server_error
,server_ca_key
,server_ca_cert
,server_tls_key
,server_tls_cert
,server_created
,server_updated
,server_started
,server_stopped
) VALUES (
 :server_name
,:server_id
,:server_provider
,:server_state
,:server_image
,:server_region
,:server_size
,:server_platform
,:server_address
,:server_capacity
,:server_secret
,:server_error
,:server_ca_key
,:server_ca_cert
,:server_tls_key
,:server_tls_cert
,:server_created
,:server_updated
,:server_started
,:server_stopped
)
`

const serverUpdateStmt = `
UPDATE servers SET 
 server_id=:server_id
,server_provider=:server_provider
,server_state=:server_state
,server_image=:server_image
,server_region=:server_region
,server_size=:server_size
,server_platform=:server_platform
,server_address=:server_address
,server_capacity=:server_capacity
,server_secret=:server_secret
,server_error=:server_error
,server_ca_key=:server_ca_key
,server_ca_cert=:server_ca_cert
,server_tls_key=:server_tls_key
,server_tls_cert=:server_tls_cert
,server_updated=:server_updated
,server_started=:server_started
,server_stopped=:server_stopped
WHERE server_name=:server_name
`

const serverDeleteStmt = `
DELETE FROM servers WHERE server_name=:server_name
`

const serverPurgeStmt = `
DELETE FROM servers
WHERE server_state = 'stopped'
  AND server_stopped < :server_stopped
`
