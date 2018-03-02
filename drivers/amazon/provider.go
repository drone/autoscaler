// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/internal/scripts"
	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	maxRetries = 50

	defaultDeviceName = "/dev/sda1"
	defaultImage      = "ami-66506c1c"
	defaultRootSize   = 16
	defaultVolumeType = "gp2"
)

// Provider defines the Amazon provider.
type Provider struct {
	autoscaler.Provider

	config config.Config
	signer ssh.Signer
}

func (p *Provider) getClient() *ec2.EC2 {
	return buildClient(p.config)
}

func buildClient(conf config.Config) *ec2.EC2 {
	config := aws.NewConfig()
	config = config.WithRegion(conf.Amazon.Region)
	config = config.WithMaxRetries(maxRetries)
	return ec2.New(session.New(config))
}

func (p *Provider) setupScriptOpts(instance *autoscaler.Instance) scripts.SetupOpts {
	opts := scripts.SetupOpts{}
	opts.Server.Host = p.config.Agent.Host
	opts.Server.Secret = p.config.Agent.Token
	opts.Agent.Image = p.config.Agent.Image
	opts.Agent.Capacity = p.config.Agent.Concurrency
	opts.Instance.Addr = instance.Address
	opts.Instance.Name = instance.Name
	opts.Cadvisor.Disable = false
	opts.Cadvisor.Secret = instance.Secret
	return opts
}
