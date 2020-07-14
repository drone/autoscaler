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
	_, desErr := p.getClient().DescribeInstances(describe)
	// if we are able to describe the instance it confirms the
	// instance still exists and could not be terminated. Return
	// an error so that the instance is flagged as being in an
	// error state and requires manual attention.
	if desErr == nil {
		logger.Errorln("describe instance was successful. instance still exists")
		return err
	}

	// if the ware unable to describe the instance because the
	// instance no longer exists, we can return a not found error.
	// this will result in the instance being deleted from the
	// system, since we will have confirmed it no longer exists.
	if awsErr, ok := desErr.(awserr.Error); ok {
		switch awsErr.Code() {
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdMalformed:
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		case ec2.UnsuccessfulInstanceCreditSpecificationErrorCodeInvalidInstanceIdNotFound:
			logger.Debugln("instance does not exist")
			return autoscaler.ErrInstanceNotFound
		}
	}

	// otherwise we return the original error returned when
	// attempting to delete the instance.
	return err
}
