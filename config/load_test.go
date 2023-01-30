// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"io/ioutil"
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
	if got, want := conf.CapacityBuffer, 0; got != want {
		t.Errorf("Want default DRONE_CAPACITY_BUFFER of %d, got %d", want, got)
	}
	if got, want := conf.Pool.Max, 4; got != want {
		t.Errorf("Want default DRONE_POOL_MIN of %d, got %d", want, got)
	}
	if got, want := conf.Pool.Min, 2; got != want {
		t.Errorf("Want default DRONE_POOL_MAX of %d, got %d", want, got)
	}
	if got, want := conf.Pool.MinAge, time.Minute*55; got != want {
		t.Errorf("Want default DRONE_POOL_MIN_AGE of %d, got %d", want, got)
	}

	if got, want := conf.Check.Interval, time.Minute; got != want {
		t.Errorf("Want default DRONE_INSTALL_CHECK_INTERVAL of %s, got %s", want, got)
	}
	if got, want := conf.Check.Deadline, time.Minute*30; got != want {
		t.Errorf("Want default DRONE_INSTALL_CHECK_DEADLINE of %s, got %s", want, got)
	}

	if got, want := conf.HTTP.Port, ":8080"; got != want {
		t.Errorf("Want default DRONE_HTTP_PORT of %s, got %s", want, got)
	}
	if got, want := conf.HTTP.Root, "/"; got != want {
		t.Errorf("Want default DRONE_HTTP_ROOT of %s, got %s", want, got)
	}
	if got, want := conf.Database.Driver, "sqlite3"; got != want {
		t.Errorf("Want default DRONE_DATABASE_DRIVER of %s, got %s", want, got)
	}
	if got, want := conf.Database.Datasource, "database.sqlite?cache=shared&mode=rwc&_busy_timeout=9999999"; got != want {
		t.Errorf("Want default DRONE_DATABASE_DATASOURCE of %s, got %s", want, got)
	}
	if got, want := conf.Agent.Concurrency, 2; got != want {
		t.Errorf("Want default DRONE_AGENT_CONCURRENCY of %d, got %d", want, got)
	}
	if got, want := conf.Agent.Image, "drone/drone-runner-docker:1"; got != want {
		t.Errorf("Want default DRONE_AGENT_IMAGE of %s, got %s", want, got)
	}
}

