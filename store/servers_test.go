// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/drone/autoscaler"
)

func TestServer(t *testing.T) {
	temp := tempfile()
	defer os.Remove(temp)

	t.Logf("create boltdb database %s", temp)

	db := Must(temp)
	store := NewServerStore(db).(*serverStore)
	t.Run("Create", testServerCreate(store))
	t.Run("Find", testServerFind(store))
	t.Run("List", testServerList(store))
	t.Run("Update", testServerUpdate(store))
	t.Run("Delete", testServerDelete(store))
}

func testServerCreate(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		server := &autoscaler.Server{
			Provider: autoscaler.ProviateGoogle,
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

func testServerUpdate(store *serverStore) func(t *testing.T) {
	return func(t *testing.T) {
		server := &autoscaler.Server{
			Provider: autoscaler.ProviateGoogle,
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
		if got, want := err, autoscaler.ErrServerNotFound; got != want {
			t.Errorf("Want ErrServerNotFound, got %s", got)
		}
	}
}

func testServer(server *autoscaler.Server) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := server.Name, "i-5203422c"; got != want {
			t.Errorf("Want server Name %q, got %q", want, got)
		}
		if got, want := server.Address, "54.194.252.215"; got != want {
			t.Errorf("Want server Address %q, got %q", want, got)
		}
		if got, want := server.Capacity, 2; got != want {
			t.Errorf("Want server Capacity %d, got %d", want, got)
		}
		if got, want := server.Provider, autoscaler.ProviateGoogle; got != want {
			t.Errorf("Want server Provider %v, got %v", want, got)
		}
	}
}
