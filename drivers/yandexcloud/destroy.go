package yandexcloud

import (
	"context"
	"fmt"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"

	"github.com/drone/autoscaler"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	op, err := p.service.WrapOperation(p.deleteInstance(ctx, instance.ID))
	if err != nil {
		return fmt.Errorf("make delete operation: %w", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("wait delete operation: %w", err)
	}

	return nil
}

func (p *provider) deleteInstance(ctx context.Context, id string) (*operation.Operation, error) {
	return p.service.Compute().Instance().Delete(ctx, &compute.DeleteInstanceRequest{
		InstanceId: id,
	})
}
