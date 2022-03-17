package yandexcloud_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/yandexcloud"
)

var noContext = context.TODO()

type providerDestroySuite struct {
	suite.Suite

	token      string
	folderID   string
	subnetID   string
	imageID    string
	zoneID     string
	instanceID string

	sdk      *ycsdk.SDK
	provider autoscaler.Provider
}

func (s *providerDestroySuite) SetupSuite() {
	var err error

	token := os.Getenv("TEST_YANDEX_CLOUD_TOKEN")
	if token == "" {
		s.T().Skipf("Skip yandex cloud provider integration test. No token provided.")
	}
	folderID := os.Getenv("TEST_YANDEX_CLOUD_FOLDER_ID")
	if folderID == "" {
		s.T().Skipf("Skip yandex cloud provider integration test. No folder id provided.")
	}
	subnetID := os.Getenv("TEST_YANDEX_CLOUD_SUBNET_ID")
	if subnetID == "" {
		s.T().Skipf("Skip yandex cloud provider integration test. No subnet id provided.")
	}

	s.sdk, err = ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: ycsdk.OAuthToken(token),
	})
	s.NoError(err)

	s.provider, err = yandexcloud.New(
		yandexcloud.WithToken(token),
		yandexcloud.WithFolderID(folderID),
		yandexcloud.WithSubnetID(subnetID),
	)
	s.NoError(err)

	image, err := s.sdk.Compute().Image().GetLatestByFamily(noContext, &compute.GetImageLatestByFamilyRequest{
		FolderId: "standard-images",
		Family:   "debian-9",
	})
	s.NoError(err)
	s.imageID = image.Id
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
