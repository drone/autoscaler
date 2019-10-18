// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package amazon

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func (p *provider) setup(ctx context.Context) error {
	var g errgroup.Group
	if p.key == "" {
		g.Go(func() error {
			return p.setupKeypair(ctx)
		})
	}
	if p.subnet == "" {
		// TODO: find or create subnet
	}
	if len(p.groups) == 0 {
		// TODO: find or create security groups
	}
	return g.Wait()
}

func (p *provider) setupKeypair(ctx context.Context) error {
	logger := log.Ctx(ctx)

	logger.Debug().
		Msg("finding default ssh key")

	opts := new(ec2.DescribeKeyPairsInput)
	keys, err := p.getClient().DescribeKeyPairs(opts)
	if err != nil {
		return err
	}

	index := map[string]string{}
	for _, key := range keys.KeyPairs {
		index[*key.KeyName] = *key.KeyFingerprint
	}

	// if the account has multiple keys configured we will
	// attempt to use an existing key based on naming convention.
	for _, name := range []string{"drone", "id_rsa_drone"} {
		fingerprint, ok := index[name]
		if !ok {
			continue
		}
		p.key = name

		logger.Debug().
			Str("name", name).
			Str("fingerprint", fingerprint).
			Msg("using default ssh key")
		return nil
	}

	// if there were no matches but the account has at least
	// one keypair already created we will select the first
	// in the list.
	if len(keys.KeyPairs) > 0 {
		key := keys.KeyPairs[0]
		p.key = *key.KeyName

		logger.Debug().
			Str("name", *key.KeyName).
			Str("fingerprint", *key.KeyFingerprint).
			Msg("using default ssh key")
		return nil
	}

	return errors.New("No matching keys")
}
