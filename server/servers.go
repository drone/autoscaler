// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/hlog"
)

// HandleServerList returns an http.HandlerFunc that writes
// the json-encoded server list to the the response body.
func HandleServerList(servers autoscaler.ServerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		list, err := servers.List(ctx)
		if err != nil {
			hlog.FromRequest(r).
				Error().
				Err(err).
				Msg("cannot get server list")
			writeError(w, err)
			return
		}
		writeJSON(w, list, 200)
	}
}

// HandleServerFind returns an http.HandlerFunc that finds
// and writes the json-encoded server to the the response body.
func HandleServerFind(servers autoscaler.ServerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := chi.URLParam(r, "name")
		server, err := servers.Find(ctx, name)
		if err != nil {
			hlog.FromRequest(r).
				Error().
				Err(err).
				Str("server", name).
				Msg("cannot get server")
			writeNotFound(w, err)
			return
		}
		writeJSON(w, server, 200)
	}
}

// HandleServerDelete returns an http.HandlerFunc that destroys
// and then deletes the named server.
func HandleServerDelete(
	servers autoscaler.ServerStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := chi.URLParam(r, "name")
		server, err := servers.Find(ctx, name)
		if err != nil {
			hlog.FromRequest(r).
				Error().
				Err(err).
				Str("server", name).
				Msg("cannot get server")
			writeNotFound(w, err)
			return
		}
		server.State = autoscaler.StateShutdown
		err = servers.Update(ctx, server)
		if err != nil {
			hlog.FromRequest(r).
				Error().
				Err(err).
				Str("server", name).
				Msg("cannot update server")
			writeError(w, err)
			return
		}
		writeJSON(w, server, 200)
	}
}

// HandleServerCreate returns an http.HandlerFunc that creates
// and a new server.
func HandleServerCreate(
	servers autoscaler.ServerStore,
	config config.Config,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		server := &autoscaler.Server{
			Name:     "agent-" + uniuri.NewLen(8),
			State:    autoscaler.StatePending,
			Capacity: config.Agent.Concurrency,
		}
		err := servers.Create(ctx, server)
		if err != nil {
			hlog.FromRequest(r).
				Error().
				Err(err).
				Msg("cannot persist server")
			writeError(w, err)
			return
		}
		writeJSON(w, server, 200)
	}
}
