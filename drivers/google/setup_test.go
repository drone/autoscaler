// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	compute "google.golang.org/api/compute/v1"
)

func TestSetupFirewall(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/firewalls/default-allow-docker").
		Reply(404)

	gock.New("https://compute.googleapis.com").
		Post("/compute/v1/projects/my-project/global/firewalls").
		JSON(createFirewallMock).
		Reply(200).
		BodyString(`{ "name": "operation-name" }`)

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/operations/operation-name").
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	p, err := New(
		WithClient(http.DefaultClient),
		WithZones("us-central1-a"),
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

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/firewalls/default-allow-docker").
		Reply(200).
		BodyString(findFirewallRes)

	p, err := New(
		WithClient(http.DefaultClient),
		WithZones("us-central1-a"),
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

var createFirewallMock = &compute.Firewall{
	Allowed: []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports:      []string{"2376"},
		},
	},
	Direction:    "INGRESS",
	Name:         "default-allow-docker",
	Network:      "global/networks/default",
	Priority:     1000,
	SourceRanges: []string{"0.0.0.0/0"},
	TargetTags:   []string{"allow-docker"},
}

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
