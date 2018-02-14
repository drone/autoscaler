// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HandleMetrics returns an http.HandlerFunc that writes
// metrics to the response body in plain text format.
func HandleMetrics(token string) http.HandlerFunc {
	handler := promhttp.Handler()
	return func(w http.ResponseWriter, r *http.Request) {
		// if a bearer token is not configured we should
		// just server the http request.
		if token == "" {
			handler.ServeHTTP(w, r)
			return
		}
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, errInvalidToken.Error(), 401)
			return
		}
		if header != "Bearer "+token {
			http.Error(w, errInvalidToken.Error(), 401)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