func TestLoad(t *testing.T) {
	environ := map[string]string{
		"DRONE_INTERVAL":                   "1m",
		"DRONE_SLACK_WEBHOOK":              "https://hooks.slack.com/services/XXX/YYY/ZZZ",
		"DRONE_SLACK_CREATE":               "false",
		"DRONE_SLACK_DESTROY":              "false",
		"DRONE_LOGS_DEBUG":                 "true",
		"DRONE_LOGS_COLOR":                 "true",
		"DRONE_LOGS_PRETTY":                "true",
		"DRONE_CAPACITY_BUFFER":            "3",
		"DRONE_POOL_MIN_AGE":               "1h",
		"DRONE_POOL_MIN":                   "1",
		"DRONE_POOL_MAX":                   "5",
		"DRONE_SERVER_HOST":                "drone.company.com",
		"DRONE_SERVER_PROTO":               "http",
		"DRONE_SERVER_TOKEN":               "633eb230f5",
		"DRONE_HTTP_HOST":                  "autoscaler.drone.company.com",
		"DRONE_HTTP_PORT":                  "633eb230f5",
		"DRONE_HTTP_ROOT":                  "/autoscaler",
		"DRONE_AGENT_TOKEN":                "f5064039f5",
		"DRONE_AGENT_IMAGE":                "drone/drone-runner-docker:latest",
		"DRONE_AGENT_CONCURRENCY":          "2",
		"DRONE_AGENT_ARCH":                 "arm64",
		"DRONE_TLS_AUTOCERT":               "true",
		"DRONE_TLS_CERT":                   "/path/to/cert.crt",
		"DRONE_TLS_KEY":                    "/path/to/cert.key",
		"DRONE_PROMETHEUS_AUTH_TOKEN":      "b359e05e8",
		"DRONE_DATABASE_DRIVER":            "mysql",
		"DRONE_DATABASE_DATASOURCE":        "user:password@/dbname",
		"DRONE_DIGITALOCEAN_TOKEN":         "2573633eb",
		"DRONE_DIGITALOCEAN_IMAGE":         "docker-18-04",
		"DRONE_DIGITALOCEAN_REGION":        "ncy1",
		"DRONE_DIGITALOCEAN_SSHKEY":        "/path/to/ssh/key",
		"DRONE_DIGITALOCEAN_SIZE":          "s-1vcpu-1gb",
		"DRONE_DIGITALOCEAN_IPV6":          "true",
		"DRONE_DIGITALOCEAN_PRIVATE_IP":    "false",
		"DRONE_DIGITALOCEAN_FIREWALL":      "",
		"DRONE_DIGITALOCEAN_TAGS":          "drone,agent,prod",
		"DRONE_DIGITALOCEAN_USERDATA":      "#cloud-init",
		"DRONE_DIGITALOCEAN_USERDATA_FILE": "/path/to/cloud/init.yml",
		"DRONE_GOOGLE_ZONE":                "us-central1-b,us-central1-a",
		"DRONE_GOOGLE_MACHINE_TYPE":        "f1-micro",
		"DRONE_GOOGLE_MACHINE_IMAGE":       "ubuntu-1510-wily-v20151114",
		"DRONE_GOOGLE_DISK_TYPE":           "pd-standard",
		"DRONE_GOOGLE_NETWORK":             "default",
		"DRONE_GOOGLE_SUBNETWORK":          "",
		"DRONE_GOOGLE_PRIVATE_IP":          "false",
		"DRONE_GOOGLE_PREEMPTIBLE":         "true",
		"DRONE_GOOGLE_SCOPES":              "devstorage.read_only,pubsub",
		"DRONE_GOOGLE_DISK_SIZE":           "10",
		"DRONE_GOOGLE_PROJECT":             "project-foo",
		"DRONE_GOOGLE_TAGS":                "drone,agent,prod",
		"DRONE_GOOGLE_USERDATA":            "#cloud-init",
		"DRONE_GOOGLE_USERDATA_FILE":       "/path/to/cloud/init.yml",
		"DRONE_GOOGLE_READ_RATELIMIT":      "20",
		"DRONE_AMAZON_IMAGE":               "ami-07f84a50d2dec2fa4",
		"DRONE_AMAZON_INSTANCE":            "t3.medium",
		"DRONE_AMAZON_PRIVATE_IP":          "true",
		"DRONE_AMAZON_RETRIES":             "1",
		"DRONE_AMAZON_REGION":              "us-east-2",
		"DRONE_AMAZON_SSHKEY":              "id_rsa",
		"DRONE_AMAZON_SUBNET_ID":           "subnet-0b32177f",
		"DRONE_AMAZON_SECURITY_GROUP":      "sg-770eabe1",
		"DRONE_AMAZON_TAGS":                "os:linux,arch:amd64",
		"DRONE_AMAZON_USERDATA":            "#cloud-init",
		"DRONE_AMAZON_USERDATA_FILE":       "/path/to/cloud/init.yml",
		"DRONE_HETZNERCLOUD_TOKEN":         "12345678",
		"DRONE_HETZNERCLOUD_IMAGE":         "ubuntu-16.04",
		"DRONE_HETZNERCLOUD_DATACENTER":    "nbg1-dc3",
		"DRONE_HETZNERCLOUD_SSHKEY":        "12345",
		"DRONE_HETZNERCLOUD_TYPE":          "cx11",
		"DRONE_HETZNERCLOUD_USERDATA":      "#cloud-init",
		"DRONE_HETZNERCLOUD_USERDATA_FILE": "/path/to/cloud/init.yml",
		"DRONE_PACKET_APIKEY":              "12345678",
		"DRONE_PACKET_FACILITY":            "facility",
		"DRONE_PACKET_PROJECT_ID":          "project",
		"DRONE_PACKET_PLAN":                "plan",
		"DRONE_PACKET_OS":                  "ubuntu",
		"DRONE_PACKET_SSHKEY":              "id_rsa",
		"DRONE_PACKET_USERDATA":            "#cloud-init",
		"DRONE_PACKET_USERDATA_FILE":       "/path/to/cloud/init.yml",
		"DRONE_PACKET_HOSTNAME":            "agent",
		"DRONE_PACKET_TAGS":                "drone,agent,prod",
		"DRONE_OPENSTACK_NETWORK":          "my-subnet-1",
		"DRONE_OPENSTACK_IP_POOL":          "ext-ips-1",
		"DRONE_OPENSTACK_SSHKEY":           "drone-ci",
		"DRONE_OPENSTACK_SECURITY_GROUP":   "secgrp-feedface",
		"DRONE_OPENSTACK_FLAVOR":           "t1.medium",
		"DRONE_OPENSTACK_IMAGE":            "ubuntu-16.04-server-latest",
		"DRONE_OPENSTACK_METADATA":         "name:agent,owner:drone-ci",
		"DRONE_OPENSTACK_USERDATA":         "#cloud-init",
		"DRONE_OPENSTACK_USERDATA_FILE":    "/path/to/cloud/init.yml",
		"OS_REGION_NAME":                   "sto-01",
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
	err := json.Unmarshal(jsonConfig, &b)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(a, b) {
		t.Errorf("configuration mismatch")
		pretty.Ldiff(t, a, b)
	}
}

var jsonConfig = []byte(`{
  "Interval": 60000000000,
  "CapacityBuffer": 3,
  "Timeout": {
    "Stop": 3600000000000
  },
  "Slack": {
    "Webhook": "https://hooks.slack.com/services/XXX/YYY/ZZZ",
    "Create": false,
    "Destroy": false,
    "Error": true
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
    "OS": "linux",
    "Arch": "arm64",
    "Token": "f5064039f5",
    "Image": "drone/drone-runner-docker:latest",
    "Concurrency": 2,
    "KeepaliveTime": 360000000000,
    "KeepaliveTimeout": 30000000000,
    "NamePrefix": "agent-"
  },
  "HTTP": {
    "Proto": "http",
    "Host": "autoscaler.drone.company.com",
    "Port": "633eb230f5",
    "Root": "/autoscaler"
  },
  "UI": {
    "Realm": "Autoscaler"
  },
  "TLS": {
    "Autocert": true,
    "Cert": "/path/to/cert.crt",
    "Key": "/path/to/cert.key"
  },
  "Prometheus": {
    "AuthToken": "b359e05e8"
  },
  "Database": {
    "Driver": "mysql",
    "Datasource": "user:password@/dbname"
  },
  "DigitalOcean": {
    "Token": "2573633eb",
    "Image": "docker-18-04",
    "Region": "ncy1",
    "SSHKey": "/path/to/ssh/key",
    "Size": "s-1vcpu-1gb",
    "Tags": [
      "drone",
      "agent",
      "prod"
    ],
    "UserData": "#cloud-init",
    "UserDataFile": "/path/to/cloud/init.yml"
  },
  "Amazon": {
    "Image": "ami-07f84a50d2dec2fa4",
    "Instance": "t3.medium",
    "PrivateIP": true,
    "Retries": 1,
    "Region": "us-east-2",
    "SSHKey": "id_rsa",
    "SubnetID": [
		"subnet-0b32177f"
	],
    "SecurityGroup": [
      "sg-770eabe1"
    ],
    "tags": {
      "os": "linux",
      "arch": "amd64"
    },
    "UserData": "#cloud-init",
    "UserDataFile": "/path/to/cloud/init.yml"
  },
  "Google": {
    "Zone": ["us-central1-b","us-central1-a"],
    "MachineType": "f1-micro",
    "MachineImage": "ubuntu-1510-wily-v20151114",
    "DiskType": "pd-standard",
    "Address": "",
    "Network": "default",
    "Subnetwork": "",
    "Preemptible": true,
    "Scopes": [
      "devstorage.read_only",
      "pubsub"
    ],
    "DiskSize": 10,
    "Project": "project-foo",
    "Tags": [
      "drone",
      "agent",
      "prod"
    ],
    "UserData": "#cloud-init",
    "UserDataFile": "/path/to/cloud/init.yml",
    "UserDataKey": "user-data",
	"RateLimit": 20
  },
  "HetznerCloud": {
    "Token": "12345678",
    "Image": "ubuntu-16.04",
    "Datacenter": "nbg1-dc3",
    "SSHKey": 12345,
    "Type": "cx11",
    "UserData": "#cloud-init",
    "UserDataFile": "/path/to/cloud/init.yml"
  },
  "Packet": {
    "APIKey": "12345678",
    "Facility": "facility",
    "ProjectID": "project",
    "Plan": "plan",
    "OS": "ubuntu",
    "SSHKey": "id_rsa",
    "UserData": "#cloud-init",
    "UserDataFile": "/path/to/cloud/init.yml",
    "Hostname": "agent",
    "Tags": [
      "drone",
      "agent",
      "prod"
    ]
  },
  "OpenStack": {
    "Region": "sto-01",
    "Image": "ubuntu-16.04-server-latest",
    "Flavor": "t1.medium",
    "Network": "my-subnet-1",
    "Pool": "ext-ips-1",
    "SecurityGroup": [
      "secgrp-feedface"
    ],
    "SSHKey": "drone-ci",
    "Metadata": {
      "name": "agent",
      "owner": "drone-ci"
    },
    "UserData": "#cloud-init",
    "UserDataFile": "/path/to/cloud/init.yml"
  },
  "Watchtower": {
    "Image": "webhippie/watchtower",
    "Interval": 300,
    "Timeout": 7200000000000
  },
  "GC": {
    "Image": "drone/gc",
    "Interval": 1800000000000,
    "Cache": "10gb"
  },
  "Reaper": {
    "Interval": 3600000000000
  },
  "Pinger": {
    "Interval": 600000000000
  },
  "Check": {
    "Interval": 60000000000,
    "Deadline": 1800000000000
  }
}`)

func TestLoadEnvVariables(t *testing.T) {
	f, err := ioutil.TempFile("", "autoscaler-env-file-test")
	if err != nil {
		t.Error(err)
	}
	f.WriteString("ENV_FROM_FILE=FILE_VALUE")
	defer os.Remove(f.Name())

	environ := map[string]string{
		"ENV_FROM_HOST":        "HOST_VALUE",
		"DRONE_AGENT_ENVIRON":  `ENV=VALUE,ENV_FROM_HOST`,
		"DRONE_AGENT_ENV_FILE": f.Name(),
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
	want := []string{
		"ENV=VALUE",
		"ENV_FROM_HOST=HOST_VALUE",
		"ENV_FROM_FILE=FILE_VALUE",
	}
	if got, want := len(a.Agent.Environ), len(want); got != want {
		t.Errorf("Should have an environment of length %d, got %d", want, got)
	}
	for i := range a.Agent.Environ {
		if got, wantV := a.Agent.Environ[i], want[i]; got != wantV {
			t.Errorf("Wanted environ %s at index %d, got %s", wantV, i, got)
		}
	}
}
