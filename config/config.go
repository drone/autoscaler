// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package config

import "time"

type (
	// Config stores the configuration settings.
	Config struct {
		License  string
		Interval time.Duration `default:"5m"`

		Slack struct {
			Webhook string
		}

		Logs struct {
			Color  bool
			Debug  bool `default:"true"`
			Pretty bool
		}

		Pool struct {
			Min    int           `default:"2"`
			Max    int           `default:"4"`
			MinAge time.Duration `default:"55m" split_words:"true"`
		}

		Server struct {
			Host  string
			Proto string
			Token string
		}

		Agent struct {
			Host        string
			Token       string
			Image       string `default:"drone/agent:0.8"`
			Concurrency int    `default:"2"`
		}

		HTTP struct {
			Host string
			Port string `default:":8080"`
			Root string `default:"/"`
		}

		TLS struct {
			Autocert bool
			Cert     string
			Key      string
		}

		Prometheus struct {
			Token string
		}

		Database struct {
			Driver     string `default:"sqlite3"`
			Datasource string `default:"database.sqlite?cache=shared&mode=rwc&_busy_timeout=9999999"`
		}

		Amazon struct {
			Image         string
			Instance      string
			PrivateIP     bool `split_words:"true"`
			Region        string
			Retries       int
			SSHKey        string
			SubnetID      string   `split_words:"true"`
			SecurityGroup []string `split_words:"true"`
			Tags          map[string]string
			UserData      string `envconfig:"DRONE_AMAZON_USERDATA"`
			UserDataFile  string `envconfig:"DRONE_AMAZON_USERDATA_FILE"`
		}

		DigitalOcean struct {
			Token        string
			Image        string
			Region       string
			SSHKey       string
			Size         string
			Tags         []string
			UserData     string `envconfig:"DRONE_DIGITALOCEAN_USERDATA"`
			UserDataFile string `envconfig:"DRONE_DIGITALOCEAN_USERDATA_FILE"`
		}

		Google struct {
			Zone         string `default:"us-central1-a"`
			MachineType  string `split_words:"true" default:"n1-standard-1"`
			MachineImage string `split_words:"true" default:"ubuntu-1510-wily-v20151114"`
			DiskType     string `split_words:"true" default:"pd-standard"`
			Address      string
			Network      string `default:"default"`
			Subnetwork   string
			Preemptible  bool
			Scopes       string
			DiskSize     int `split_words:"true"`
			Project      string
			Tags         []string
			UserData     string `envconfig:"DRONE_GOOGLE_USERDATA"`
			UserDataFile string `envconfig:"DRONE_GOOGLE_USERDATA_FILE"`
		}

		HetznerCloud struct {
			Datacenter   string
			Image        string
			SSHKey       int
			Token        string
			Type         string
			UserData     string `envconfig:"DRONE_HETZNERCLOUD_USERDATA"`
			UserDataFile string `envconfig:"DRONE_HETZNERCLOUD_USERDATA_FILE"`
		}
	}
)
