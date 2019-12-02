// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package template

import "time"

//go:generate togo tmpl -func funcMap -format html

// mirros the func map in template.go
var funcMap = map[string]interface{}{
	"timestamp": func(v int64) string {
		return time.Unix(v, 0).UTC().Format("2006-01-02T15:04:05Z")
	},
}
