// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package web provides HTTP handlers that expose pipeline
// state and status.
package web

import (
	"net/http"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger/history"
)

// HandleServers returns a http.HandlerFunc that displays a
// list of activate servers.
func HandleServers(servers autoscaler.ServerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nocache(w)
		items, _ := servers.List(r.Context())
		filtered := []*autoscaler.Server{}
		for _, item := range items {
			if item.State != autoscaler.StateStopped {
				filtered = append(filtered, item)
			}
		}
		render(w, "index.tmpl", struct {
			Items []*autoscaler.Server
		}{filtered})
	}
}

// HandleLogging returns a http.HandlerFunc that displays a
// list recent log entries.
func HandleLogging(t *history.Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nocache(w)
		render(w, "logs.tmpl", struct {
			Entries []*history.Entry
		}{t.Entries()})
	}
}
