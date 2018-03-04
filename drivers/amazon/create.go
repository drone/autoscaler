// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"
)

// Create creates the DigitalOcean instance.
func (p *Provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	client := p.getClient()

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
		Str("region", p.config.Amazon.Region).
		Str("image", defaultImage).
		Str("size", p.config.Amazon.Instance).
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

	amazonInstance := results.Instances[0]

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderAmazon,
		ID:       *amazonInstance.InstanceId,
		Name:     opts.Name,
		Size:     *amazonInstance.InstanceType,
		Region:   *amazonInstance.Placement.AvailabilityZone,
		Image:    *amazonInstance.ImageId,
	}

	logger.Info().
		Str("name", instance.Name).
		Msg("instance create success")

	// poll the digitalocean endpoint for server updates
	// and exit when a network address is allocated.
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.Debug().
				Str("name", instance.Name).
				Msg("instance network deadline exceeded")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute

			logger.Debug().
				Str("name", instance.Name).
				Msg("check instance network")

			desc, err := client.DescribeInstances(
				&ec2.DescribeInstancesInput{
					InstanceIds: []*string{
						amazonInstance.InstanceId,
					},
				},
			)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("instance details failed")
				return nil, err
			}
			amazonInstance = desc.Reservations[0].Instances[0]

			if amazonInstance.PublicIpAddress != nil {
				instance.Address = *amazonInstance.PublicIpAddress
				break poller
			}
		}
	}

	logger.Debug().
		Str("name", instance.Name).
		Str("ip", instance.Address).
		Msg("instance network ready")

	return instance, nil
}
