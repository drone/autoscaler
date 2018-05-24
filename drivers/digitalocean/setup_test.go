// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"testing"

	"github.com/h2non/gock"
)

func TestSetupKey_Single(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Get("/v2/account/keys").
		Reply(200).
		BodyString(respSingleKey)

	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.key, "3b:16:bf:e4:8b:00:8b:b8:59:8c:a9:d3:f0:19:45:fa"; got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}
}

func TestSetupKey_FoundMatch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Get("/v2/account/keys").
		Reply(200).
		BodyString(respMultiKey)

	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.key, "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff"; got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}
}

func TestSetupKey_NoMatch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Get("/v2/account/keys").
		Reply(200).
		BodyString(respMultiKeyNoMatch)

	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.key, "3b:16:bf:e4:8b:00:8b:b8:59:8c:a9:d3:f0:19:45:fa"; got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}
}

var respSingleKey = `
{
  "ssh_keys": [
    {
      "id": 512189,
      "fingerprint": "3b:16:bf:e4:8b:00:8b:b8:59:8c:a9:d3:f0:19:45:fa",
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQDDHr/jh2Jy4yALcK4JyWbVkPRaWmhck3IgCoeOO3z1e2dBowLh64QAM+Qb72pxekALga2oi4GvT+TlWNhzPH4V example",
      "name": "My SSH Public Key"
    }
  ],
  "links": {
  },
  "meta": {
    "total": 1
  }
}
`

var respMultiKey = `
{
  "ssh_keys": [
    {
      "id": 512189,
      "fingerprint": "3b:16:bf:e4:8b:00:8b:b8:59:8c:a9:d3:f0:19:45:fa",
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQDDHr/jh2Jy4yALcK4JyWbVkPRaWmhck3IgCoeOO3z1e2dBowLh64QAM+Qb72pxekALga2oi4GvT+TlWNhzPH4V example",
      "name": "My SSH Public Key"
    },
    {
      "id": 513199,
      "fingerprint": "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff",
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQDDHr/jh2Jy4yALcK4JyWbVkPRaWmhck3IgCoeOO3z1e2dBowLh64QAM+Qb72pxekALga2oi4GvT+TlWNhzPH4V example",
      "name": "id_rsa_drone"
    }
  ],
  "links": {
  },
  "meta": {
    "total": 2
  }
}
`

var respMultiKeyNoMatch = `
{
  "ssh_keys": [
    {
      "id": 512189,
      "fingerprint": "3b:16:bf:e4:8b:00:8b:b8:59:8c:a9:d3:f0:19:45:fa",
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQDDHr/jh2Jy4yALcK4JyWbVkPRaWmhck3IgCoeOO3z1e2dBowLh64QAM+Qb72pxekALga2oi4GvT+TlWNhzPH4V example",
      "name": "My SSH Public Key"
    },
    {
      "id": 513199,
      "fingerprint": "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff",
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQDDHr/jh2Jy4yALcK4JyWbVkPRaWmhck3IgCoeOO3z1e2dBowLh64QAM+Qb72pxekALga2oi4GvT+TlWNhzPH4V example",
      "name": "My SSH Public Key2"
    }
  ],
  "links": {
  },
  "meta": {
    "total": 2
  }
}
`
