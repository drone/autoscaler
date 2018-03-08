// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package limiter

import (
	"crypto"
	"encoding/json"
	"time"

	"github.com/o1egl/paseto"
)

// License represents a software license key.
type License struct {
	Key string    `json:"key"`
	Pro string    `json:"pro"`
	Sub string    `json:"sub"`
	Lim int       `json:"lim"`
	Iss time.Time `json:"iat"`
	Exp time.Time `json:"exp"`
}

// Expired returns true if the license is expired.
func (l *License) Expired() bool {
	return l.Exp.IsZero() == false && time.Now().After(l.Exp)
}

// ParseVerify parses and verifies the token, and returns
// a License from the token payload.
func ParseVerify(token string, publicKey crypto.PublicKey) (*License, error) {
	var payload []byte
	err := paseto.NewV2().Verify(token, publicKey, &payload, nil)
	if err != nil {
		return nil, err
	}
	out := new(License)
	err = json.Unmarshal(payload, out)
	return out, err
}
