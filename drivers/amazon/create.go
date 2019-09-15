// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/drone/autoscaler"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog/log"
)

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	p.init.Do(func() {
		p.setup(ctx)
	})

	buf := new(bytes.Buffer)
	err := p.userdata.Execute(buf, &opts)
	if err != nil {
		return nil, err
	}

	client := p.getClient()

	var iamProfile *ec2.IamInstanceProfileSpecification

	if p.iamProfileArn != "" {
		iamProfile = &ec2.IamInstanceProfileSpecification{
			Arn: &p.iamProfileArn,
		}
	}

	var marketOptions *ec2.InstanceMarketOptionsRequest

	if p.spotInstance == true {
		marketOptions = &ec2.InstanceMarketOptionsRequest{
			MarketType: aws.String("spot"),
		}
	}

	tags := createCopy(p.tags)
	tags["Name"] = opts.Name

	in := &ec2.RunInstancesInput{
		KeyName:               aws.String(p.key),
		ImageId:               aws.String(p.image),
		InstanceType:          aws.String(p.size),
		MinCount:              aws.Int64(1),
		MaxCount:              aws.Int64(1),
		InstanceMarketOptions: marketOptions,
		IamInstanceProfile:    iamProfile,
		UserData:              aws.String(base64.StdEncoding.EncodeToString(buf.Bytes())),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				AssociatePublicIpAddress: aws.Bool(!p.privateIP),
				DeviceIndex:              aws.Int64(0),
				SubnetId:                 aws.String(p.subnet),
				Groups:                   aws.StringSlice(p.groups),
			},
		},
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         convertTags(tags),
			},
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String(p.deviceName),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize:          aws.Int64(p.volumeSize),
					VolumeType:          aws.String(p.volumeType),
					DeleteOnTermination: aws.Bool(true),
				},
			},
		},
	}

	logger := log.Ctx(ctx).With().
		Str("region", p.region).
		Str("image", p.image).
		Str("size", p.size).
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

	// poll the amazon endpoint for server updates
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
				logger.Warn().
					Err(err).
					Msg("instance details failed")
				continue
			}

			if len(desc.Reservations) == 0 {
				logger.Warn().Msg("empty reservations in details")
				continue
			}
			if len(desc.Reservations[0].Instances) == 0 {
				logger.Warn().Msg("empty instances in reservations")
				continue
			}

			amazonInstance = desc.Reservations[0].Instances[0]

			if p.privateIP {
				if amazonInstance.PrivateIpAddress != nil {
					instance.Address = *amazonInstance.PrivateIpAddress
					break poller
				}
			}

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
