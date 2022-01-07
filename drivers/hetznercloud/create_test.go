// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"testing"

	"github.com/drone/autoscaler"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Post("/v1/servers").
		Reply(200).
		BodyString(respInstanceCreate)

	p := New(
		WithToken("LRK9DAWQ1ZAEFSrCNEEzLCUwhYX1U3g7wMg4dTlkkDC96fyDuyJ39nVbVjCKSDfj"),
	).(*provider)
	p.init.Do(func() {}) // pre-initialize

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err != nil {
		t.Error(err)
	}

	t.Run("Attributes", testInstance(instance))
}

func TestCreate_CreateError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Post("/v1/servers").
		Reply(500)

	p := New(
		WithToken("LRK9DAWQ1ZAEFSrCNEEzLCUwhYX1U3g7wMg4dTlkkDC96fyDuyJ39nVbVjCKSDfj"),
	).(*provider)
	p.init.Do(func() {}) // pre-initialize

	_, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expect error returned from hetzner cloud")
	}
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, "544037"; got != want {
			t.Errorf("Want instance ID %v, got %v", want, got)
		}
		if got, want := instance.Image, "ubuntu-20.04"; got != want {
			t.Errorf("Want instance Image %v, got %v", want, got)
		}
		if got, want := instance.Name, "test"; got != want {
			t.Errorf("Want instance Name %v, got %v", want, got)
		}
		if got, want := instance.Region, "unknown"; got != want {
			t.Errorf("Want instance Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderHetznerCloud; got != want {
			t.Errorf("Want instance Provider %v, got %v", want, got)
		}
	}
}

func testInstanceAddress(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.Address, "195.201.93.137"; got != want {
			t.Errorf("Want instance Address %v, got %v", want, got)
		}
	}
}

// sample response for POST /v1/servers
const respInstanceCreate = `
{
  "server": {
    "id": 544037,
    "name": "test",
    "status": "initializing",
    "created": "2018-03-02T08:44:07+00:00",
    "public_net": {
      "ipv4": {
        "ip": "195.201.93.137",
        "blocked": false,
        "dns_ptr": "static.137.93.201.195.clients.your-server.de"
      },
      "ipv6": {
        "ip": "2a01:4f8:1c0c:6996::/64",
        "blocked": false,
        "dns_ptr": []
      },
      "floating_ips": []
    },
    "server_type": {
      "id": 1,
      "name": "cx11",
      "description": "CX11",
      "cores": 1,
      "memory": 2.0,
      "disk": 20,
      "prices": [
        {
          "location": "fsn1",
          "price_hourly": {
            "net": "0.0040000000",
            "gross": "0.0047600000000000"
          },
          "price_monthly": {
            "net": "2.4900000000",
            "gross": "2.9631000000000000"
          }
        },
        {
          "location": "nbg1",
          "price_hourly": {
            "net": "0.0040000000",
            "gross": "0.0047600000000000"
          },
          "price_monthly": {
            "net": "2.4900000000",
            "gross": "2.9631000000000000"
          }
        }
      ],
      "storage_type": "local"
    },
    "datacenter": {
      "id": 2,
      "name": "nbg1-dc3",
      "description": "Nuremberg 1 DC 3",
      "location": {
        "id": 2,
        "name": "nbg1",
        "description": "Nuremberg DC Park 1",
        "country": "DE",
        "city": "Nuremberg",
        "latitude": 49.452102,
        "longitude": 11.076665
      },
      "server_types": {
        "supported": [
          1,
          2,
          3,
          4,
          5,
          6,
          7,
          8,
          9,
          10
        ],
        "available": [
          1,
          2,
          3,
          4,
          5,
          6,
          7,
          8,
          9,
          10
        ]
      }
    },
    "image": {
      "id": 1,
      "type": "system",
      "status": "available",
      "name": "ubuntu-20.04",
      "description": "Ubuntu 20.04",
      "image_size": null,
      "disk_size": 5,
      "created": "2018-01-15T11:34:45+00:00",
      "created_from": null,
      "bound_to": null,
      "os_flavor": "ubuntu",
      "os_version": "20.04",
      "rapid_deploy": true
    },
    "iso": null,
    "rescue_enabled": false,
    "locked": false,
    "backup_window": null,
    "outgoing_traffic": 0,
    "ingoing_traffic": 0,
    "included_traffic": 21990232555520
  },
  "action": {
    "id": 279192,
    "command": "create_server",
    "status": "running",
    "progress": 0,
    "started": "2018-03-02T08:44:07+00:00",
    "finished": null,
    "resources": [
      {
        "id": 544037,
        "type": "server"
      }
    ],
    "error": null
  },
  "root_password": null
}
`
