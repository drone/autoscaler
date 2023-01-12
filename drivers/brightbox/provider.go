// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import (
	"sync"
	"text/template"

	"github.com/brightbox/gobrightbox"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
)

type provider struct {
	init sync.Once

	region string

	apiURL       string
	clientID     string
	clientSecret string

	image        string
	serverType   string
	serverGroups []string
	userdata     *template.Template

	client *gobrightbox.Client
}

func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			panic(err)
		}
	}

	p.region = "gb1"

	if p.apiURL == "" {
		p.apiURL = "https://api.gb1.brightbox.com"
	}
	if p.image == "" {
		p.image = "img-sm72o" // ubuntu-jammy-22.04-amd64-server
	}
	if p.serverType == "" {
		p.serverType = "typ-1ni84" // 8gb.ssd-ram-opt
	}
	if p.userdata == nil {
		p.userdata = userdata.T
	}

	return p
}
