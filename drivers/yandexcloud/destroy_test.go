package yandexcloud_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"

	"github.com/drone/autoscaler"
)

type providerDestroySuite struct {
	baseProviderSuite
}

func (s *providerDestroySuite) SetupTest() {
	request := &compute.CreateInstanceRequest{
		FolderId:   s.folderID,
		Name:       "drone-integration-destroy-test",
		ZoneId:     s.zoneID,
		PlatformId: "standard-v3",
		ResourcesSpec: &compute.ResourcesSpec{
			Cores:  2,
			Memory: 2 * 1024 * 1024 * 1024,
		},
		BootDiskSpec: &compute.AttachedDiskSpec{
			AutoDelete: true,
			Disk: &compute.AttachedDiskSpec_DiskSpec_{
				DiskSpec: &compute.AttachedDiskSpec_DiskSpec{
					TypeId: "network-hdd",
					Size:   10 * 1024 * 1024 * 1024,
					Source: &compute.AttachedDiskSpec_DiskSpec_ImageId{
						ImageId: s.imageID,
					},
				},
			},
		},
		NetworkInterfaceSpecs: []*compute.NetworkInterfaceSpec{
			{
				SubnetId: s.subnetID,
				PrimaryV4AddressSpec: &compute.PrimaryAddressSpec{
					OneToOneNatSpec: &compute.OneToOneNatSpec{
						IpVersion: compute.IpVersion_IPV4,
					},
				},
			},
		},
	}

	op, err := s.sdk.WrapOperation(s.sdk.Compute().Instance().Create(noContext, request))
	s.NoError(err)

	err = op.Wait(noContext)
	s.NoError(err)

	resp, err := op.Response()
	s.NoError(err)

	ycInstance := resp.(*compute.Instance)
	s.instanceID = ycInstance.Id
}

func (s *providerDestroySuite) TestDestroy() {
	err := s.provider.Destroy(noContext, &autoscaler.Instance{ID: s.instanceID})
	s.NoError(err)
}

func Test_provider_DestroyIntegration(t *testing.T) {
	suite.Run(t, new(providerDestroySuite))
}
