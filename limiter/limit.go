// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package limiter

import (
	"encoding/pem"

	"github.com/drone/autoscaler"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ed25519"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
GB/hFnXEg63vDZ2W6mKFhLxZTuxMrlN/C/0iVZ2LfPQ=
-----END PUBLIC KEY-----
`)

// Limit wraps the ServerStore to limit server creation
// within the limitions of the license.
func Limit(server autoscaler.ServerStore, token string) autoscaler.ServerStore {
	if token == "" {
		// if the token is empty the software is being
		// used without a license. We assume this is for
		// trial purposes and grant limited trial access.
		return &limiter{server, &License{
			Lim: 10,
		}}
	}
	block, _ := pem.Decode(publicKey)
	license, err := ParseVerify(token, ed25519.PublicKey(block.Bytes))
	if err != nil {
		panic(err)
	}
	log.Info().
		Str("key", license.Key).
		Str("pro", license.Pro).
		Str("sub", license.Sub).
		Int("lim", license.Lim).
		Time("iat", license.Iss).
		Time("exp", license.Exp).
		Msg("license verified")
	return &limiter{server, license}
}
