// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"context"

	"github.com/drone/autoscaler"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog/log"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	logger.Debug().
		Msg("terminate instsance")

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
