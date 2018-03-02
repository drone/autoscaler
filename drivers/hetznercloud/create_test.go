// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/mocks"

	"github.com/golang/mock/gomock"
	"github.com/h2non/gock"
	"golang.org/x/crypto/ssh"
)

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Post("/v1/servers").
		Reply(200).
		BodyString(respInstanceCreate)

	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}

	// base provider to mock SSH calls.
	mockProvider := mocks.NewMockProvider(controller)
	mockProvider.EXPECT().Ping(gomock.Any(), gomock.Any()).Return(nil)
	mockProvider.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	p := Provider{
		Provider: mockProvider,
		config:   mockConfig,
		signer:   mockSigner,
	}

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

	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}

	p := Provider{
		Provider: nil,
		config:   mockConfig,
		signer:   mockSigner,
	}

	_, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expect error returned from hetzner cloud")
	}
}

func TestCreate_PingTimeout(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Post("/v1/servers").
		Reply(200).
		BodyString(respInstanceCreate)

	mockError := errors.New("oh no")
	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// base provider to mock SSH calls.
	mockProvider := mocks.NewMockProvider(controller)
	mockProvider.EXPECT().Ping(ctx, gomock.Any()).Return(mockError)

	p := Provider{
		Provider: mockProvider,
		config:   mockConfig,
		signer:   mockSigner,
	}

	instance, err := p.Create(ctx, autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expected context deadline exceeded, got nil")
	} else if err.Error() != "context deadline exceeded" {
		t.Errorf("Expected context deadline exceeded, got %s", err)
	}

	t.Run("Attributes", testInstance(instance))
	t.Run("Address", testInstanceAddress(instance))
}

func TestCreate_ExecError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Post("/v1/servers").
		Reply(200).
		BodyString(respInstanceCreate)

	mockContext := context.Background()
	mockLogs := []byte("-bash: curl: command not found")
	mockError := errors.New("uh oh")
	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}

	// base provider to mock SSH calls.
	mockProvider := mocks.NewMockProvider(controller)
	mockProvider.EXPECT().Ping(mockContext, gomock.Any()).Return(nil)
	mockProvider.EXPECT().Execute(mockContext, gomock.Any(), gomock.Any()).Return(mockLogs, mockError)

	p := Provider{
		Provider: mockProvider,
		config:   mockConfig,
		signer:   mockSigner,
	}

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if lerr, ok := err.(*autoscaler.InstanceError); !ok {
		t.Errorf("Want InstanceError")
	} else if err == nil {
		t.Errorf("Want InstanceError got nil")
	} else if lerr.Err != mockError {
		t.Errorf("Want InstanceError to wrap the ssh error")
	} else if !bytes.Equal(lerr.Logs, mockLogs) {
		t.Errorf("Want InstanceError to include the logs")
	}

	t.Run("Attributes", testInstance(instance))
	t.Run("Address", testInstanceAddress(instance))
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, "544037"; got != want {
			t.Errorf("Want instance ID %v, got %v", want, got)
		}
		if got, want := instance.Image, "ubuntu-16.04"; got != want {
			t.Errorf("Want instance Image %v, got %v", want, got)
		}
		if got, want := instance.Name, "test"; got != want {
			t.Errorf("Want instance Name %v, got %v", want, got)
		}
		if got, want := instance.Region, "nbg1-dc3"; got != want {
			t.Errorf("Want instance Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderHetznerCloud; got != want {
			t.Errorf("Want instance Provider %v, got %v", want, got)
		}
		if instance.Secret == "" {
			t.Errorf("Want instance secret populated, got empty")
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
      "name": "ubuntu-16.04",
      "description": "Ubuntu 16.04",
      "image_size": null,
      "disk_size": 5,
      "created": "2018-01-15T11:34:45+00:00",
      "created_from": null,
      "bound_to": null,
      "os_flavor": "ubuntu",
      "os_version": "16.04",
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
