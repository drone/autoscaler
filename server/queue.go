// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/rs/zerolog/hlog"
)

// HandleQueueList returns an http.HandlerFunc that writes
// the json-encoded queue details to the the response body.
func HandleQueueList(client drone.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := client.BuildQueue()
		if err != nil {
			hlog.FromRequest(r).
				Error().
				Err(err).
				Msg("cannot get build queue")
			writeError(w, err)
			return
		}
		writeJSON(w, list, 200)
	}
}
