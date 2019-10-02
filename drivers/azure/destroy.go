// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"
	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"
	disks "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func (p *provider) getDiskClient() (disks.DisksClient, error) {
	disksClient := disks.NewDisksClient(p.subscriptionId)
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return disksClient, err
	}
	disksClient.Authorizer = auth
	return disksClient, nil
}

func (p *provider) deleteDisks(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("Azure Destroy:", instance.ID).
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	diskClient, err := p.getDiskClient()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	future, err := diskClient.Delete(ctx, p.resourceGroup, instance.ID + "-osdisk")
	err = future.WaitForCompletionRef(ctx, diskClient.Client)
	logger.Debug().Msg("waitforcompletionref")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}
	
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	err = future.WaitForCompletionRef(ctx, diskClient.Client)
	logger.Debug().Msg("waitforcompletionref")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}
	return nil
}

func (p *provider) deleteVM(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("Azure Destroy:", instance.ID).
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	vmClient, err := p.getClient()
	if err != nil {
		return err
	}
	
	future, err := vmClient.Delete(ctx, p.resourceGroup, instance.ID + "-vm")
	fmt.Println(future, err)
	logger.Debug().Msg("Deallocate call")

	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	logger.Debug().Msg("waitforcompletionref")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	result, err := future.Result(vmClient)
	fmt.Println(result, err)
	logger.Debug().Msg("future.result")
	return nil
}

func (p *provider) deleteIP(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("Azure Destroy:", instance.ID).
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	ipClient, err := p.getIPClient()
	if err != nil {
		return err
	}
	future, err := ipClient.Delete(ctx, p.resourceGroup, instance.ID + "-ip")
	fmt.Println(future, err)
	logger.Debug().Msg("Deallocate call")

	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	err = future.WaitForCompletionRef(ctx, ipClient.Client)
	logger.Debug().Msg("waitforcompletionref")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	result, err := future.Result(ipClient)
	fmt.Println(result, err)
	logger.Debug().Msg("future.result")
	return nil
}

func (p *provider) deleteNIC(ctx context.Context, instance *autoscaler.Instance) error {
	logger := log.Ctx(ctx).With().
		Str("Azure Destroy:", instance.ID).
		Str("id", instance.ID).
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Str("zone", instance.Region).
		Logger()

	nicClient, err := p.getNicClient()
	if err != nil {
		return err
	}
	future, err := nicClient.Delete(ctx, p.resourceGroup, instance.ID + "-nic")
	fmt.Println(future, err)
	logger.Debug().Msg("Deallocate call")

	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	err = future.WaitForCompletionRef(ctx, nicClient.Client)
	logger.Debug().Msg("waitforcompletionref")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot terminate instance")
		return err
	}

	result, err := future.Result(nicClient)
	fmt.Println(result, err)
	logger.Debug().Msg("future.result")
	return nil
}

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	err := p.deleteVM(ctx, instance)
	if err != nil {
		return err
	}

	err = p.deleteNIC(ctx, instance)
	if err != nil {
		return err
	}
	
	err = p.deleteIP(ctx, instance)
	if err != nil {
		return err
	}

	err = p.deleteDisks(ctx, instance)
	if err != nil {
		return err
	}

	return nil
}
