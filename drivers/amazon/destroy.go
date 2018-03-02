// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"context"

	"github.com/drone/autoscaler"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog/log"
)

// Destroy destroyes the DigitalOcean instance.
func (p *Provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	logger.Debug().
		Msg("shutdown agent")

	_, err := p.Provider.Execute(ctx, instance, teardownScript)
	if err != nil {
		// if we cannot gracefully shutdown the agent we should
		// still continue and destroy the instance. I think.
		logger.Error().
			Err(err).
			Msg("cannot shutdown agent")

		// TODO(bradrydzewski) we should snapshot the error logs
	}

	logger.Debug().
		Msg("terminate instsance")

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.ID),
		},
	}
	_, err = p.getClient().TerminateInstances(input)
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

var teardownScript = `
set -x;

sudo docker stop -t 3600 agent
sudo docker ps -a
`
