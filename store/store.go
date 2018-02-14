// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package store

import (
	"github.com/boltdb/bolt"
)

var serverKey = []byte("servers")

// New returns a new bolt database connection.
func New(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(serverKey)
		return nil
	})
	return db, err
}

// Must returns the bolt database connection and
// panics on error.
func Must(path string) *bolt.DB {
	db, err := New(path)
	if err != nil {
		panic(err)
	}
	return db
}
