// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"

	"github.com/drone/autoscaler"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	return nil // TODO
}
