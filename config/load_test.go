// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestDefaults(t *testing.T) {
	conf := MustLoad()
	if got, want := conf.Logs.Debug, true; got != want {
		t.Errorf("Want default DRONE_LOGS_DEBUG of %v, got %v", want, got)
	}
	if got, want := conf.Interval, time.Minute*5; got != want {
		t.Errorf("Want default DRONE_INTERVAL of %s, got %s", want, got)
	}
	if got, want := conf.Pool.Max, 4; got != want {
		t.Errorf("Want default DRONE_POOL_MIN of %d, got %d", want, got)
	}
	if got, want := conf.Pool.Min, 2; got != want {
		t.Errorf("Want default DRONE_POOL_MAX of %d, got %d", want, got)
	}
	if got, want := conf.Pool.MinAge, time.Minute*60; got != want {
		t.Errorf("Want default DRONE_POOL_MIN_AGE of %d, got %d", want, got)
	}
	if got, want := conf.HTTP.Port, ":8080"; got != want {
		t.Errorf("Want default DRONE_HTTP_PORT of %s, got %s", want, got)
	}
	if got, want := conf.Database.Driver, "sqlite3"; got != want {
		t.Errorf("Want default DRONE_DATABASE_DRIVER of %s, got %s", want, got)
	}
	if got, want := conf.Database.Datasource, "database.sqlite?cache=shared&mode=rwc&_busy_timeout=9999999"; got != want {
		t.Errorf("Want default DRONE_DATABASE_DATASOURCE of %s, got %s", want, got)
	}
	//
	// Digital Ocean
	if got, want := conf.DigitalOcean.Image, "docker-16-04"; got != want {
		t.Errorf("Want default DRONE_DIGITALOCEAN_IMAGE of %s, got %s", want, got)
	}
	if got, want := conf.DigitalOcean.Region, "nyc3"; got != want {
		t.Errorf("Want default DRONE_DIGITALOCEAN_REGION of %s, got %s", want, got)
	}
	if got, want := conf.DigitalOcean.Size, "s-1vcpu-3gb"; got != want {
		t.Errorf("Want default DRONE_DIGITALOCEAN_SIZE of %s, got %s", want, got)
	}
	//
	//
	// Google
	//
	if got, want := conf.Google.Zone, "us-central1-a"; got != want {
		t.Errorf("Want default DRONE_GOOGLE_ZONE of %s, got %s", want, got)
	}
	if got, want := conf.Google.MachineType, "n1-standard-1"; got != want {
		t.Errorf("Want default DRONE_GOOGLE_MACHINE_TYPE of %s, got %s", want, got)
	}
	if got, want := conf.Google.MachineImage, "ubuntu-1510-wily-v20151114"; got != want {
		t.Errorf("Want default DRONE_GOOGLE_MACHINE_IMAGE of %s, got %s", want, got)
	}
	if got, want := conf.Google.DiskType, "pd-standard"; got != want {
		t.Errorf("Want default DRONE_GOOGLE_DISK_TYPE of %s, got %s", want, got)
	}
	//
	// Hetzner Cloud
	if got, want := conf.HetznerCloud.Image, "ubuntu-16.04"; got != want {
		t.Errorf("Want default DRONE_HETZNERCLOUD_IMAGE of %s, got %s", want, got)
	}
	if got, want := conf.HetznerCloud.Datacenter, "nbg1-dc3"; got != want {
		t.Errorf("Want default DRONE_HETZNERCLOUD_DATACENTER of %s, got %s", want, got)
	}
	if got, want := conf.HetznerCloud.ServerType, "cx11"; got != want {
		t.Errorf("Want default DRONE_HETZNERCLOUD_SERVER_TYPE of %s, got %s", want, got)
	}
}

