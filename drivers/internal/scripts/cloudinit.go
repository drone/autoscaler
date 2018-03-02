// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"bytes"
	"text/template"
)

type TemplateOpts struct {
	Instance struct {
		Name string
	}
	Server struct {
		Host   string
		Secret string
	}
	Agent struct {
		Image    string
		Capacity int
	}
	Cadvisor struct {
		Disable bool
		Secret  string
	}
}

func GenerateCloudInit(tmpl string, opts TemplateOpts) (string, error) {
	buf := new(bytes.Buffer)
	err := template.Must(template.New("_").Funcs(funcs).Parse(tmpl)).Execute(buf, &opts)

	return buf.String(), err
}
