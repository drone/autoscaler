// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"strings"

	"github.com/drone/autoscaler"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	op, err := p.service.Instances.Delete(p.project, p.zone, instance.ID).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Error 404") {
			// the instance doesn't exist. which is okay, that's our goal anyway
			return nil
		}
		return err
	}
	return p.waitZoneOperation(ctx, op.Name)
}
