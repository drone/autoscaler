// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/drivers/amazon"
	"github.com/drone/autoscaler/drivers/digitalocean"
	"github.com/drone/autoscaler/drivers/google"
	"github.com/drone/autoscaler/drivers/hetznercloud"
	"github.com/drone/autoscaler/drivers/openstack"
	"github.com/drone/autoscaler/drivers/packet"
	"github.com/drone/autoscaler/drivers/scaleway"
	"github.com/drone/autoscaler/engine"
	"github.com/drone/autoscaler/metrics"
	"github.com/drone/autoscaler/server"
	"github.com/drone/autoscaler/slack"
	"github.com/drone/autoscaler/store"
	"github.com/drone/drone-go/drone"
	"github.com/drone/signal"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).
			Msg("Invalid or missing hosting provider")
	}

	// instruments the provider with prometheus metrics.
	provider = metrics.ServerCreate(provider)
	provider = metrics.ServerDelete(provider)

	db, err := store.Connect(conf.Database.Driver, conf.Database.Datasource)
	if err != nil {
		log.Fatal().Err(err).
			Msg("Cannot establish database connection")
	}

	servers := store.NewServerStore(db)
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
	)

	r := chi.NewRouter()
	r.Use(hlog.NewHandler(log.Logger))
	r.Use(hlog.RemoteAddrHandler("ip"))
	r.Use(hlog.URLHandler("path"))
	r.Use(hlog.MethodHandler("method"))
	r.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	r.Route(conf.HTTP.Root, func(root chi.Router) {
		root.Get("/metrics", server.HandleMetrics(conf.Prometheus.AuthToken))
		root.Get("/version", server.HandleVersion(source, version, commit))
		root.Get("/healthz", server.HandleHealthz())
		root.Get("/varz", server.HandleVarz(enginex))
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

	ctx := log.Logger.WithContext(context.Background())
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
		log.Fatal().Err(err).Msg("Program terminated")
	}
}

// helper funciton configures the http server.
func setupServer(c config.Config) *http.Server {
	return &http.Server{
		Addr: c.HTTP.Port,
	}
}

// helper funciton configures the logging.
func setupLogging(c config.Config) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if c.Logs.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if c.Logs.Pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !c.Logs.Color,
			},
		)
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
			google.WithProject(c.Google.Project),
			google.WithTags(c.Google.Tags...),
			google.WithUserData(c.Google.UserData),
			google.WithUserDataFile(c.Google.UserDataFile),
			google.WithZone(c.Google.Zone),
		)
	case c.DigitalOcean.Token != "":
		return digitalocean.New(
			digitalocean.WithSSHKey(c.DigitalOcean.SSHKey),
			digitalocean.WithImage(c.DigitalOcean.Image),
			digitalocean.WithRegion(c.DigitalOcean.Region),
			digitalocean.WithSize(c.DigitalOcean.Size),
			digitalocean.WithUserDataFile(c.DigitalOcean.UserDataFile),
			digitalocean.WithUserData(c.DigitalOcean.UserData),
			digitalocean.WithToken(c.DigitalOcean.Token),
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
			amazon.WithSubnet(c.Amazon.SubnetID),
			amazon.WithTags(c.Amazon.Tags),
			amazon.WithUserData(c.Amazon.UserData),
			amazon.WithUserDataFile(c.Amazon.UserDataFile),
			amazon.WithVolumeSize(c.Amazon.VolumeSize),
			amazon.WithVolumeType(c.Amazon.VolumeType),
			amazon.WithIamProfileArn(c.Amazon.IamProfileArn),
			amazon.WithMarketType(c.Amazon.MarketType),
		), nil
	case os.Getenv("OS_USERNAME") != "":
		return openstack.New(
			openstack.WithImage(c.OpenStack.Image),
			openstack.WithRegion(c.OpenStack.Region),
			openstack.WithFlavor(c.OpenStack.Flavor),
			openstack.WithFloatingIpPool(c.OpenStack.Pool),
			openstack.WithSSHKey(c.OpenStack.SSHKey),
			openstack.WithSecurityGroup(c.OpenStack.SecurityGroup...),
			openstack.WithMetadata(c.OpenStack.Metadata),
			openstack.WithUserData(c.OpenStack.UserData),
			openstack.WithUserDataFile(c.OpenStack.UserDataFile),
		)
	default:
		return nil, errors.New("missing provider configuration")
	}
}
