// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package web

import (
	"net/http"

	"github.com/drone/autoscaler/server/web/template"
)

// render writes the template to the response body.
func render(w http.ResponseWriter, t string, v interface{}) {
	w.Header().Set("Content-Type", "text/html")
	template.T.ExecuteTemplate(w, t, v)
}
