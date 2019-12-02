// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package web

import (
	"net/http"
	"time"
)

// unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// http headers to disable caching.
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

// helper function to prevent http response caching.
func nocache(w http.ResponseWriter) {
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}
}
