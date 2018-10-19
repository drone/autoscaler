// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"context"
	"errors"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func (p *provider) setup(ctx context.Context) error {
	var g errgroup.Group
	if p.key == "" {
		g.Go(func() error {
			return p.findKeyPair(ctx)
		})
	}

	return g.Wait()
}

func (p *provider) findKeyPair(ctx context.Context) error {
	logger := log.Ctx(ctx)

	logger.Debug().
		Msg("finding default ssh key")

	allPages, err := keypairs.List(p.computeClient).AllPages()
	if err != nil {
		return err
	}
	keys, err := keypairs.ExtractKeyPairs(allPages)

	if err != nil {
		return err
	}

	index := map[string]keypairs.KeyPair{}
	for _, key := range keys {
		index[key.Name] = key
	}

	// if the account has multiple keys configured we will
	// attempt to use an existing key based on naming convention.
	for _, name := range []string{"drone", "id_rsa_drone"} {
		key, ok := index[name]
		if !ok {
			continue
		}
		p.key = key.Name

		logger.Debug().
			Str("name", name).
			Str("fingerprint", key.Fingerprint).
			Msg("using default ssh key")
		return nil
	}

	// if there were no matches but the account has at least
	// one keypair already created we will select the first
	// in the list.
	if len(keys) > 0 {
		key := keys[0]
		p.key = key.Name

		logger.Debug().
			Str("name", key.Name).
			Str("fingerprint", key.Fingerprint).
			Msg("using default ssh key")
		return nil
	}
	return errors.New("no matching keys")
}
