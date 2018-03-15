// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"errors"
	"sync"
	"text/template"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/compute/v1"
)

var (
	defaultTags = []string{
		"allow-docker",
	}

	defaultScopes = []string{
		"https://www.googleapis.com/auth/devstorage.read_only",
		"https://www.googleapis.com/auth/logging.write",
		"https://www.googleapis.com/auth/monitoring.write",
		"https://www.googleapis.com/auth/trace.append",
	}
)

// provider implements a DigitalOcean provider.
type provider struct {
	init sync.Once

	diskSize int64
	diskType string
	image    string
	labels   map[string]string
	network  string
	project  string
	scopes   []string
	size     string
	tags     []string
	zone     string
	userdata *template.Template

	service *compute.Service
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) (autoscaler.Provider, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.diskSize == 0 {
		p.diskSize = 50
	}
	if p.diskType == "" {
		p.diskType = "pd-standard"
	}
	if p.zone == "" {
		p.zone = "us-central1-a"
	}
	if p.size == "" {
		p.size = "n1-standard-1"
	}
	if p.image == "" {
		p.image = "ubuntu-os-cloud/global/images/ubuntu-1604-xenial-v20170721"
	}
	if p.network == "" {
		p.network = "global/networks/default"
	}
	if len(p.tags) == 0 {
		p.tags = defaultTags
	}
	if len(p.scopes) == 0 {
		p.scopes = defaultScopes
	}
	if p.service == nil {
		client, err := google.DefaultClient(oauth2.NoContext, p.scopes...)
		if err != nil {
			return nil, err
		}
		p.service, err = compute.New(client)
		if err != nil {
			return nil, err
		}
	}
	if p.userdata == nil {
		d, err := userdata.DetectUserdata(p.image)
		p.userdata = d
		return p, err
	}
	return p, nil
}

func (p *provider) waitZoneOperation(ctx context.Context, name string) error {
	for {
		op, err := p.service.ZoneOperations.Get(p.project, p.zone, name).Context(ctx).Do()
		if err != nil {
			return err
		}
		if op.Error != nil {
			return errors.New(op.Error.Errors[0].Message)
		}
		if op.Status == "DONE" {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func (p *provider) waitGlogalOperation(ctx context.Context, name string) error {
	for {
		op, err := p.service.GlobalOperations.Get(p.project, name).Context(ctx).Do()
		if err != nil {
			return err
		}
		if op.Error != nil {
			return errors.New(op.Error.Errors[0].Message)
		}
		if op.Status == "DONE" {
			return nil
		}
		time.Sleep(time.Second)
	}
}
