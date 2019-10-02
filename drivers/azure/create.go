// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package azure

import (
	"context"
	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"
	"fmt"
	"time"

	"bytes"
	"encoding/base64"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
)

func (p *provider) computeMachineRequest(ctx context.Context, base string, opts autoscaler.InstanceCreateOpts) (compute.VirtualMachine, *string, error) {
	buf := new(bytes.Buffer)
	err := p.userdata.Execute(buf, &opts)
	cloudInit := base64.StdEncoding.EncodeToString(buf.Bytes())

	nic, ip, err := p.CreateNIC(ctx, p.vnet, p.subnet, p.nsg, base+"-ip", base+"-nic")
	if err != nil {
		return compute.VirtualMachine{}, nil, err
	}

	return compute.VirtualMachine{
		Location: to.StringPtr(p.location),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: p.vmSize,
			},
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Name:         to.StringPtr(base + "-" + "osdisk"),
					CreateOption: compute.DiskCreateOptionTypesFromImage,
					DiskSizeGB:   to.Int32Ptr(p.volumeSize),
				},
				ImageReference: &compute.ImageReference{
					Offer:     to.StringPtr(p.imageOffer),
					Publisher: to.StringPtr(p.imagePublisher),
					Sku:       to.StringPtr(p.imageSKU),
					Version:   to.StringPtr(p.imageVersion),
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: nic.ID,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
			OsProfile: &compute.OSProfile{
				ComputerName:  to.StringPtr(base + "-" + "vm"),
				AdminUsername: to.StringPtr(p.adminUsername),
				AdminPassword: to.StringPtr(p.adminPassword),
				CustomData:    &cloudInit,
			},
		},
	}, ip.PublicIPAddressPropertiesFormat.IPAddress, nil
}

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	logger := log.Ctx(ctx)
	logger.Debug().Msg("Azure creating an instance.")

	client, err := p.getClient()
	if err != nil {
		return nil, err
	}
	
	base := opts.Name
	computeRequest, ipaddress, err := p.computeMachineRequest(ctx, base, opts)
	if err != nil {
		return nil, err
	}

	future, err := client.CreateOrUpdate(
		ctx,
		p.resourceGroup,
		base+"-vm",
		computeRequest,
	)
	
	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return nil, err
	}

	azureInstance, err := future.Result(client)
	client.Start(ctx, p.resourceGroup, p.vmName)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Starting sleep: %v\n", time.Now().Unix())
	time.Sleep(180 * time.Second)
	fmt.Printf("Ending sleep: %v\n", time.Now().Unix())

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderAzure,
		Address:  *ipaddress,
		ID:       base, // Set it to name because required for deallocating.
		Size:     string(azureInstance.HardwareProfile.VMSize),
		Region:   *azureInstance.Location,
		Image:    *azureInstance.StorageProfile.ImageReference.Sku,
	}

	return instance, nil
}
