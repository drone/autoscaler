package yandexcloud_test

import (
	"context"
	"os"

	"github.com/stretchr/testify/suite"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/yandexcloud"
)

var noContext = context.TODO()

type baseProviderSuite struct {
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

func (s *baseProviderSuite) SetupSuite() {
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

	s.sdk, err = ycsdk.Build(noContext, ycsdk.Config{
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
		Family:   "container-optimized-image",
	})
	s.NoError(err)
	s.imageID = image.Id
}
