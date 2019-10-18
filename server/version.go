// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
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
	Commit  string `json:"commit,omitempty"`
}

// HandleVersion creates an http.HandlerFunc that returns the
// version number and build details.
func HandleVersion(source, version, commit string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := versionInfo{
			Source:  source,
			Version: version,
			Commit:  commit,
		}
		writeJSON(w, &data, 200)
	}
}
