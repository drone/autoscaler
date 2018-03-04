// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package certs

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Error(err)
	}

	_, err = GenerateCert("company.com", ca)
	if err != nil {
		t.Error(err)
	}
}
