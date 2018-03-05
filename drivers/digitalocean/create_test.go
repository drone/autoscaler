// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/drone/autoscaler"
	"github.com/golang/mock/gomock"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Post("/v2/droplets").
		Reply(200).
		BodyString(respDropletCreate)

	gock.New("https://api.digitalocean.com").
		Get("/v2/droplets/3164494").
		Reply(200).
		BodyString(respDropletDesc)

	p := New(
		WithFingerprint("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	)

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err != nil {
		t.Error(err)
	}

	t.Run("Attributes", testInstance(instance))
}

func TestCreate_CreateError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Post("/v2/droplets").
		Reply(500)

	p := New(
		WithFingerprint("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	)

	_, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expect error returned from digital ocean")
	} else if _, ok := err.(*godo.ErrorResponse); !ok {
		t.Errorf("Expect ErrorResponse digital ocean")
	}
}

func TestCreate_DescribeError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Post("/v2/droplets").
		Reply(200).
		BodyString(respDropletCreate)

	gock.New("https://api.digitalocean.com").
		Get("/v2/droplets/3164494").
		Reply(500)

	p := New(
		WithFingerprint("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	)

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expect error returned from digital ocean")
	} else if _, ok := err.(*godo.ErrorResponse); !ok {
		t.Errorf("Expect ErrorResponse digital ocean")
	}

	t.Run("Attributes", testInstance(instance))
}

func TestCreate_DescribeTimeout(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Post("/v2/droplets").
		Reply(200).
		BodyString(respDropletCreate)

	gock.New("https://api.digitalocean.com").
		Get("/v2/droplets/3164494").
		Reply(200).
		BodyString(respDropletCreate) // no network data

	p := New(
		WithFingerprint("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	instance, err := p.Create(ctx, autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expected context deadline exceeded, got nil")
	} else if err.Error() != "context deadline exceeded" {
		t.Errorf("Expected context deadline exceeded, got %s", err)
	}

	t.Run("Attributes", testInstance(instance))
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, "3164494"; got != want {
			t.Errorf("Want droplet ID %v, got %v", want, got)
		}
		if got, want := instance.Image, "docker-16-04"; got != want {
			t.Errorf("Want droplet Image %v, got %v", want, got)
		}
		if got, want := instance.Name, "example.com"; got != want {
			t.Errorf("Want droplet Name %v, got %v", want, got)
		}
		if got, want := instance.Region, "sfo1"; got != want {
			t.Errorf("Want droplet Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderDigitalOcean; got != want {
			t.Errorf("Want droplet Provider %v, got %v", want, got)
		}
	}
}

func testInstanceAddress(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.Address, "104.131.186.241"; got != want {
			t.Errorf("Want droplet Address %v, got %v", want, got)
		}
	}
}

// sample response for POST /v2/droplets
const respDropletCreate = `
{
  "droplet": {
    "id": 3164494,
    "name": "example.com",
    "memory": 1024,
    "vcpus": 1,
    "disk": 25,
    "locked": true,
    "status": "new",
    "kernel": {
      "id": 2233,
      "name": "Ubuntu 14.04 x64 vmlinuz-3.13.0-37-generic",
      "version": "3.13.0-37-generic"
    },
    "created_at": "2014-11-14T16:36:31Z",
    "features": [
      "virtio"
    ],
    "backup_ids": [
      
    ],
    "snapshot_ids": [
      
    ],
    "image": {
      
    },
    "volume_ids": [
      
    ],
    "size": {
      
    },
    "size_slug": "s-1vcpu-1gb",
    "networks": {

    },
    "region": {
      
    },
    "tags": [
      "web"
    ]
  },
  "links": {
    "actions": [
      {
        "id": 36805096,
        "rel": "create",
        "href": "https:\/\/api.digitalocean.com\/v2\/actions\/36805096"
      }
    ]
  }
}
`

// sample response for POST /v2/droplets/:id
const respDropletDesc = `
{
  "droplet": {
    "id": 3164494,
    "name": "example.com",
    "memory": 1024,
    "vcpus": 1,
    "disk": 25,
    "locked": true,
    "status": "new",
    "kernel": {
      "id": 2233,
      "name": "Ubuntu 14.04 x64 vmlinuz-3.13.0-37-generic",
      "version": "3.13.0-37-generic"
    },
    "created_at": "2014-11-14T16:36:31Z",
    "features": [
      "virtio"
    ],
    "backup_ids": [
      
    ],
    "snapshot_ids": [
      
    ],
    "image": {
      
    },
    "volume_ids": [
      
    ],
    "size": {
      
    },
    "size_slug": "s-1vcpu-1gb",
    "networks": {
      "v4": [
        {
          "ip_address": "104.131.186.241",
          "netmask": "255.255.240.0",
          "gateway": "104.131.176.1",
          "type": "public"
        }
      ]
    },
    "region": {
      
    },
    "tags": [
      "web"
    ]
  },
  "links": {
    "actions": [
      {
        "id": 36805096,
        "rel": "create",
        "href": "https:\/\/api.digitalocean.com\/v2\/actions\/36805096"
      }
    ]
  }
}
`
