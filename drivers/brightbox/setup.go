// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import (
	"context"

	"github.com/drone/autoscaler/logger"

	"github.com/brightbox/gobrightbox"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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
	logger := logger.FromContext(ctx)

	// Setup OAuth2 authentication
	conf := clientcredentials.Config{
		ClientID:     p.clientID,
		ClientSecret: p.clientSecret,
		Scopes:       []string{},
		TokenURL:     p.apiURL + "/token",
	}
	oc := conf.Client(oauth2.NoContext)

	// Setup API client
	client, err := gobrightbox.NewClient(p.apiURL, "", oc)
	if err != nil {
		logger.WithError(err).Errorln("unable to connect to Brightbox API")
		return err
	}

	p.client = client
	return nil
}
