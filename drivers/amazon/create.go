// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package amazon

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	instance, err := p.create(ctx, opts, false)
	// if the instance was successfully provisioned,
	// return the instance.
	if err == nil {
		return instance, err
	}

	// if the instance was provisioned with errors,
	// return the instance and the error
	if instance != nil {
		return instance, err
	}

	// if the instance was not provisioned, and fallback
	// parameters were provided, retry using the fallback
	if p.sizeAlt != "" {
		instance, err = p.create(ctx, opts, true)
	}

	// if there is no fallback logic do not retry
	return instance, err
}

func (p *provider) create(ctx context.Context, opts autoscaler.InstanceCreateOpts, retry bool) (*autoscaler.Instance, error) {
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

	if p.volumeType == "io1" || p.volumeType == "io2" || p.volumeType == "gp3" {
		for _, blockDeviceMapping := range in.BlockDeviceMappings {
			if p.volumeIops > 0 {
				blockDeviceMapping.Ebs.Iops = aws.Int64(p.volumeIops)
			}
		}
	}

	if p.volumeType == "gp3" {
		for _, blockDeviceMapping := range in.BlockDeviceMappings {
			if p.volumeThroughput > 0 {
				blockDeviceMapping.Ebs.Throughput = aws.Int64(p.volumeThroughput)
			}
		}
	}

	logger := logger.FromContext(ctx).
		WithField("attempt", 1).
		WithField("region", p.region).
		WithField("image", p.image).
		WithField("size", p.size).
		WithField("name", opts.Name)

	// TODO(bradyrdzewski) instead of passing a re-try flag
	// and then setting parameters, we should instead accept
	// an struct that specifies the size, image and any other
	// alternate values that one may want to try

	// if this is our second attempt to create the instance,
	// re-create using the alternate instance size.
	if retry {
		in.InstanceType = aws.String(p.sizeAlt)

		// update the logger to reflect this is a retry.
		logger = logger.
			WithField("size", p.sizeAlt).
			WithField("attempt", 2)
	}

	logger.Debug("instance create")

	results, err := client.RunInstances(in)
	if err != nil {
		logger.WithError(err).
			Error("instance create failed")
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

	logger.WithField("name", instance.Name).
		Infoln("instance create success")

	// poll the amazon endpoint for server updates
	// and exit when a network address is allocated.
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.WithField("name", instance.Name).
				Debugln("instance network deadline exceeded")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute

			logger.WithField("name", instance.Name).
				Debugln("check instance network")

			desc, err := client.DescribeInstances(
				&ec2.DescribeInstancesInput{
					InstanceIds: []*string{
						amazonInstance.InstanceId,
					},
				},
			)
			if err != nil {
				logger.WithError(err).
					Warnln("instance details failed")
				continue
			}

			if len(desc.Reservations) == 0 {
				logger.Warnln("empty reservations in details")
				continue
			}
			if len(desc.Reservations[0].Instances) == 0 {
				logger.Warnln("empty instances in reservations")
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

	logger.
		WithField("name", instance.Name).
		WithField("ip", instance.Address).
		Debugln("instance network ready")

	return instance, nil
}
