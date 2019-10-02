package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func (p *provider) getSubnetsClient() (network.SubnetsClient, error) {
	subnetsClient := network.NewSubnetsClient(p.subscriptionId)
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return subnetsClient, err
	}
	subnetsClient.Authorizer = auth
	return subnetsClient, nil
}

// getVirtualNetworkSubnet returns an existing subnet from a virtual network
func (p *provider) getVirtualNetworkSubnet(ctx context.Context, vnetName string, subnetName string) (network.Subnet, error) {
	subnetsClient, err := p.getSubnetsClient()
	if err != nil {
		return network.Subnet{}, err
	}
	return subnetsClient.Get(ctx, p.resourceGroup, vnetName, subnetName, "")
}


func (p *provider) getIPClient() (network.PublicIPAddressesClient, error) {
	ipClient := network.NewPublicIPAddressesClient(p.subscriptionId)
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return ipClient, err
	}
	ipClient.Authorizer = auth
	return ipClient, nil
}

// CreatePublicIP creates a new public IP
func (p *provider) createPublicIP(ctx context.Context, ipName string) (ip network.PublicIPAddress, err error) {
	ipClient, err := p.getIPClient()
	if err != nil {
		return network.PublicIPAddress{}, err
	}
	future, err := ipClient.CreateOrUpdate(
		ctx,
		p.resourceGroup,
		ipName,
		network.PublicIPAddress{
			Name:     to.StringPtr(ipName),
			Location: to.StringPtr(p.location),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   network.IPv4,
				PublicIPAllocationMethod: network.Static,
			},
		},
	)

	if err != nil {
		return ip, fmt.Errorf("cannot create public ip address: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, ipClient.Client)
	if err != nil {
		return ip, fmt.Errorf("cannot get public ip address create or update future response: %v", err)
	}

	return future.Result(ipClient)
}

func (p *provider) getNicClient() (network.InterfacesClient, error) {
	nicClient := network.NewInterfacesClient(p.subscriptionId)
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nicClient, err
	}

	nicClient.Authorizer = auth
	return nicClient, nil
}

// CreateNIC creates a new network interface. The Network Security Group is not a required parameter
func (p *provider) CreateNIC(ctx context.Context, vnetName, subnetName, nsgName, ipName, nicName string) (nic network.Interface, ipres network.PublicIPAddress,  err error) {
	subnet, err := p.getVirtualNetworkSubnet(ctx, vnetName, subnetName)
	if err != nil {
		log.Fatalf("failed to get subnet: %v", err)
	}

	ip, err := p.createPublicIP(ctx, ipName)
	
	if err != nil {
		log.Fatalf("failed to get ip address: %v", err)
	}

	nicParams := network.Interface{
		Name:     to.StringPtr(nicName),
		Location: to.StringPtr(p.location),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				{
					Name: to.StringPtr("ipConfig1"),
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						Subnet:                    &subnet,
						PrivateIPAllocationMethod: network.Dynamic,
						PublicIPAddress:           &ip,
					},
				},
			},
		},
	}

	nicClient, err := p.getNicClient()
	if err != nil {
		return nic, ip, err
	}
	future, err := nicClient.CreateOrUpdate(ctx, p.resourceGroup, nicName, nicParams)
	if err != nil {
		return nic, ip, fmt.Errorf("cannot create nic: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, nicClient.Client)
	if err != nil {
		return nic, ip, fmt.Errorf("cannot get nic create or update future response: %v", err)
	}

	nicresponse, err := future.Result(nicClient)
	return nicresponse, ip, err
}
