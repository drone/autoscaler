// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/drone/autoscaler"
)

// HandleEnginePause returns an http.HandlerFunc that pauses
// scaling engine.
func HandleEnginePause(engine autoscaler.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		engine.Pause()
		w.WriteHeader(204)
	}
}

// HandleEngineResume returns an http.HandlerFunc that resumes
// scaling engine.
func HandleEngineResume(engine autoscaler.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		engine.Resume()
		w.WriteHeader(204)
	}
}
