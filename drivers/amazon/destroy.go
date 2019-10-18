// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package amazon

import (
	"context"
	"os"

	"github.com/drone/autoscaler"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog/log"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	if os.Getenv("DRONE_FLAG_ALTERNATE_DESTROY") == "true" {
		return p.destroy2(ctx, instance)
	}

	logger := log.Ctx(ctx).With().
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	logger.Debug().
		Msg("terminate instance")

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.ID),
		},
	}
	_, err := p.getClient().TerminateInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdMalformed:
			fallthrough
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdNotFound:
			logger.Debug().Msg("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	logger.Debug().
		Msg("terminated")

	return nil
}

func (p *provider) destroy2(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	logger.Debug().
		Msg("terminate instance")

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.ID),
		},
	}
	_, err := p.getClient().TerminateInstances(input)
	if err == nil {
		logger.Debug().Msg("terminated")
		return nil
	}

	// if terminate instance returns an error indicating
	// the instance no longer exists, return a not found
	// error.
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdMalformed:
			logger.Debug().Msg("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdNotFound:
			logger.Debug().Msg("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
	}

	logger.Debug().
		Msg("describe instance")

	describe := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.ID),
		},
	}
	_, err = p.getClient().DescribeInstances(describe)
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdMalformed:
			logger.Debug().Msg("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdNotFound:
			logger.Debug().Msg("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}

	logger.Error().
		Err(err).
		Msg("cannot describe instance")
	return err
}
