// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package azure

import (
	"sync"
	"github.com/drone/autoscaler"
	"text/template"

	"github.com/drone/autoscaler/drivers/internal/userdata"

	// Azure apis 
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type provider struct {
	init sync.Once

	subscriptionId string
	resourceGroup  string
	image          string
	location       string
	adminUsername  string
	adminPassword  string
	vmName         string
	sshKey         string
	vmSize         compute.VirtualMachineSizeTypes
	userdata      *template.Template
	vnet           string
	nsg            string
	subnet         string
	volumeSize     int32
	imageOffer     string
	imagePublisher string
	imageSKU       string
	imageVersion   string
}

// Note requires the following environment variables:
//   - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
//   - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
//   - `AZURE_CLIENT_SECRET`: Specifies the app secret to use.

func (p *provider) getClient() (compute.VirtualMachinesClient, error) {
	vmClient := compute.NewVirtualMachinesClient(p.subscriptionId)
	// pulled from [https://github.com/Azure/azure-sdk-for-go#more-authentication-details]
	// Uses environment variables for the authorization client_secret/client_id
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}
	
	vmClient.Authorizer = auth
	return vmClient, nil
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	p.vmSize = compute.VirtualMachineSizeTypesBasicA0
	for _, opt := range opts {
		opt(p)
	}
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	if p.volumeSize == 0 {
		p.volumeSize = 1032
	}
	
	return p
}
