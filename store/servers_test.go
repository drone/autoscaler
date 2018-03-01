// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/drone/autoscaler"
)

func TestServer(t *testing.T) {
	conn, err := connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	store := NewServerStore(conn).(*serverStore)
	t.Run("Create", testServerCreate(store))
	t.Run("Find", testServerFind(store))
	t.Run("List", testServerList(store))
	t.Run("ListState", testServerListState(store))
	t.Run("Update", testServerUpdate(store))
	t.Run("Delete", testServerDelete(store))
	t.Run("Purge", testServerPurge(store))
}

func testServerCreate(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		server := &autoscaler.Server{
			Provider: autoscaler.ProviderGoogle,
			State:    autoscaler.StateRunning,
			Name:     "i-5203422c",
			Address:  "54.194.252.215",
			Capacity: 2,
			Created:  time.Now().Unix(),
			Updated:  time.Now().Unix(),
		}
		err := store.Create(context.TODO(), server)
		if err != nil {
			t.Error(err)
		}
	}
}

func testServerFind(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		server, err := store.Find(context.TODO(), "i-5203422c")
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Fields", testServer(server))
		}
	}
}

func testServerList(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		servers, err := store.List(context.TODO())
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := len(servers), 1; got != want {
			t.Errorf("Want server count %d, got %d", want, got)
		} else {
			t.Run("Fields", testServer(servers[0]))
		}
	}
}

func testServerListState(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		// seed the database with two servers with shutdown state.
		// to confirm we can list servers by state. These will be
		// used in a subsequent purge test.
		store.Create(context.TODO(), &autoscaler.Server{
			Provider: autoscaler.ProviderGoogle,
			State:    autoscaler.StateStopped,
			Name:     "agent-123456789",
		})
		store.Create(context.TODO(), &autoscaler.Server{
			Provider: autoscaler.ProviderGoogle,
			State:    autoscaler.StateStopped,
			Name:     "agent-987654321",
		})
		servers, err := store.ListState(context.TODO(), autoscaler.StateStopped)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := len(servers), 2; got != want {
			t.Errorf("Want server count %d, got %d", want, got)
		}
	}
}

func testServerUpdate(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		server := &autoscaler.Server{
			Provider: autoscaler.ProviderGoogle,
			Name:     "i-5203422c",
			Address:  "54.194.252.215",
			Capacity: 2,
			Created:  time.Now().Unix(),
			Updated:  time.Now().Unix(),
		}
		err := store.Update(context.TODO(), server)
		if err != nil {
			t.Error(err)
			return
		}
		updated, err := store.Find(context.TODO(), server.Name)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := updated.Capacity, server.Capacity; got != want {
			t.Errorf("Want updated capacity %d, got %d", want, got)
		}
	}
}

func testServerDelete(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := store.Find(context.TODO(), "i-5203422c")
		if err != nil {
			t.Error(err)
			return
		}

		err = store.Delete(context.TODO(), &autoscaler.Server{Name: "i-5203422c"})
		if err != nil {
			t.Error(err)
			return
		}

		_, err = store.Find(context.TODO(), "i-5203422c")
		if got, want := err, sql.ErrNoRows; got != want {
			t.Errorf("Want ErrNoRows, got %s", got)
		}
	}
}

func testServerPurge(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		// this test attempts to purge the database of all
		// servers with a state of stopped. The database was
		// seeded with stopped servers in testServerListState.
		before, _ := store.List(context.TODO())
		if got, want := len(before), 2; got != want {
			t.Errorf("Want %d servers, got %d", want, got)
			return
		}

		err := store.Purge(context.TODO(), time.Now().Unix()+1)
		if err != nil {
			t.Error(err)
			return
		}

		after, err := store.List(context.TODO())
		if err != nil {
			t.Error(err)
			return
		}

		if got, want := len(after), 0; got != want {
			t.Errorf("Want 0 remaining servers, got %d", got)
		}
	}
}

func testServer(server *autoscaler.Server) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := server.Name, "i-5203422c"; got != want {
			t.Errorf("Want server Name %q, got %q", want, got)
		}
		if got, want := server.State, autoscaler.StateRunning; got != want {
			t.Errorf("Want server State %v, got %v", want, got)
		}
		if got, want := server.Address, "54.194.252.215"; got != want {
			t.Errorf("Want server Address %q, got %q", want, got)
		}
		if got, want := server.Capacity, 2; got != want {
			t.Errorf("Want server Capacity %d, got %d", want, got)
		}
		if got, want := server.Provider, autoscaler.ProviderGoogle; got != want {
			t.Errorf("Want server Provider %v, got %v", want, got)
		}
	}
}
