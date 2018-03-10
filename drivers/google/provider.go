// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"errors"
	"text/template"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"

	"google.golang.org/api/compute/v1"
)

// provider implements a DigitalOcean provider.
type provider struct {
	diskSize int64
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
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.diskSize == 0 {
		p.diskSize = 50
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
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	if len(p.scopes) == 0 {
		p.scopes = []string{
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring.write",
		}
	}
	return p
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
		time.Sleep(1 * time.Second)
	}
}
