// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
)

func TestSetupFirewall(t *testing.T) {
	defer gock.Off()

	gock.New("https://www.googleapis.com").
		Get("/compute/v1/projects/my-project/global/firewalls/default-allow-docker").
		Reply(404)

	gock.New("https://www.googleapis.com").
		Post("/compute/v1/projects/my-project/global/firewalls").
		Reply(200).
		BodyString(`{ "name": "operation-name" }`)

	gock.New("https://www.googleapis.com").
		Get("/compute/v1/projects/my-project/global/operations/operation-name").
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	p, err := New(
		WithClient(http.DefaultClient),
		WithZone("us-central1-a"),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = p.(*provider).setupFirewall(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

func TestSetupFirewall_Exists(t *testing.T) {
	defer gock.Off()

	gock.New("https://www.googleapis.com").
		Get("/compute/v1/projects/my-project/global/firewalls/default-allow-docker").
		Reply(200).
		BodyString(findFirewallRes)

	p, err := New(
		WithClient(http.DefaultClient),
		WithZone("us-central1-a"),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = p.(*provider).setupFirewall(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

var createFirewallReq = `
{
  "name": "default-allow-docker",
  "selfLink": "projects/drone-1191/global/firewalls/default-allow-docker",
  "network": "projects/drone-1191/global/networks/default",
  "direction": "INGRESS",
  "priority": 1000,
  "targetTags": [
    "allow-docker"
  ],
  "allowed": [
    {
      "IPProtocol": "tcp",
      "ports": [
        "2376"
      ]
    }
  ],
  "sourceRanges": [
    "0.0.0.0/0"
  ]
}
`

var findFirewallRes = `
{
  "allowed": [
    {
      "IPProtocol": "tcp",
      "ports": [
        "2376"
      ]
    }
  ],
  "creationTimestamp": "2018-03-10T11:31:09.445-08:00",
  "description": "",
  "direction": "INGRESS",
  "id": "3206167972979853122",
  "kind": "compute#firewall",
  "name": "default-allow-docker",
  "network": "projects/my-project/global/networks/default",
  "priority": 1000,
  "selfLink": "projects/my-project/global/firewalls/default-allow-docker",
  "sourceRanges": [
    "0.0.0.0/0"
  ],
  "targetTags": [
    "allow-docker"
  ]
}
`

/*

POST https://www.googleapis.com/compute/v1/projects/drone-1191/global/firewalls
{
  "name": "default-allow-docker",
  "selfLink": "projects/drone-1191/global/firewalls/default-allow-docker",
  "network": "projects/drone-1191/global/networks/default",
  "direction": "INGRESS",
  "priority": 1000,
  "targetTags": [
    "allow-docker"
  ],
  "allowed": [
    {
      "IPProtocol": "tcp",
      "ports": [
        "2376"
      ]
    }
  ],
  "sourceRanges": [
    "0.0.0.0/0"
  ]
}

*/
