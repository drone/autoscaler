// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/drone/signal"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/amazon"
	"github.com/drone/autoscaler/drivers/digitalocean"
	"github.com/drone/autoscaler/drivers/google"
	"github.com/drone/autoscaler/drivers/hetznercloud"
	"github.com/drone/autoscaler/drivers/openstack"
	"github.com/drone/autoscaler/drivers/packet"
	"github.com/drone/autoscaler/drivers/scaleway"
	"github.com/drone/autoscaler/drivers/yandexcloud"
	"github.com/drone/autoscaler/engine"
	"github.com/drone/autoscaler/logger"
	"github.com/drone/autoscaler/logger/history"
	"github.com/drone/autoscaler/logger/request"
	"github.com/drone/autoscaler/metrics"
	"github.com/drone/autoscaler/server"
	"github.com/drone/autoscaler/server/web"
	"github.com/drone/autoscaler/server/web/static"
	"github.com/drone/autoscaler/slack"
	"github.com/drone/autoscaler/store"

	"github.com/99designs/basicauth-go"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	source  = "https://github.com/drone/autoscaler.git"
	version string
	commit  string
)

func main() {
	conf := config.MustLoad()
	setupLogging(conf)

	provider, err := setupProvider(conf)
	if err != nil {
		logrus.WithError(err).
			Fatalln("Invalid or missing hosting provider")
	}

	// instruments the provider with prometheus metrics.
	provider = metrics.ServerCreate(provider)
	provider = metrics.ServerDelete(provider)

	db, err := store.Connect(
		conf.Database.Driver,
		conf.Database.Datasource,
		conf.Database.MaxIdle,
		conf.Database.MaxLifetime,
	)
	if err != nil {
		logrus.WithError(err).
			Fatalln("Cannot establish database connection")
	}

	mu := store.NewLocker(conf.Database.Driver)
	servers := store.NewServerStore(db, mu)

	// instruments the provider with slack notifications
	// instance creation and termination events.
	if conf.Slack.Webhook != "" {
		servers = slack.New(conf, servers)
	}
	servers = metrics.ServerCount(servers)
	defer db.Close()

	client := setupClient(conf)

	enginex := engine.New(
		client,
		conf,
		servers,
		provider,
		metrics.New(),
	)

	//
	// Setup the router
	//

	r := chi.NewRouter()
	r.Use(request.Logger)

	// middleware to require basic authentication.
	auth := basicauth.New(conf.UI.Realm, map[string][]string{
		conf.UI.Username: {conf.UI.Password},
	})

	r.Route(conf.HTTP.Root, func(root chi.Router) {
		// handler to serve static assets for the dashboard.
		fs := http.FileServer(static.New())

		root.Handle("/", http.RedirectHandler("/ui", http.StatusSeeOther))
		root.Get("/metrics", server.HandleMetrics(conf.Prometheus.AuthToken))
		root.Get("/version", server.HandleVersion(source, version, commit))
		root.Get("/healthz", server.HandleHealthz())
		root.Get("/varz", server.HandleVarz(enginex))
		root.Handle("/static/*", http.StripPrefix("/static/", fs))

		if conf.UI.Password != "" {
			// register the history handler
			history := history.New()
			logrus.AddHook(history)

			root.Route("/ui", func(ui chi.Router) {
				ui.Use(auth)
				ui.Get("/", web.HandleServers(servers))
				ui.Get("/logs", web.HandleLogging(history))
			})
		}
		root.Route("/api", func(api chi.Router) {
			api.Use(server.CheckDrone(conf))

			api.Post("/pause", server.HandleEnginePause(enginex))
			api.Post("/resume", server.HandleEngineResume(enginex))
			api.Get("/servers", server.HandleServerList(servers))
			api.Post("/servers", server.HandleServerCreate(servers, conf))
			api.Get("/servers/{name}", server.HandleServerFind(servers))
			api.Delete("/servers/{name}", server.HandleServerDelete(servers))
		})
	})

	//
	// starts the web server.
	//

	srv := &http.Server{
		Handler: r,
	}

	ctx := context.Background()
	ctx = signal.WithContextFunc(ctx, func() {
		srv.Shutdown(ctx)
	})

	var g errgroup.Group
	g.Go(func() error {
		if conf.TLS.Autocert {
			return srv.Serve(
				autocert.NewListener(conf.HTTP.Host),
			)
		} else if conf.TLS.Cert != "" {
			return srv.ListenAndServeTLS(
				conf.TLS.Cert,
				conf.TLS.Key,
			)
		}
		srv.Addr = conf.HTTP.Port

		logrus.WithField("addr", conf.HTTP.Port).
			Infoln("starting the server")

		return srv.ListenAndServe()
	})

	//
	// starts the auto-scaler routine.
	//

	g.Go(func() error {
		enginex.Start(ctx)
		return nil
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Fatalln("Program terminated")
	}
}

// helper funciton configures the logging.
func setupLogging(c config.Config) {
	logger.Default = logger.Logrus(
		logrus.NewEntry(
			logrus.StandardLogger(),
		),
	)
	if c.Logs.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if c.Logs.Trace {
		logrus.SetLevel(logrus.TraceLevel)
	}
	if c.Logs.Pretty == false {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

// helper function configures the drone client.
func setupClient(c config.Config) drone.Client {
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: c.Server.Token,
		},
	)
	uri := new(url.URL)
	uri.Scheme = c.Server.Proto
	uri.Host = c.Server.Host
	return drone.NewClient(uri.String(), auther)
}

// helper function configures the hosting provider.
func setupProvider(c config.Config) (autoscaler.Provider, error) {
	switch {
	case c.Google.Project != "":
		return google.New(
			google.WithDiskSize(c.Google.DiskSize),
			google.WithDiskType(c.Google.DiskType),
			google.WithMachineImage(c.Google.MachineImage),
			google.WithMachineType(c.Google.MachineType),
			google.WithLabels(c.Google.Labels),
			google.WithNetwork(c.Google.Network),
			google.WithSubnetwork(c.Google.Subnetwork),
			google.WithPrivateIP(c.Google.PrivateIP),
			google.WithServiceAccountEmail(c.Google.ServiceAccountEmail),
			google.WithProject(c.Google.Project),
			google.WithTags(c.Google.Tags...),
			google.WithScopes(c.Google.Scopes...),
			google.WithUserData(c.Google.UserData),
			google.WithUserDataFile(c.Google.UserDataFile),
			google.WithZones(c.Google.Zone...),
			google.WithUserDataKey(c.Google.UserDataKey),
		)
	case c.DigitalOcean.Token != "":
		return digitalocean.New(
			digitalocean.WithSSHKey(c.DigitalOcean.SSHKey),
			digitalocean.WithImage(c.DigitalOcean.Image),
			digitalocean.WithRegion(c.DigitalOcean.Region),
			digitalocean.WithSize(c.DigitalOcean.Size),
			digitalocean.WithFirewall(c.DigitalOcean.Firewall),
			digitalocean.WithUserDataFile(c.DigitalOcean.UserDataFile),
			digitalocean.WithUserData(c.DigitalOcean.UserData),
			digitalocean.WithToken(c.DigitalOcean.Token),
			digitalocean.WithPrivateIP(c.DigitalOcean.PrivateIP),
			digitalocean.WithTags(c.DigitalOcean.Tags...),
		), nil
	case c.Scaleway.AccessKey != "":
		return scaleway.New(
			scaleway.WithAccessKey(c.Scaleway.AccessKey),
			scaleway.WithSecretKey(c.Scaleway.SecretKey),
			scaleway.WithOrganisationID(c.Scaleway.OrganisationID),
			scaleway.WithZone(c.Scaleway.Zone),
			scaleway.WithSize(c.Scaleway.Size),
			scaleway.WithImage(c.Scaleway.Image),
			scaleway.WithDynamicIP(c.Scaleway.DynamicIP),
			scaleway.WithTags(c.Scaleway.Tags...),
			scaleway.WithUserData(c.Scaleway.UserData),
			scaleway.WithUserDataFile(c.Scaleway.UserDataFile),
		)
	case c.HetznerCloud.Token != "":
		return hetznercloud.New(
			hetznercloud.WithDatacenter(c.HetznerCloud.Datacenter),
			hetznercloud.WithImage(c.HetznerCloud.Image),
			hetznercloud.WithUserDataFile(c.HetznerCloud.UserDataFile),
			hetznercloud.WithUserData(c.HetznerCloud.UserData),
			hetznercloud.WithServerType(c.HetznerCloud.Type),
			hetznercloud.WithSSHKey(c.HetznerCloud.SSHKey),
			hetznercloud.WithToken(c.HetznerCloud.Token),
		), nil
	case c.Packet.APIKey != "":
		return packet.New(
			packet.WithAPIKey(c.Packet.APIKey),
			packet.WithFacility(c.Packet.Facility),
			packet.WithProject(c.Packet.ProjectID),
			packet.WithPlan(c.Packet.Plan),
			packet.WithOS(c.Packet.OS),
			packet.WithSSHKey(c.Packet.SSHKey),
			packet.WithUserData(c.Packet.UserData),
			packet.WithUserDataFile(c.Packet.UserDataFile),
			packet.WithHostname(c.Packet.Hostname),
			packet.WithTags(c.Packet.Tags...),
		), nil
	case os.Getenv("AWS_ACCESS_KEY_ID") != "" || os.Getenv("AWS_IAM") != "":
		return amazon.New(
			amazon.WithDeviceName(c.Amazon.DeviceName),
			amazon.WithImage(c.Amazon.Image),
			amazon.WithRegion(c.Amazon.Region),
			amazon.WithRetries(c.Amazon.Retries),
			amazon.WithPrivateIP(c.Amazon.PrivateIP),
			amazon.WithSSHKey(c.Amazon.SSHKey),
			amazon.WithSecurityGroup(c.Amazon.SecurityGroup...),
			amazon.WithSize(c.Amazon.Instance),
			amazon.WithSizeAlt(c.Amazon.InstanceAlt),
			amazon.WithSubnet(c.Amazon.SubnetID),
			amazon.WithTags(c.Amazon.Tags),
			amazon.WithUserData(c.Amazon.UserData),
			amazon.WithUserDataFile(c.Amazon.UserDataFile),
			amazon.WithVolumeSize(c.Amazon.VolumeSize),
			amazon.WithVolumeType(c.Amazon.VolumeType),
			amazon.WithVolumeIops(c.Amazon.VolumeIops),
			amazon.WithIamProfileArn(c.Amazon.IamProfileArn),
			amazon.WithMarketType(c.Amazon.MarketType),
		), nil
	case os.Getenv("OS_USERNAME") != "":
		return openstack.New(
			openstack.WithImage(c.OpenStack.Image),
			openstack.WithRegion(c.OpenStack.Region),
			openstack.WithFlavor(c.OpenStack.Flavor),
			openstack.WithNetwork(c.OpenStack.Network),
			openstack.WithFloatingIpPool(c.OpenStack.Pool),
			openstack.WithSSHKey(c.OpenStack.SSHKey),
			openstack.WithSecurityGroup(c.OpenStack.SecurityGroup...),
			openstack.WithMetadata(c.OpenStack.Metadata),
			openstack.WithUserData(c.OpenStack.UserData),
			openstack.WithUserDataFile(c.OpenStack.UserDataFile),
		)
	case c.YandexCloud.Token != "" || c.YandexCloud.ServiceAccount != "":
		return yandexcloud.New(
			yandexcloud.WithToken(c.YandexCloud.Token),
			yandexcloud.WithServiceAccountJSON(c.YandexCloud.ServiceAccount),
			yandexcloud.WithFolderID(c.YandexCloud.FolderID),
			yandexcloud.WithSubnetID(c.YandexCloud.SubnetID),
			yandexcloud.WithZone(c.YandexCloud.Zone),
			yandexcloud.WithDiskSize(c.YandexCloud.DiskSize),
			yandexcloud.WithDiskType(c.YandexCloud.DiskType),
			yandexcloud.WithResourceCoreFraction(c.YandexCloud.ResourceCoreFraction),
			yandexcloud.WithPreemptible(c.YandexCloud.Preemptible),
			yandexcloud.WithPrivateIP(c.YandexCloud.PrivateIP),
			yandexcloud.WithResourceCores(c.YandexCloud.ResourceCores),
			yandexcloud.WithResourceMemory(c.YandexCloud.ResourceMemory),
			yandexcloud.WithPlatformID(c.YandexCloud.PlatformID),
			yandexcloud.WithImageFolderID(c.YandexCloud.ImageFolderID),
			yandexcloud.WithImageFamily(c.YandexCloud.ImageFamily),
			yandexcloud.WithDockerComposeConfig(c.YandexCloud.DockerComposeConfig),
			yandexcloud.WithSSHUserKeyPair(c.YandexCloud.SSHUserKeyPair),
			yandexcloud.WithSecurityGroups(c.YandexCloud.SecurityGroupIDs),
		)
	default:
		return nil, errors.New("missing provider configuration")
	}
}