func TestLoad(t *testing.T) {
	environ := map[string]string{
		"DRONE_INTERVAL":                 "1m",
		"DRONE_SLACK_WEBHOOK":            "https://hooks.slack.com/services/XXX/YYY/ZZZ",
		"DRONE_LOGS_DEBUG":               "true",
		"DRONE_LOGS_COLOR":               "true",
		"DRONE_LOGS_PRETTY":              "true",
		"DRONE_POOL_MIN_AGE":             "1h",
		"DRONE_POOL_MIN":                 "1",
		"DRONE_POOL_MAX":                 "5",
		"DRONE_SERVER_HOST":              "drone.company.com",
		"DRONE_SERVER_PROTO":             "http",
		"DRONE_SERVER_TOKEN":             "633eb230f5",
		"DRONE_HTTP_HOST":                "autoscaler.drone.company.com",
		"DRONE_HTTP_PORT":                "633eb230f5",
		"DRONE_AGENT_HOST":               "drone.company.com:9000",
		"DRONE_AGENT_TOKEN":              "f5064039f5",
		"DRONE_AGENT_IMAGE":              "drone/agent:0.8",
		"DRONE_AGENT_CONCURRENCY":        "2",
		"DRONE_TLS_AUTOCERT":             "true",
		"DRONE_TLS_CERT":                 "/path/to/cert.crt",
		"DRONE_TLS_KEY":                  "/path/to/cert.key",
		"DRONE_PROMETHEUS_TOKEN":         "b359e05e8",
		"DRONE_DATABASE_DRIVER":          "mysql",
		"DRONE_DATABASE_DATASOURCE":      "user:password@/dbname",
		"DRONE_DIGITALOCEAN_TOKEN":       "2573633eb",
		"DRONE_DIGITALOCEAN_IMAGE":       "docker-16-04",
		"DRONE_DIGITALOCEAN_REGION":      "ncy1",
		"DRONE_DIGITALOCEAN_SSHKEY":      "/path/to/ssh/key",
		"DRONE_DIGITALOCEAN_SIZE":        "s-1vcpu-1gb",
		"DRONE_DIGITALOCEAN_IPV6":        "true",
		"DRONE_DIGITALOCEAN_TAGS":        "drone,agent,prod",
		"DRONE_GOOGLE_ZONE":              "us-central1-b",
		"DRONE_GOOGLE_MACHINE_TYPE":      "f1-micro",
		"DRONE_GOOGLE_MACHINE_IMAGE":     "ubuntu-1510-wily-v20151114",
		"DRONE_GOOGLE_DISK_TYPE":         "pd-standard",
		"DRONE_GOOGLE_NETWORK":           "default",
		"DRONE_GOOGLE_SUBNETWORK":        "",
		"DRONE_GOOGLE_PREEMPTIBLE":       "true",
		"DRONE_GOOGLE_SCOPES":            "devstorage.read_only",
		"DRONE_GOOGLE_DISK_SIZE":         "10",
		"DRONE_GOOGLE_PROJECT":           "project-foo",
		"DRONE_GOOGLE_TAGS":              "drone,agent,prod",
		"DRONE_AMAZON_INSTANCE":          "t2.medium",
		"DRONE_AMAZON_REGION":            "us-east-2",
		"DRONE_AMAZON_SSHKEY":            "/path/to/ssh/key",
		"DRONE_AMAZON_SSHKEY_NAME":       "id_rsa",
		"DRONE_AMAZON_SUBNET_ID":         "subnet-0b32177f",
		"DRONE_AMAZON_SECURITY_GROUP":    "sg-770eabe1",
		"DRONE_HETZNERCLOUD_TOKEN":       "12345678",
		"DRONE_HETZNERCLOUD_IMAGE":       "ubuntu-16.04",
		"DRONE_HETZNERCLOUD_DATACENTER":  "nbg1-dc3",
		"DRONE_HETZNERCLOUD_SSHKEY":      "/path/to/ssh/key",
		"DRONE_HETZNERCLOUD_TYPE":        "cx11",
		"DRONE_HETZNERCLOUD_SSHKEY_ID":   "12345",
	}

	defer func() {
		// reset the environment.
		for k := range environ {
			os.Unsetenv(k)
		}
	}()

	// set test environment variables
	for k, v := range environ {
		os.Setenv(k, v)
	}

	a := MustLoad()
	b := Config{}
	json.Unmarshal(jsonConfig, &b)

	if !reflect.DeepEqual(a, b) {
		t.Errorf("configuration mismatch")
		pretty.Ldiff(t, a, b)
	}
}

var jsonConfig = []byte(`{
  "Interval": 60000000000,
  "Slack": {
    "Webhook": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
  },
  "Logs": {
    "Color": true,
    "Debug": true,
    "Pretty": true
  },
  "Pool": {
    "Min": 1,
    "Max": 5,
    "MinAge": 3600000000000
  },
  "Server": {
    "Host": "drone.company.com",
    "Proto": "http",
    "Token": "633eb230f5"
  },
  "Agent": {
    "Host": "drone.company.com:9000",
    "Token": "f5064039f5",
    "Image": "drone/agent:0.8",
    "Concurrency": 2
  },
  "HTTP": {
    "Host": "autoscaler.drone.company.com",
    "Port": "633eb230f5"
  },
  "TLS": {
    "Autocert": true,
    "Cert": "/path/to/cert.crt",
    "Key": "/path/to/cert.key"
  },
  "Prometheus": {
    "Token": "b359e05e8"
  },
  "Database": {
    "Driver": "mysql",
    "Datasource": "user:password@/dbname"
  },
  "DigitalOcean": {
    "Token": "2573633eb",
    "Image": "docker-16-04",
    "Region": "ncy1",
    "SSHKey": "/path/to/ssh/key",
    "Size": "s-1vcpu-1gb",
    "IPv6": true,
    "Tags": [
      "drone",
      "agent",
      "prod"
    ]
	},
	"Amazon": {
    "Instance": "t2.medium",
		"Region": "us-east-2",
		"SSHKey": "/path/to/ssh/key",
		"SSHKeyName": "id_rsa",
		"SubnetID": "subnet-0b32177f",
		"SecurityGroup": [ "sg-770eabe1" ]
	},
  "Google": {
    "Zone": "us-central1-b",
    "MachineType": "f1-micro",
    "MachineImage": "ubuntu-1510-wily-v20151114",
    "DiskType": "pd-standard",
    "Address": "",
    "Network": "default",
    "Subnetwork": "",
    "Preemptible": true,
    "Scopes": "devstorage.read_only",
    "DiskSize": 10,
    "Project": "project-foo",
    "Tags": [
      "drone",
      "agent",
      "prod"
    ]
  },
  "HetznerCloud": {
    "Token": "12345678",
    "Image": "ubuntu-16.04",
    "Datacenter": "nbg1-dc3",
    "SSHKey": "/path/to/ssh/key",
    "ServerType": "cx11",
    "SSHKeyID": "12345"
  }
}`)
