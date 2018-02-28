// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/drone/autoscaler"
)

// NewServerStore returns a new server store.
func NewServerStore(db *bolt.DB) autoscaler.ServerStore {
	return &serverStore{db}
}

type serverStore struct {
	*bolt.DB
}

func (db *serverStore) Find(ctx context.Context, name string) (*autoscaler.Server, error) {
	key := []byte(name)
	val := new(autoscaler.Server)
	err := db.DB.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(serverKey).Get(key)
		if len(data) == 0 {
			return autoscaler.ErrServerNotFound
		}
		return json.Unmarshal(data, val)
	})
	return val, err
}

func (db *serverStore) List(ctx context.Context) ([]*autoscaler.Server, error) {
	items := []*autoscaler.Server{}
	err := db.DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(serverKey).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := new(autoscaler.Server)
			json.Unmarshal(v, item)
			items = append(items, item)
		}
		return nil
	})
	return items, err
}

func (db *serverStore) ListState(ctx context.Context, state autoscaler.ServerState) ([]*autoscaler.Server, error) {
	items := []*autoscaler.Server{}
	all, err := db.List(ctx)
	if err != nil {
		return items, err
	}
	for _, item := range all {
		if item.State == state {
			items = append(items, item)
		}
	}
	return items, err
}

func (db *serverStore) Create(ctx context.Context, server *autoscaler.Server) error {
	return db.Update(ctx, server)
}

func (db *serverStore) Update(ctx context.Context, server *autoscaler.Server) error {
	key := []byte(server.Name)
	val, _ := json.Marshal(server)
	return db.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(serverKey).Put(key, val)
	})
}

func (db *serverStore) Delete(ctx context.Context, server *autoscaler.Server) error {
	key := []byte(server.Name)
	return db.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(serverKey).Delete(key)
	})
}
