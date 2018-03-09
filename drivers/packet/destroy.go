// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package packet

import (
	"context"

	"github.com/drone/autoscaler"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	_, err := p.client.Devices.Delete(instance.ID)
	return err
}
