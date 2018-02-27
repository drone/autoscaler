// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/drone/autoscaler"
)

// HandleScalerPause returns an http.HandlerFunc that pauses
// automatic scaling.
func HandleScalerPause(scaler autoscaler.Scaler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scaler.Pause()
		w.WriteHeader(204)
	}
}

// HandleScalerResume returns an http.HandlerFunc that resumed
// automatic scaling.
func HandleScalerResume(scaler autoscaler.Scaler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scaler.Resume()
		w.WriteHeader(204)
	}
}
