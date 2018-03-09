// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"errors"
	"time"

	"github.com/drone/autoscaler"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	op, err := p.service.Instances.Delete(p.proj, p.zone, instance.ID).Context(ctx).Do()
	if err != nil {
		return err
	}
	return p.waitZoneOperation(ctx, op.Name)
}

func (p *provider) waitZoneOperation(ctx context.Context, name string) error {
	for {
		op, err := p.service.ZoneOperations.Get(p.proj, p.zone, name).Context(ctx).Do()
		if err != nil {
			return err
		}
		if op.Error != nil {
			return errors.New(op.Error.Errors[0].Message)
		}
		if op.Status == "DONE" {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}
