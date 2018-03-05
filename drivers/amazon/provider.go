// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"github.com/drone/autoscaler"

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
	retries int
	key     string
	region  string
	image   string
	size    string
	subnet  string
	groups  []string
	tags    map[string]string
}

func (p *provider) getClient() *ec2.EC2 {
	config := aws.NewConfig()
	config = config.WithRegion(p.region)
	config = config.WithMaxRetries(p.retries)
	return ec2.New(session.New(config))
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := &provider{
		retries: 25,
	}
	for _, opt := range opts {
		opt(p)
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
	return p
}
