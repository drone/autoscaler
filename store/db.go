// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"database/sql"
	"time"

	ddl "github.com/drone/autoscaler/store/migrate"

	"github.com/jmoiron/sqlx"
)

var noContext = context.Background()

// Connect to a database and verify with a ping.
func Connect(driver, datasource string) (*sqlx.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}
	switch driver {
	case "postgres":
		db.SetMaxIdleConns(0)
	case "mysql":
		db.SetMaxIdleConns(0)
	case "sqlite3":
		db.SetMaxOpenConns(1)
	}
	dbx := sqlx.NewDb(db, driver)
	if err := pingDatabase(dbx); err != nil {
		return nil, err
	}
	if err := setupDatabase(dbx); err != nil {
		return nil, err
	}
	return dbx, nil
}

// Must is a helper function that wraps a call to Connect
// and panics if the error is non-nil.
func Must(db *sqlx.DB, err error) *sqlx.DB {
	if err != nil {
		panic(err)
	}
	return db
}

// helper function to ping the database with backoff to ensure
// a connection can be established before we proceed with the
// database setup and migration.
func pingDatabase(db *sqlx.DB) (err error) {
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}

// helper function to setup the databsae by performing automated
// database migration steps.
func setupDatabase(db *sqlx.DB) error {
	return ddl.Migrate(db)
}
