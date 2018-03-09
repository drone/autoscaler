// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	defaultDeviceName = "/dev/sda1"
	defaultImage      = "ami-66506c1c"
	defaultRootSize   = 16
	defaultVolumeType = "gp2"
)

type provider struct {
	init sync.Once

	retries   int
	key       string
	region    string
	image     string
	privateIP bool
	userdata  *template.Template
	size      string
	subnet    string
	groups    []string
	tags      map[string]string
}

func (p *provider) getClient() *ec2.EC2 {
	config := aws.NewConfig()
	config = config.WithRegion(p.region)
	config = config.WithMaxRetries(p.retries)
	return ec2.New(session.New(config))
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.retries == 0 {
		p.retries = 10
	}
	if p.region == "" {
		p.region = "us-east-1"
	}
	if p.size == "" {
		p.size = "t2.medium"
	}
	if p.image == "" {
		p.image = "ami-66506c1c"
	}
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	return p
}
