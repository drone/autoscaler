// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"errors"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

func (p *provider) setup(ctx context.Context) error {
	logger := log.Ctx(ctx)

	client := newClient(ctx, p.token)
	keys, _, err := client.Keys.List(ctx, &godo.ListOptions{})
	if err != nil {
		return err
	}

	// if the account has a single keypair created we
	// will use this as the default.
	if len(keys) == 1 {
		key := keys[0]
		p.key = key.Fingerprint

		logger.Debug().
			Str("name", key.Name).
			Str("fingerprint", key.Fingerprint).
			Msg("using default ssh key")
		return nil
	}

	index := map[string]string{}
	for _, key := range keys {
		index[key.Name] = key.Fingerprint
	}

	// if the account has multiple keys configured we will
	// attempt to use an existing key based on naming convention.
	for _, name := range []string{"id_rsa_drone", "id_rsa"} {
		fingerprint, ok := index[name]
		if !ok {
			continue
		}

		p.key = fingerprint
		logger.Debug().
			Str("name", name).
			Str("fingerprint", fingerprint).
			Msg("using default ssh key")
	}

	return errors.New("No matching keys")
}
