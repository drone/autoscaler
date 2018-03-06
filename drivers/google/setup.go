// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (p *provider) setup(ctx context.Context) error {
	var g errgroup.Group
	if p.key == "" {
		g.Go(func() error {
			return p.setupKeypair(ctx)
		})
	}
	return g.Wait()
}

func (p *provider) setupKeypair(ctx context.Context) error {
	return nil // TODO
}
