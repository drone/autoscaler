// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/drone/autoscaler"
)

type varz struct {
	Paused bool `json:"paused"`
}

// HandleVarz creates an http.HandlerFunc that returns system
// configuration and runtime information.
func HandleVarz(engine autoscaler.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := varz{
			Paused: engine.Paused(),
		}
		writeJSON(w, &data, 200)
	}
}
