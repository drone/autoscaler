// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package amazon

import (
	"context"
	"os"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	if os.Getenv("DRONE_FLAG_ALTERNATE_DESTROY") == "true" {
		return p.destroy2(ctx, instance)
	}

	logger := logger.FromContext(ctx).
		WithField("id", instance.ID).
		WithField("ip", instance.Address).
		WithField("name", instance.Name).
		WithField("zone", instance.Region)

	logger.Debugln("terminate instance")

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
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}
	if err != nil {
		logger.WithError(err).
			Errorln("cannot terminate instance")
		return err
	}

	logger.Debugln("terminated")

	return nil
}

func (p *provider) destroy2(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("id", instance.ID).
		WithField("ip", instance.Address).
		WithField("name", instance.Name).
		WithField("zone", instance.Region)

	logger.Debugln("terminate instance")

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.ID),
		},
	}
	_, err := p.getClient().TerminateInstances(input)
	if err == nil {
		logger.Debugln("terminated")
		return nil
	}

	// if terminate instance returns an error indicating
	// the instance no longer exists, return a not found
	// error.
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdMalformed:
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdNotFound:
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}
	if err != nil {
		logger.WithError(err).
			Errorln("cannot terminate instance")
	}

	logger.Debugln("describe instance")

	describe := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.ID),
		},
	}
	_, err = p.getClient().DescribeInstances(describe)
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdMalformed:
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdNotFound:
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}

	logger.WithError(err).
		Errorln("cannot describe instance")
	return err
}
