// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var noContext = context.Background()

// connect opens a new test database connection.
func connect() (*sqlx.DB, error) {
	var (
		driver = "sqlite3"
		config = ":memory:"
	)
	if os.Getenv("DATABASE_DRIVER") != "" {
		driver = os.Getenv("DATABASE_DRIVER")
		config = os.Getenv("DATABASE_CONFIG")
	}
	return Connect(driver, config)
}

// locker returns a new text locker.
func locker() sync.Locker {
	driver := "sqlite3"
	if os.Getenv("DATABASE_DRIVER") != "" {
		driver = os.Getenv("DATABASE_DRIVER")
	}
	return NewLocker(driver)
}
