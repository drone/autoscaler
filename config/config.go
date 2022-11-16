// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package config

import "time"

// TODO change DRONE_HTTP_PORT to DRONE_HTTP_BIND

type (
	// Config stores the configuration settings.
	Config struct {
		Interval       time.Duration `default:"5m"`
		CapacityBuffer int           `default:"0" split_words:"true"`

		Timeout struct {
			Stop time.Duration `envconfig:"DRONE_TIMEOUT_STOP" default:"1h"`
		}

		Slack struct {
			Webhook string
			Create  bool `default:"true"`
			Destroy bool `default:"true"`
			Error   bool `default:"true"`
		}

		Logs struct {
			Debug  bool `default:"true"`
			Trace  bool
			Pretty bool
		}

		Pool struct {
			Min    int           `default:"2"`
			Max    int           `default:"4"`
			MinAge time.Duration `default:"55m" split_words:"true"`
		}

		Check struct {
			Interval time.Duration `envconfig:"DRONE_INSTALL_CHECK_INTERVAL" default:"1m"`
			Deadline time.Duration `envconfig:"DRONE_INSTALL_CHECK_DEADLINE" default:"30m"`
		}

		Server struct {
			Host  string
			Proto string
			Token string
		}

		Agent struct {
			Token       string
			Image       string `default:"drone/drone-runner-docker:1"`
			Concurrency int    `default:"2"`
			OS          string `default:"linux"`
			Arch        string `default:"amd64"`
			Version     string
			Kernel      string
			EnvironFile string `envconfig:"DRONE_AGENT_ENV_FILE"`
			Environ     []string
			Volumes     []string
			Ports       []string          `envconfig:"DRONE_AGENT_PUBLISHED_PORTS"`
			Labels      map[string]string `envconfig:"DRONE_AGENT_LABELS"`
			NamePrefix  string            `envconfig:"DRONE_AGENT_NAME_PREFIX" default:"agent-"`
		}

		Runner Runner

		GC struct {
			Enabled  bool          `envconfig:"DRONE_GC_ENABLED"`
			Image    string        `envconfig:"DRONE_GC_IMAGE" default:"drone/gc"`
			Debug    bool          `envconfig:"DRONE_GC_DEBUG"`
			Images   []string      `envconfig:"DRONE_GC_IGNORE_IMAGES"`
			Interval time.Duration `envconfig:"DRONE_GC_INTERVAL" default:"30m"`
			Cache    string        `envconfig:"DRONE_GC_CACHE" default:"10gb"`
		}

		Reaper struct {
			Enabled  bool          `envconfig:"DRONE_REAPER_ENABLED", default:"false"`
			Interval time.Duration `envconfig:"DRONE_REAPER_INTERVAL" default:"1h"`
		}

		Pinger struct {
			Enabled  bool          `envconfig:"DRONE_PINGER_ENABLED", default:"false"`
			Interval time.Duration `envconfig:"DRONE_PINGER_INTERVAL" default:"10m"`
		}

		Watchtower struct {
			Enabled  bool          `envconfig:"DRONE_WATCHTOWER_ENABLED"`
			Image    string        `envconfig:"DRONE_WATCHTOWER_IMAGE" default:"webhippie/watchtower"`
			Interval int           `envconfig:"DRONE_WATCHTOWER_INTERVAL" default:"300"`
			Timeout  time.Duration `envconfig:"DRONE_WATCHTOWER_TIMEOUT" default:"120m"`
		}

		HTTP struct {
			Proto string `envconfig:"DRONE_HTTP_PROTO" default:"http"`
			Host  string `envconfig:"DRONE_HTTP_HOST"`
			Port  string `envconfig:"DRONE_HTTP_PORT" default:":8080"`
			Root  string `envconfig:"DRONE_HTTP_ROOT" default:"/"`
		}

		UI struct {
			Username string `envconfig:"DRONE_UI_USERNAME"`
			Password string `envconfig:"DRONE_UI_PASSWORD"`
			Realm    string `envconfig:"DRONE_UI_REALM" default:"Autoscaler"`
		}

		TLS struct {
			Autocert bool
			Cert     string
			Key      string
		}

		Prometheus struct {
			AuthToken string `split_words:"true"`
		}

		Database struct {
			Driver      string        `default:"sqlite3"`
			Datasource  string        `default:"database.sqlite?cache=shared&mode=rwc&_busy_timeout=9999999"`
			MaxIdle     int           `envconfig:"DRONE_DATABASE_MAX_IDLE" default:"0"`
			MaxLifetime time.Duration `envconfig:"DRONE_DATABASE_MAX_LIFETIME"`
		}

		Amazon struct {
			DeviceName    string `envconfig:"DRONE_AMAZON_DEVICE_NAME"`
			Image         string `envconfig:"DRONE_AMAZON_IMAGE"`
			Instance      string `envconfig:"DRONE_AMAZON_INSTANCE"`
			InstanceAlt   string `envconfig:"DRONE_AMAZON_INSTANCE_ALT"`
			PrivateIP     bool   `split_words:"true"`
			Region        string
			Retries       int
			SSHKey        string
			SubnetID      string   `split_words:"true"`
			SecurityGroup []string `split_words:"true"`
			Tags          map[string]string
			UserData      string `envconfig:"DRONE_AMAZON_USERDATA"`
			UserDataFile  string `envconfig:"DRONE_AMAZON_USERDATA_FILE"`
			VolumeSize    int64  `envconfig:"DRONE_AMAZON_VOLUME_SIZE"`
			VolumeType    string `envconfig:"DRONE_AMAZON_VOLUME_TYPE"`
			VolumeIops    int64  `envconfig:"DRONE_AMAZON_VOLUME_IOPS"`
			IamProfileArn string `envconfig:"DRONE_AMAZON_IAM_PROFILE_ARN"`
			MarketType    string `envconfig:"DRONE_AMAZON_MARKET_TYPE"`
		}

		DigitalOcean struct {
			Token        string
			Image        string
			Region       string
			SSHKey       string
			Size         string
			Firewall     string
			Tags         []string
			PrivateIP    bool   `split_words:"true"`
			UserData     string `envconfig:"DRONE_DIGITALOCEAN_USERDATA"`
			UserDataFile string `envconfig:"DRONE_DIGITALOCEAN_USERDATA_FILE"`
		}

		Google struct {
			MachineType         string            `envconfig:"DRONE_GOOGLE_MACHINE_TYPE"`
			MachineImage        string            `envconfig:"DRONE_GOOGLE_MACHINE_IMAGE"`
			Network             string            `envconfig:"DRONE_GOOGLE_NETWORK"`
			Subnetwork          string            `envconfig:"DRONE_GOOGLE_SUBNETWORK"`
			Labels              map[string]string `envconfig:"DRONE_GOOGLE_LABELS"`
			Scopes              []string          `envconfig:"DRONE_GOOGLE_SCOPES"`
			ServiceAccountEmail string            `envconfig:"DRONE_GOOGLE_SERVICE_ACCOUNT_EMAIL"`
			DiskSize            int64             `envconfig:"DRONE_GOOGLE_DISK_SIZE"`
			DiskType            string            `envconfig:"DRONE_GOOGLE_DISK_TYPE"`
			Project             string            `envconfig:"DRONE_GOOGLE_PROJECT"`
			PrivateIP           bool              `split_words:"true"`
			Tags                []string          `envconfig:"DRONE_GOOGLE_TAGS"`
			UserData            string            `envconfig:"DRONE_GOOGLE_USERDATA"`
			UserDataFile        string            `envconfig:"DRONE_GOOGLE_USERDATA_FILE"`
			Zone                []string          `envconfig:"DRONE_GOOGLE_ZONE"`
			UserDataKey         string            `envconfig:"DRONE_GOOGLE_USERDATA_KEY" default:"user-data"`
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

		Packet struct {
			APIKey       string
			Facility     string
			Plan         string
			OS           string
			ProjectID    string `split_words:"true"`
			Tags         []string
			SSHKey       string
			UserData     string `envconfig:"DRONE_PACKET_USERDATA"`
			UserDataFile string `envconfig:"DRONE_PACKET_USERDATA_FILE"`
			Hostname     string
		}

		OpenStack struct {
			Region        string `envconfig:"OS_REGION_NAME"`
			Image         string
			Flavor        string
			Network       string
			Pool          string   `envconfig:"DRONE_OPENSTACK_IP_POOL"`
			SecurityGroup []string `split_words:"true"`
			SSHKey        string
			Metadata      map[string]string
			UserData      string `envconfig:"DRONE_OPENSTACK_USERDATA"`
			UserDataFile  string `envconfig:"DRONE_OPENSTACK_USERDATA_FILE"`
		}

		Scaleway struct {
			AccessKey      string `split_words:"true"`
			SecretKey      string `split_words:"true"`
			OrganisationID string `split_words:"true"`
			Zone           string
			Size           string
			Image          string
			DynamicIP      bool `split_words:"true"`
			Tags           []string
			UserData       string `envconfig:"DRONE_SCALEWAY_USERDATA"`
			UserDataFile   string `envconfig:"DRONE_SCALEWAY_USERDATA_FILE"`
		}

		YandexCloud struct {
			Token                string   `envconfig:"DRONE_YANDEX_CLOUD_TOKEN"`
			ServiceAccount       string   `envconfig:"DRONE_YANDEX_CLOUD_SERVICE_ACCOUNT"`
			FolderID             string   `envconfig:"DRONE_YANDEX_CLOUD_FOLDER_ID"`
			SubnetID             string   `envconfig:"DRONE_YANDEX_CLOUD_SUBNET_ID"`
			Zone                 []string `envconfig:"DRONE_YANDEX_CLOUD_ZONE"`
			DiskSize             int64    `envconfig:"DRONE_YANDEX_CLOUD_DISK_SIZE"`
			DiskType             string   `envconfig:"DRONE_YANDEX_CLOUD_DISK_TYPE"`
			ResourceCores        int64    `envconfig:"DRONE_YANDEX_CLOUD_RESOURCE_CORES"`
			ResourceMemory       int64    `envconfig:"DRONE_YANDEX_CLOUD_RESOURCE_MEMORY"`
			ResourceCoreFraction int64    `envconfig:"DRONE_YANDEX_CLOUD_RESOURCE_CORE_FRACTION"`
			Preemptible          bool     `envconfig:"DRONE_YANDEX_CLOUD_PREEMPTIBLE"`
			PrivateIP            bool     `envconfig:"DRONE_YANDEX_CLOUD_PRIVATE_IP"`
			PlatformID           string   `envconfig:"DRONE_YANDEX_CLOUD_PLATFORM_ID"`
			ImageFolderID        string   `envconfig:"DRONE_YANDEX_CLOUD_IMAGE_FOLDER_ID"`
			ImageFamily          string   `envconfig:"DRONE_YANDEX_CLOUD_IMAGE_FAMILY"`
			SecurityGroupIDs     []string `envconfig:"DRONE_YANDEX_CLOUD_SECURITY_GROUP_IDS"`
			SSHUserKeyPair       string   `envconfig:"DRONE_YANDEX_CLOUD_SSH_USER_KEY_PAIR"`
			DockerComposeConfig  string   `envconfig:"DRONE_YANDEX_CLOUD_DOCKER_COMPOSE"`
		}
	}

	Runner struct {
		Volumes    string
		Devices    string
		Privileged string
		EnvFile    string
	}
)
