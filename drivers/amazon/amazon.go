// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"context"
	"time"

	"github.com/dchest/uniuri"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/internal/scripts"
	"github.com/drone/autoscaler/drivers/internal/sshutil"
	"github.com/rs/zerolog/log"

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
	config config.Config
}

// Create creates the Amazon instance.
func (p *Provider) Create(ctx context.Context, opts *autoscaler.ServerOpts) (*autoscaler.Server, error) {
	client := p.getClient()

	signer, err := sshutil.ParsePrivateKey(p.config.Amazon.SSHKey)
	if err != nil {
		return nil, err
	}

	in := &ec2.RunInstancesInput{
		KeyName:      aws.String(p.config.Amazon.SSHKeyName),
		ImageId:      aws.String(defaultImage),
		InstanceType: aws.String(p.config.Amazon.Instance),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				AssociatePublicIpAddress: aws.Bool(true),
				DeviceIndex:              aws.Int64(0),
				SubnetId:                 aws.String(p.config.Amazon.SubnetID),
				Groups:                   aws.StringSlice(p.config.Amazon.SecurityGroup),
			},
		},
	}

	logger := log.Ctx(ctx).With().
		Str("image", *in.ImageId).
		Str("type", *in.InstanceType).
		Str("name", opts.Name).
		Logger()

	logger.Debug().
		Msg("instance create")

	results, err := client.RunInstances(in)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("instance create failed")
		return nil, err
	}

	instance := results.Instances[0]

	// tag the instance with user-defined tags.
	if tags := p.config.Amazon.Tags; tags != nil && len(tags) != 0 {
		p.getClient().CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{instance.InstanceId},
			Tags:      convertTags(p.config.Amazon.Tags),
		})
	}

	server := &autoscaler.Server{
		Provider: autoscaler.ProviderAmazon,
		UID:      *instance.InstanceId,
		Name:     opts.Name,
		Size:     *in.InstanceType,
		Region:   *instance.Placement.AvailabilityZone,
		Image:    *in.ImageId,
		Capacity: opts.Capacity,
		Secret:   opts.Secret,
	}

	// wait for the server to be available
	for {
		logger.Debug().
			Str("name", server.Name).
			Msg("instance network check")

		desc, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{instance.InstanceId},
		})
		if err != nil {
			logger.Error().
				Err(err).
				Msg("instance details failed")
			return nil, err
		}
		instance = desc.Reservations[0].Instances[0]

		if instance.PublicIpAddress == nil {
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	server.Address = *instance.PublicIpAddress

	logger.Debug().
		Str("name", server.Name).
		Str("ip", server.Address).
		Msg("instance network address assigned")

	// ping the server in a loop until we can successfully
	// authenticate.
	for {
		logger.Debug().
			Str("name", server.Name).
			Str("ip", server.Address).
			Str("port", "22").
			Str("user", "ubuntu").
			Msg("ping server")
		_, err = sshutil.Execute(server.Address, "22", "ubuntu", "whoami", signer)
		if err == nil {
			break
		} else {
			time.Sleep(time.Minute)
		}
	}

	server.Secret = uniuri.New()
	server.Created = time.Now().Unix()
	server.Updated = time.Now().Unix()

	script, err := scripts.GenerateInstall(p.config, server)
	if err != nil {
		return server, err
	}

	logger.Debug().
		Str("name", server.Name).
		Str("ip", server.Address).
		Msg("install agent")

	out, err := sshutil.Execute(server.Address, "22", "ubuntu", script, signer)
	server.Logs = string(out)
	if err != nil {
		logger.Error().
			Err(err).
			Str("name", server.Name).
			Str("ip", server.Address).
			Msg("install failed")
		return server, err
	}

	logger.Debug().
		Str("name", server.Name).
		Str("ip", server.Address).
		Msg("install complete")

	return server, nil
}

// Destroy terminates the AWS instance.
func (p *Provider) Destroy(ctx context.Context, server *autoscaler.Server) error {
	logger := log.Ctx(ctx).With().
		Str("region", server.Region).
		Str("image", server.Image).
		Str("size", server.Size).
		Str("name", server.Name).
		Logger()

	script, err := scripts.GenerateTeardown(p.config)
	if err != nil {
		return err
	}

	signer, err := sshutil.ParsePrivateKey(p.config.Amazon.SSHKey)
	if err != nil {
		return err
	}

	logger.Debug().
		Msg("teardown instance")

	_, err = sshutil.Execute(server.Address, "22", "ubuntu", script, signer)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("teardown failed")
		return err
	}

	logger.Debug().
		Msg("teardown instance complete")

	logger.Debug().
		Msg("terminate instance")

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(server.UID),
		},
	}
	_, err = p.getClient().TerminateInstances(input)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("terminate instance failed")
		return err
	}

	logger.Debug().
		Msg("instance terminated")

	return nil
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
