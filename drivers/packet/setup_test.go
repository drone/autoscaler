// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package packet

import (
	"context"
	"os"
	"testing"

	"github.com/h2non/gock"
)

const (
	baseURL      = "https://api.packet.net/"
	getDevice    = "/devices"
	getSSH       = "/ssh-keys"
	projectID    = "x"
	createDevice = "/projects/" + projectID + getDevice
	instanceID   = "92b0facf-189e-4bbf-81a8-bc56c0c4dc88"
	apiKey       = "apiKey"
	sshKey       = "sshKey"
	hostname     = "hostname"
	tag          = "tag"
)

var (
	prov               *provider
	respCreate         string
	respCreateInactive string
	respSSHKeys        string
)

func TestMain(m *testing.M) {
	prov = New(
		WithProject(projectID),
		WithTags(tag),
		WithHostname(hostname),
		WithAPIKey(apiKey),
	).(*provider)

	respCreate = `
{
  "id": "` + instanceID + `",
  "state": "active",
  "tags": ["` + tag + `"],
  "hostname": "` + hostname + `",
  "operating_system": {
    "slug": "` + prov.os + `"
  },
  "facility": {
    "code": "ewr1"
  },
  "ip_addresses": [
    {
      "address_family": 4,
      "public": true,
      "address": "147.75.77.155"
    }
  ],
  "plan": {
    "slug": "baremetal_0"
  }
}
`
	respCreateInactive = `
{
  "id": "` + instanceID + `",
  "state": "inactive",
  "tags": ["` + tag + `"],
  "hostname": "` + hostname + `",
  "operating_system": {
    "slug": "` + prov.os + `"
  },
  "facility": {
    "code": "ewr1"
  },
  "ip_addresses": [
    {
      "address_family": 4,
      "public": true,
      "address": "147.75.77.155"
    }
  ],
  "plan": {
    "slug": "baremetal_0"
  }
}
`

	respSSHKeys = `
{
  "ssh_keys": [
    {
      "id": "` + sshKey + `",
      "label": "label",
      "key": "key",
      "fingerprint": "fingerprint"
    }
  ]
}
`
	os.Exit(m.Run())

}

func TestSetup_Keypair(t *testing.T) {
	defer gock.Off()

	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Get(getSSH).
		Reply(200).
		BodyString(respSSHKeys)

	if err := prov.setupKeypair(context.Background()); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if prov.sshkey != sshKey {
		t.Errorf("expected: %s, got: %s", sshKey, prov.sshkey)
	}
}
