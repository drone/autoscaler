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
		"DRONE_INTERVAL":                            "1m",
		"DRONE_SLACK_WEBHOOK":                       "https://hooks.slack.com/services/XXX/YYY/ZZZ",
		"DRONE_SLACK_CREATE":                        "false",
		"DRONE_SLACK_DESTROY":                       "false",
		"DRONE_LOGS_DEBUG":                          "true",
		"DRONE_LOGS_COLOR":                          "true",
		"DRONE_LOGS_PRETTY":                         "true",
		"DRONE_CAPACITY_BUFFER":                     "3",
		"DRONE_POOL_MIN_AGE":                        "1h",
		"DRONE_POOL_MIN":                            "1",
		"DRONE_POOL_MAX":                            "5",
		"DRONE_SERVER_HOST":                         "drone.company.com",
		"DRONE_SERVER_PROTO":                        "http",
		"DRONE_SERVER_TOKEN":                        "633eb230f5",
		"DRONE_HTTP_HOST":                           "autoscaler.drone.company.com",
		"DRONE_HTTP_PORT":                           "633eb230f5",
		"DRONE_HTTP_ROOT":                           "/autoscaler",
		"DRONE_AGENT_TOKEN":                         "f5064039f5",
		"DRONE_AGENT_IMAGE":                         "drone/drone-runner-docker:latest",
		"DRONE_AGENT_CONCURRENCY":                   "2",
		"DRONE_AGENT_ARCH":                          "arm64",
		"DRONE_TLS_AUTOCERT":                        "true",
		"DRONE_TLS_CERT":                            "/path/to/cert.crt",
		"DRONE_TLS_KEY":                             "/path/to/cert.key",
		"DRONE_PROMETHEUS_AUTH_TOKEN":               "b359e05e8",
		"DRONE_DATABASE_DRIVER":                     "mysql",
		"DRONE_DATABASE_DATASOURCE":                 "user:password@/dbname",
		"DRONE_DIGITALOCEAN_TOKEN":                  "2573633eb",
		"DRONE_DIGITALOCEAN_IMAGE":                  "docker-18-04",
		"DRONE_DIGITALOCEAN_REGION":                 "ncy1",
		"DRONE_DIGITALOCEAN_SSHKEY":                 "/path/to/ssh/key",
		"DRONE_DIGITALOCEAN_SIZE":                   "s-1vcpu-1gb",
		"DRONE_DIGITALOCEAN_IPV6":                   "true",
		"DRONE_DIGITALOCEAN_PRIVATE_IP":             "false",
		"DRONE_DIGITALOCEAN_FIREWALL":               "",
		"DRONE_DIGITALOCEAN_TAGS":                   "drone,agent,prod",
		"DRONE_DIGITALOCEAN_USERDATA":               "#cloud-init",
		"DRONE_DIGITALOCEAN_USERDATA_FILE":          "/path/to/cloud/init.yml",
		"DRONE_GOOGLE_ZONE":                         "us-central1-b,us-central1-a",
		"DRONE_GOOGLE_MACHINE_TYPE":                 "f1-micro",
		"DRONE_GOOGLE_MACHINE_IMAGE":                "ubuntu-1510-wily-v20151114",
		"DRONE_GOOGLE_DISK_TYPE":                    "pd-standard",
		"DRONE_GOOGLE_NETWORK":                      "default",
		"DRONE_GOOGLE_SUBNETWORK":                   "",
		"DRONE_GOOGLE_PRIVATE_IP":                   "false",
		"DRONE_GOOGLE_PREEMPTIBLE":                  "true",
		"DRONE_GOOGLE_SCOPES":                       "devstorage.read_only,pubsub",
		"DRONE_GOOGLE_DISK_SIZE":                    "10",
		"DRONE_GOOGLE_PROJECT":                      "project-foo",
		"DRONE_GOOGLE_TAGS":                         "drone,agent,prod",
		"DRONE_GOOGLE_USERDATA":                     "#cloud-init",
		"DRONE_GOOGLE_USERDATA_FILE":                "/path/to/cloud/init.yml",
		"DRONE_AMAZON_IMAGE":                        "ami-80ca47e6",
		"DRONE_AMAZON_INSTANCE":                     "t2.medium",
		"DRONE_AMAZON_PRIVATE_IP":                   "true",
		"DRONE_AMAZON_RETRIES":                      "1",
		"DRONE_AMAZON_REGION":                       "us-east-2",
		"DRONE_AMAZON_SSHKEY":                       "id_rsa",
		"DRONE_AMAZON_SUBNET_ID":                    "subnet-0b32177f",
		"DRONE_AMAZON_SECURITY_GROUP":               "sg-770eabe1",
		"DRONE_AMAZON_TAGS":                         "os:linux,arch:amd64",
		"DRONE_AMAZON_USERDATA":                     "#cloud-init",
		"DRONE_AMAZON_USERDATA_FILE":                "/path/to/cloud/init.yml",
		"DRONE_HETZNERCLOUD_TOKEN":                  "12345678",
		"DRONE_HETZNERCLOUD_IMAGE":                  "ubuntu-16.04",
		"DRONE_HETZNERCLOUD_DATACENTER":             "nbg1-dc3",
		"DRONE_HETZNERCLOUD_SSHKEY":                 "12345",
		"DRONE_HETZNERCLOUD_TYPE":                   "cx11",
		"DRONE_HETZNERCLOUD_USERDATA":               "#cloud-init",
		"DRONE_HETZNERCLOUD_USERDATA_FILE":          "/path/to/cloud/init.yml",
		"DRONE_PACKET_APIKEY":                       "12345678",
		"DRONE_PACKET_FACILITY":                     "facility",
		"DRONE_PACKET_PROJECT_ID":                   "project",
		"DRONE_PACKET_PLAN":                         "plan",
		"DRONE_PACKET_OS":                           "ubuntu",
		"DRONE_PACKET_SSHKEY":                       "id_rsa",
		"DRONE_PACKET_USERDATA":                     "#cloud-init",
		"DRONE_PACKET_USERDATA_FILE":                "/path/to/cloud/init.yml",
		"DRONE_PACKET_HOSTNAME":                     "agent",
		"DRONE_PACKET_TAGS":                         "drone,agent,prod",
		"DRONE_OPENSTACK_NETWORK":                   "my-subnet-1",
		"DRONE_OPENSTACK_IP_POOL":                   "ext-ips-1",
		"DRONE_OPENSTACK_SSHKEY":                    "drone-ci",
		"DRONE_OPENSTACK_SECURITY_GROUP":            "secgrp-feedface",
		"DRONE_OPENSTACK_FLAVOR":                    "t1.medium",
		"DRONE_OPENSTACK_IMAGE":                     "ubuntu-16.04-server-latest",
		"DRONE_OPENSTACK_METADATA":                  "name:agent,owner:drone-ci",
		"DRONE_OPENSTACK_USERDATA":                  "#cloud-init",
		"DRONE_OPENSTACK_USERDATA_FILE":             "/path/to/cloud/init.yml",
		"OS_REGION_NAME":                            "sto-01",
		"DRONE_YANDEX_CLOUD_TOKEN":                  "example token",
		"DRONE_YANDEX_CLOUD_SERVICE_ACCOUNT":        `{"id": "ajeboa0du6edu6m43c3t","service_account_id": "ajeq7dsmihqple6761c5","created_at": "2018-10-31T09:30:52Z","key_algorithm": "RSA_4096", "public_key": "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0q4HXY6/7jzA3iwofyTq\nxJ+VPR2fQGrhd+nEi32lemw8XEqVpr4wvzaJHdW2921Z7nsEKJPJq9ZCaDpnkdpS\npGPnJCUtJSYQvjs5yrYbEB00LGGi4pERDNWbsX+MyyMl+0Mqd3G3wiu/T8k31T5F\ngmOK1KPnlDQ6JjZ+OQWkBojBTGGkaCsYKuDwuKfsHZQCqhTt8pLIN7ZiiXWRIB4Q\n4GfBuBWUfhgncbNCj+PBEBvy1auFnI0CHQ8T9cHqnh9UQIi0qsxVslICv4Z5iX4y\nYCrRfSw3UJOqQ+mkttSNBjnJC7TpC4uQyc98XC+kLzP1i8/nNv967K9eWA6MVsHF\nZqAkFJYcUn6Bx/f3FDiIcW0tR5P/FgDtVTvQAAdUW+l02P27JOqRyD9oDX/y889/\nc1TGbXlhmaWCqjIVoUnUBlnlHAB7v8X4aqlCu9vwP0DUaXdI/Yxf7VcG6B7wFFZN\noM2k6X1J3LSMdrFTXSLbduv/n0mMLUurUx1D0YIIrk2Kv1N63YNiVtPWdFkGHQs6\nshVgrpiTBUm0VBME7EYKwKQUK7pZ0gn6/IeZpgel0aPCQtaF9FIffLi8KJaMVbJi\nNGgvr4HTejzn/jabWuLc3rN62AexNYUqnRMfmfNPXyArJ0A54tl2u/TKoPmw0w5t\nYAwgJ+mSGlylBJbZy2CBp2sCAwEAAQ==\n-----END PUBLIC KEY-----\n","private_key": "-----BEGIN PRIVATE KEY-----\n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\nawxDbYxnTunk1YxP7KA0/bDnu/OZlQ==\n-----END PRIVATE KEY-----\n"}`,
		"DRONE_YANDEX_CLOUD_FOLDER_ID":              "klfsd8gs9fisd9f",
		"DRONE_YANDEX_CLOUD_SUBNET_ID":              "asdasf6a7s6f7a6",
		"DRONE_YANDEX_CLOUD_ZONE":                   "ru-central1-a,ru-central1-b",
		"DRONE_YANDEX_CLOUD_DISK_SIZE":              "10737418240",
		"DRONE_YANDEX_CLOUD_DISK_TYPE":              "network-hdd",
		"DRONE_YANDEX_CLOUD_RESOURCE_CORES":         "2",
		"DRONE_YANDEX_CLOUD_RESOURCE_MEMORY":        "1073741824",
		"DRONE_YANDEX_CLOUD_RESOURCE_CORE_FRACTION": "100",
		"DRONE_YANDEX_CLOUD_PLATFORM_ID":            "standard-v3",
		"DRONE_YANDEX_CLOUD_IMAGE_FOLDER_ID":        "standard-images",
		"DRONE_YANDEX_CLOUD_IMAGE_FAMILY":           "debian-9",
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
    "Image": "ami-80ca47e6",
    "Instance": "t2.medium",
    "PrivateIP": true,
    "Retries": 1,
    "Region": "us-east-2",
    "SSHKey": "id_rsa",
    "SubnetID": "subnet-0b32177f",
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
    "UserDataKey": "user-data"
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
  },
  "YandexCloud": {
    "token": "example token",
	"serviceAccount": "{\"id\": \"ajeboa0du6edu6m43c3t\",\"service_account_id\": \"ajeq7dsmihqple6761c5\",\"created_at\": \"2018-10-31T09:30:52Z\",\"key_algorithm\": \"RSA_4096\", \"public_key\": \"-----BEGIN PUBLIC KEY-----\\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0q4HXY6/7jzA3iwofyTq\\nxJ+VPR2fQGrhd+nEi32lemw8XEqVpr4wvzaJHdW2921Z7nsEKJPJq9ZCaDpnkdpS\\npGPnJCUtJSYQvjs5yrYbEB00LGGi4pERDNWbsX+MyyMl+0Mqd3G3wiu/T8k31T5F\\ngmOK1KPnlDQ6JjZ+OQWkBojBTGGkaCsYKuDwuKfsHZQCqhTt8pLIN7ZiiXWRIB4Q\\n4GfBuBWUfhgncbNCj+PBEBvy1auFnI0CHQ8T9cHqnh9UQIi0qsxVslICv4Z5iX4y\\nYCrRfSw3UJOqQ+mkttSNBjnJC7TpC4uQyc98XC+kLzP1i8/nNv967K9eWA6MVsHF\\nZqAkFJYcUn6Bx/f3FDiIcW0tR5P/FgDtVTvQAAdUW+l02P27JOqRyD9oDX/y889/\\nc1TGbXlhmaWCqjIVoUnUBlnlHAB7v8X4aqlCu9vwP0DUaXdI/Yxf7VcG6B7wFFZN\\noM2k6X1J3LSMdrFTXSLbduv/n0mMLUurUx1D0YIIrk2Kv1N63YNiVtPWdFkGHQs6\\nshVgrpiTBUm0VBME7EYKwKQUK7pZ0gn6/IeZpgel0aPCQtaF9FIffLi8KJaMVbJi\\nNGgvr4HTejzn/jabWuLc3rN62AexNYUqnRMfmfNPXyArJ0A54tl2u/TKoPmw0w5t\\nYAwgJ+mSGlylBJbZy2CBp2sCAwEAAQ==\\n-----END PUBLIC KEY-----\\n\",\"private_key\": \"-----BEGIN PRIVATE KEY-----\\nMIIJQgIBADANBgkqhkiG9w0BAQEFAASCCSwwggkoAgEAAoICAQDSrgddjr/uPMDe\\nLCh/JOrEn5U9HZ9AauF36cSLfaV6bDxcSpWmvjC/Nokd1bb3bVnuewQok8mr1kJo\\nOmeR2lKkY+ckJS0lJhC+OznKthsQHTQsYaLikREM1Zuxf4zLIyX7Qyp3cbfCK79P\\nyTfVPkWCY4rUo+eUNDomNn45BaQGiMFMYaRoKxgq4PC4p+wdlAKqFO3yksg3tmKJ\\ndZEgHhDgZ8G4FZR+GCdxs0KP48EQG/LVq4WcjQIdDxP1weqeH1RAiLSqzFWyUgK/\\nhnmJfjJgKtF9LDdQk6pD6aS21I0GOckLtOkLi5DJz3xcL6QvM/WLz+c2/3rsr15Y\\nDoxWwcVmoCQUlhxSfoHH9/cUOIhxbS1Hk/8WAO1VO9AAB1Rb6XTY/bsk6pHIP2gN\\nf/Lzz39zVMZteWGZpYKqMhWhSdQGWeUcAHu/xfhqqUK72/A/QNRpd0j9jF/tVwbo\\nHvAUVk2gzaTpfUnctIx2sVNdItt26/+fSYwtS6tTHUPRggiuTYq/U3rdg2JW09Z0\\nWQYdCzqyFWCumJMFSbRUEwTsRgrApBQrulnSCfr8h5mmB6XRo8JC1oX0Uh98uLwo\\nloxVsmI0aC+vgdN6POf+Npta4tzes3rYB7E1hSqdEx+Z809fICsnQDni2Xa79Mqg\\n+bDTDm1gDCAn6ZIaXKUEltnLYIGnawIDAQABAoICACXvsmHVZ5gllnErMGucoS2g\\nssXbhKab2Fe4X2zixh5iSQgxYfsxeiOkVVJq/lRVe4Em45vO6NypazHLeoTX9FOn\\nraJjk1qCHTe0AHcRDZR8Pb3UIvl7N7/A4xU2K4sUnC0/bfEuJ/Gt4Pgj+orKeMe+\\n1uvtS7DzKplg7J+l9WA71drEJk+fmu11rcMCcdDtqwEnXaV1atolXF72LZjD8TQH\\nWumj8SY3gTrHFbBFSal17ucsyJVlCsFiyqxRK8cnSwuH0kiDHNdMTzRfqZjpgXax\\nnyFUCe3XeSxbcQ5+/ZnmY95YyDINAphkZTdQWNcrGwb++9p6bI8cEPf4PqsMn1fE\\ngCvhjTpfVLhTbjvA6QH1Za6j4BQxR3+zUw1yzCsuma00r9m1H0/4aQwBhRN6bbKs\\nf2OwPmwSMtqzrklqODmLSQbXyPIkb6xoky6zRNJdq3uwBORILgOvzNiM4jUrBjue\\nfUmlblu7n2Qw2POJRjrNGEUasYUhZCk9UAlGl8ThTqMo5r2pfT0p0MuQrCzsg4ch\\nylM7vzqKhrGdHlIEBr//AJ2C/NHqZjiX2h3H9gfoFiNCx0meOtjDmen9z6vCtiOs\\n7JcAHO30z5KuRd59nPO3GuWeXe6SlWWymgv9DiNAQXx163lyBpYvXI8F8KzgKvM+\\nbtQzUXpr/Q2/KOcSgl9hAoIBAQD9Vo4PZQ9FOzJSxXPOHeMbpP3lizHoPwmDbRjT\\nCu5d8+a9IuHmsINm5bzShPJQauZMsaJvyfHd6Q1RlieJ6l/qObSHrN9LNwVYZxJ8\\nTZiyufk4FrgPX0K6rzNUHxCG/2R5A/qfg32BgQ8+h+rK0ZSPVGnVUMRpofTPubBE\\nP9lzy/3pcqvJjljI2BtbFYYanB4WWrtDCEZOcLU6Nny7ll/N5Y1ZFFcrL1eiB2kh\\nFvFp1ScgAtjTWoxEPIPJK0hVUkMnWnbMMCpJ+WXg+2blDEY2H8BTk47EyN7P/Y3G\\nFhujduNLCJbwCuGn8wchCgvSFuy56UcbLpdu8e1R/TEpL5vxAoIBAQDU5LqS7TyS\\n0p7tOXQuLgNuPXezuu8yC/eIg9wfUWPEMNaFS6er75AKOqzykFMHnNimPnuxMNLI\\nmXwb2bd6J9PfRJ6eqsf60E1iDD+F9KiKSd+btGxh5soc8wVjTWbghW47DfWkAZGX\\nw0dWZTh2PY41HupU6lizrBMPnx2xFQz4CgU3DjDlhSVDSbkJ64Wt0RDKICe+OQWe\\nlhKmC8fjo9IO5W/QRIcOmqO+1W3nUko84CDVlh+mJ2Hxi2N3+xUT9XJee0id81u2\\nVO6tH20+zDTfsu8RFJXBL2Mx2WJ4MTV3ASxvBbUMG5sAOlFb4SL6p6ODzZwOwlDZ\\ny/T3tQHr+IUbAoIBAAp9vSBSFRHO48SdvK/6eN86M/F/lC+D/Mbei7qhp0FoylNm\\n0GgXQznNpcYqD0bZRnRCnvF2MXf5IL4SM8z4UcSHYzyDIjQhMS16Bz/yjrJIFVQH\\nTNQGI+NLQhrntm2Awg5o5cYZUec9Cv6R7l071KUi38cfsyKUvGilzfDlnAG5nug+\\nAXM1W+PlXyykdYtAj9ZpJ3wdKZwx+q9QdlXmYk1KhlH8D6gQK9bf67CdHJ4/X4Fp\\n3MTT6R8iSmrYSgSOhY1pp6XJENdDZr6sapRtr7KqGfLcF3t6vg9q9qYPYFGiqMMA\\ntg92w+WKoO7zVY37uQ3x5SnxAgBsMGHG1HRaLmECggEBAIEuVYQIDkRtJ2B9B2Fq\\nLEy9YaAeozv0BPzCPlSGl4oZtGHnuVNcJ0P9vKnnJ2qsIs4lhfrLzGtKrwNbRbkK\\n58ZHphRTPsuTkBEZq4YGIirfjp61iTqSxztvv2o1MmK0tGGDI/Wjugujw+rJuswM\\np/jVzI1AMhi8JkjJXUPxqQ/tTKLOqp7q/uRonK5HSrNg89YiUttbUGydVa2J4n3g\\nDvtY/1MZ8fXLoeaPLYQ6668qtOHFmWjB5u2hjfbk1TJqMj7ggfzOCW2G9dj5A9oi\\nIUdIFUaA/ineLku2Q8j42x9eB+9KQESbj59Aw9ODtizwggjdP3+5K0QtPXT9UbA0\\n+dcCggEAfpDfGOMurnFZxh7AYU3HfFLR0LIDqc4JX/SA8WlsbFUfM8ujaHhKoR+5\\nnWWTouuOy8lJlXVnqUfKvG4Ty0+2QvcTFE50h167AewsHmDqJ4oELJ6kCbMEUzpk\\nzILaeiCJlbldkfi4ztA7hT8Dfv+yKmi9GA2pyoMbdsVwG4xPDkA/R0jj7H5kkrh1\\nAv/K+T674XEr0ReHxEIxRFFQ0K/lyOxRIdxGssb30SNS3VvKTFtvFDKTm0uP7MYD\\ndSc0bk6fmeN0bR/Og2/S1ZEkQNxUFBPx92e9T4g/bi/2rIOdl4xcpwjW1By2UkNG\\nawxDbYxnTunk1YxP7KA0/bDnu/OZlQ==\\n-----END PRIVATE KEY-----\\n\"}",
	"folderID": "klfsd8gs9fisd9f",
	"subnetID": "asdasf6a7s6f7a6",
	"zone": ["ru-central1-a","ru-central1-b"],
    "diskSize": 10737418240,
    "diskType": "network-hdd",
    "resourceCores": 2,
    "resourceMemory": 1073741824,
	"resourceCoreFraction": 100, 
    "platformID": "standard-v3",
    "imageFolderID": "standard-images",
	"imageFamily": "debian-9",
	"preemptible": false
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
