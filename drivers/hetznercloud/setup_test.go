// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"testing"

	"github.com/h2non/gock"
)

func TestSetupKey_ChooseFirst(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Get("/v1/ssh_keys").
		Reply(200).
		BodyString(respSingleKey)

	p := New(
		WithToken("LRK9DAWQ1ZAEFSrCNEEzLCUwhYX1U3g7wMg4dTlkkDC96fyDuyJ39nVbVjCKSDfj"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.key, 2323; got != want {
		t.Errorf("Want key id %d, got %d", want, got)
	}
}

func TestSetupKey_ChooseMatch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Get("/v1/ssh_keys").
		Reply(200).
		BodyString(respMultiKey)

	p := New(
		WithToken("LRK9DAWQ1ZAEFSrCNEEzLCUwhYX1U3g7wMg4dTlkkDC96fyDuyJ39nVbVjCKSDfj"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.key, 2324; got != want {
		t.Errorf("Want key id %d, got %d", want, got)
	}
}

const respSingleKey = `
{
  "ssh_keys": [
    {
      "id": 2323,
      "name": "My ssh key",
      "fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2f",
      "public_key": "ssh-rsa AAAjjk76kgf...Xt"
    }
  ]
}
`

const respMultiKey = `
{
  "ssh_keys": [
    {
      "id": 2323,
      "name": "My ssh key",
      "fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2f",
      "public_key": "ssh-rsa AAAjjk76kgf...Xt"
    },
    {
      "id": 2324,
      "name": "drone",
      "fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2f",
      "public_key": "ssh-rsa AAAjjk76kgf...Xt"
    }
  ]
}
`
