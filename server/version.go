// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http"
)

// version information, loosely based on
// https://github.com/mozilla-services/Dockerflow
type versionInfo struct {
	Source  string `json:"source,omitempty"`
	Version string `json:"version,omitempty"`
}

// HandleVersion creates an http.HandlerFunc that returns the
// version number and build details.
func HandleVersion(source, version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := versionInfo{
			Source:  source,
			Version: version,
		}
		writeJSON(w, &data, 200)
	}
}
