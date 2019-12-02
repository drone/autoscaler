// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package packet

import (
	"context"
	"errors"

	"github.com/drone/autoscaler/logger"

	"github.com/packethost/packngo"
	"golang.org/x/sync/errgroup"
)

func (p *provider) setup(ctx context.Context) error {
	var g errgroup.Group
	if p.sshkey == "" {
		g.Go(func() error {
			return p.setupKeypair(ctx)
		})
	}
	return g.Wait()
}

// helper funciton to ascertain the ID of an existing SSH
// key to use when provisioning instances. This is only
// necessary when the user has not provided the ID.
func (p *provider) setupKeypair(ctx context.Context) error {
	logger := logger.FromContext(ctx)
	logger.Debugln("finding default ssh key")

	keys, _, err := p.client.SSHKeys.List()
	if err != nil {
		return err
	}

	index := map[string]packngo.SSHKey{}
	for _, key := range keys {
		index[key.Label] = key
	}

	// if the account has multiple keys configured we will
	// attempt to use an existing key based on naming convention.
	for _, name := range []string{"drone", "id_rsa_drone"} {
		key, ok := index[name]
		if !ok {
			continue
		}
		p.sshkey = key.Key

		logger.
			WithField("id", key.ID).
			WithField("label", key.Key).
			WithField("fingerprint", key.FingerPrint).
			Debugln("using default ssh key")
		return nil
	}

	// if there were no matches but the account has at least
	// one keypair already created we will select the first
	// in the list.
	if len(keys) > 0 {
		key := keys[0]
		p.sshkey = key.ID

		logger.
			WithField("id", key.ID).
			WithField("label", key.Label).
			WithField("fingerprint", key.FingerPrint).
			Debugln("using default ssh key")
		return nil
	}

	return errors.New("No matching keys")
}
