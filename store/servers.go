// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/drone/autoscaler"

	"github.com/avast/retry-go"
	"github.com/jmoiron/sqlx"
)

// NewServerStore returns a new server store.
func NewServerStore(db *sqlx.DB, mu sync.Locker) autoscaler.ServerStore {
	return &serverStore{mu, db}
}

type serverStore struct {
	mu sync.Locker
	db *sqlx.DB
}

func (s *serverStore) Find(_ context.Context, name string) (*autoscaler.Server, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dest := &autoscaler.Server{Name: name}
	stmt, args, err := s.db.BindNamed(serverFindStmt, dest)
	if err != nil {
		return nil, err
	}
	err = s.db.GetContext(noContext, dest, stmt, args...)
	return dest, err
}

func (s *serverStore) List(_ context.Context) ([]*autoscaler.Server, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dest := []*autoscaler.Server{}
	err := s.db.SelectContext(noContext, &dest, serverListStmt)
	return dest, err
}

func (s *serverStore) ListState(_ context.Context, state autoscaler.ServerState) ([]*autoscaler.Server, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dest := []*autoscaler.Server{}
	stmt, args, err := s.db.BindNamed(serverListStateStmt, map[string]interface{}{"server_state": state})
	if err != nil {
		return nil, err
	}
	err = s.db.SelectContext(noContext, &dest, stmt, args...)
	if err == sql.ErrNoRows {
		return dest, nil
	}
	return dest, err
}

func (s *serverStore) Create(_ context.Context, server *autoscaler.Server) error {
	return retry.Do(
		func() error {
			if err := s.create(server); isConnReset(err) {
				return err
			} else {
				return retry.Unrecoverable(err)
			}
		},
		retry.Attempts(5),
		retry.MaxDelay(time.Second*5),
		retry.LastErrorOnly(true),
	)
}

func (s *serverStore) create(server *autoscaler.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	server.Created = time.Now().Unix()
	server.Updated = time.Now().Unix()
	stmt, args, err := s.db.BindNamed(serverInsertStmt, server)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(noContext, stmt, args...)
	return err
}

func (s *serverStore) Update(_ context.Context, server *autoscaler.Server) error {
	return retry.Do(
		func() error {
			if err := s.update(server); isConnReset(err) {
				return err
			} else {
				return retry.Unrecoverable(err)
			}
		},
		retry.Attempts(5),
		retry.MaxDelay(time.Second*5),
		retry.LastErrorOnly(true),
	)
}

func (s *serverStore) update(server *autoscaler.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	server.Updated = time.Now().Unix()
	stmt, args, err := s.db.BindNamed(serverUpdateStmt, server)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(noContext, stmt, args...)
	return err
}

func (s *serverStore) Busy(_ context.Context, server *autoscaler.Server) error {
	return retry.Do(
		func() error {
			if err := s.busy(server); isConnReset(err) {
				return err
			} else {
				return retry.Unrecoverable(err)
			}
		},
		retry.Attempts(5),
		retry.MaxDelay(time.Second*5),
		retry.LastErrorOnly(true),
	)
}

func (s *serverStore) busy(server *autoscaler.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	server.LastBusy = time.Now().Unix()
	stmt, args, err := s.db.BindNamed(serverUpdateStmt, server)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(noContext, stmt, args...)
	return err
}

func (s *serverStore) Delete(_ context.Context, server *autoscaler.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stmt, args, err := s.db.BindNamed(serverDeleteStmt, server)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(noContext, stmt, args...)
	return err
}

func (s *serverStore) Purge(_ context.Context, before int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stmt, args, err := s.db.BindNamed(serverPurgeStmt, &autoscaler.Server{Stopped: before})
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(noContext, stmt, args...)
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
,server_lastbusy
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
,server_lastbusy
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
,server_lastbusy
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
,:server_lastbusy
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
,server_lastbusy=:server_lastbusy
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
