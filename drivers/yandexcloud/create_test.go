package yandexcloud_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"

	"github.com/drone/autoscaler"
)

type providerCreteSuite struct {
	baseProviderSuite
}

func (s *providerCreteSuite) deleteInstance(instanceID string) {
	op, err := s.sdk.WrapOperation(s.sdk.Compute().Instance().Delete(noContext, &compute.DeleteInstanceRequest{
		InstanceId: instanceID,
	}))
	s.NoError(err)

	err = op.Wait(noContext)
	s.NoError(err)
}

func (s *providerCreteSuite) TestCreate() {
	var instanceID string

	defer func() {
		s.deleteInstance(instanceID)
	}()

	instance, err := s.provider.Create(context.Background(), autoscaler.InstanceCreateOpts{
		Name: "drone-runner-test",
	})
	s.NoError(err)

	s.NotEqual(nil, instance)
	instanceID = instance.ID
}

func Test_provider_CreteIntegration(t *testing.T) {
	suite.Run(t, new(providerCreteSuite))
}
