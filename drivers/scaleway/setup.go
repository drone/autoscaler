// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"
)

func (p *provider) setup(ctx context.Context) error {
	var g errgroup.Group
	if p.client == nil {
		g.Go(func() error {
			return p.newClient(ctx)
		})
	}
	return g.Wait()
}

func (p *provider) newClient(ctx context.Context) error {
	client, err := scw.NewClient(
		scw.WithDefaultProjectID(p.orgID),
		scw.WithAuth(p.accessKey, p.secretKey),
	)
	if err != nil {
		return err
	}

	p.client = client
	return nil
}
