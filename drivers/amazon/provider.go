// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
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

type provider struct {
	init sync.Once

	deviceName       string
	volumeSize       int64
	volumeType       string
	volumeIops       int64
	volumeThroughput int64
	retries          int
	key              string
	region           string
	image            string
	privateIP        bool
	userdata         *template.Template
	size             string
	sizeAlt          string
	subnets          []string
	groups           []string
	tags             map[string]string
	iamProfileArn    string
	spotInstance     bool
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
		p.size = "t3.medium"
	}
	if p.image == "" {
		p.image = defaultImage(p.region)
	}
	if p.deviceName == "" {
		p.deviceName = "/dev/sda1"
	}
	if p.volumeSize == 0 {
		p.volumeSize = 32
	}
	if p.volumeType == "" {
		p.volumeType = "gp2"
	}
	if (p.volumeType == "io1" || p.volumeType == "io2") && p.volumeIops == 0 {
		p.volumeIops = 100
	}
	if p.volumeType == "gp3" && p.volumeIops == 0 {
		p.volumeIops = 3000 // 3000 is the minimum for gp3
	}
	if p.volumeType == "gp3" && p.volumeThroughput == 0 {
		p.volumeThroughput = 125 // 125 is the minimum for gp3
	}
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	return p
}
