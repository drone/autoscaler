// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package server

import (
	"net/http"
	"strconv"

	"github.com/dchest/uniuri"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/logger"

	"github.com/go-chi/chi"
)

// HandleServerList returns an http.HandlerFunc that writes
// the json-encoded server list to the the response body.
func HandleServerList(servers autoscaler.ServerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		list, err := servers.List(ctx)
		if err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Errorln("cannot get server list")
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
			logger.FromContext(ctx).
				WithError(err).
				WithField("server", name).
				Errorln("cannot get server")
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
		force, _ := strconv.ParseBool(r.FormValue("force"))

		server, err := servers.Find(ctx, name)
		if err != nil {
			logger.FromContext(ctx).
				WithError(err).
				WithField("server", name).
				Errorln("cannot get server")
			writeNotFound(w, err)
			return
		}

		// in some cases the server fails to create and is stuck
		// in an error state. In this case we force-delete from
		// the database.
		if server.State == autoscaler.StateError && (server.ID == "" || force) {
			logger.FromContext(ctx).
				WithField("server", server.Name).
				WithField("state", string(server.State)).
				WithField("force", force).
				Infoln("force delete server from database")

			err = servers.Delete(ctx, server)
			if err != nil {
				logger.FromContext(ctx).
					WithError(err).
					WithField("server", server.Name).
					Errorln("cannot delete instance")
				writeError(w, err)
				return
			}
			w.WriteHeader(204)
			return
		}

		logger.FromContext(ctx).
			WithField("server", server.Name).
			WithField("state", string(server.State)).
			WithField("force", force).
			Infoln("schedule server shutdown")

		server.State = autoscaler.StateShutdown
		err = servers.Update(ctx, server)
		if err != nil {
			logger.FromContext(ctx).
				WithError(err).
				WithField("server", server.Name).
				WithField("state", "shutdown").
				Errorln("cannot update server")
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
			logger.FromContext(ctx).
				WithError(err).
				Errorln("cannot persist server")
			writeError(w, err)
			return
		}
		writeJSON(w, server, 200)
	}
}
