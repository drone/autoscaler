// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"

	"docker.io/go-docker"
	"docker.io/go-docker/api"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
)

type installer struct {
	wg sync.WaitGroup

	image  string
	secret string
	server string

	servers autoscaler.ServerStore
}

func (i *installer) Install(ctx context.Context) error {
	logger := log.Ctx(ctx)

	servers, err := i.servers.ListState(ctx, autoscaler.StateCreated)
	if err != nil {
		return err
	}

	for _, server := range servers {
		server.State = autoscaler.StateStaging
		err = i.servers.Update(ctx, server)
		if err != nil {
			logger.Error().
				Err(err).
				Str("server", server.Name).
				Str("state", "staging").
				Msg("failed to update server state")
			return err
		}

		i.wg.Add(1)
		go func(server *autoscaler.Server) {
			i.install(ctx, server)
			i.wg.Done()
		}(server)
	}
	return nil
}

func (i *installer) install(ctx context.Context, instance *autoscaler.Server) error {
	logger := log.Ctx(ctx).With().
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Logger()

	client, err := newDockerClient(instance)
	if err != nil {
		logger.Error().Err(err).
			Msg("cannot create docker client")
		return i.errorUpdate(ctx, instance, err)
	}

	logger.Debug().
		Str("name", instance.Name).
		Msg("check docker connectivity")

	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logger.Debug().
				Str("name", instance.Name).
				Msg("docker timeout")

			return i.errorUpdate(ctx, instance, ctx.Err())
		case <-time.After(interval):
			interval = time.Minute

			logger.Debug().
				Str("name", instance.Name).
				Msg("docker ping")

			_, err := client.ContainerList(ctx, types.ContainerListOptions{})
			if err != nil {
				logger.Debug().
					Str("error", err.Error()).
					Str("name", instance.Name).
					Msg("cannot ping docker")
				continue
			}
			break poller
		}
	}

	logger.Debug().
		Str("image", i.image).
		Msg("pull docker image")

	rc, err := client.ImagePull(ctx, i.image, types.ImagePullOptions{})
	if err != nil {
		logger.Error().Err(err).
			Str("image", i.image).
			Msg("cannot pull docker image")
		return i.errorUpdate(ctx, instance, err)
	}
	io.Copy(ioutil.Discard, rc)
	rc.Close()

	logger.Debug().
		Str("image", i.image).
		Msg("create agent container")

	res, err := client.ContainerCreate(ctx,
		&container.Config{
			Image:        i.image,
			AttachStdout: true,
			AttachStderr: true,
			Env: []string{
				fmt.Sprintf("DRONE_SECRET=%s", i.secret),
				fmt.Sprintf("DRONE_SERVER=%s", i.server),
				fmt.Sprintf("DRONE_MAX_PROCS=%v", instance.Capacity),
				fmt.Sprintf("DRONE_HOSTNAME=%s", instance.Name),
			},
			Volumes: map[string]struct{}{
				"/var/run/docker.sock": struct{}{},
			},
		},
		&container.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		}, nil, "agent")

	if err != nil {
		logger.Error().Err(err).
			Str("image", i.image).
			Msg("cannot create agent container")
		return i.errorUpdate(ctx, instance, err)
	}

	logger.Debug().
		Str("image", i.image).
		Msg("start the agent container")

	err = client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	if err != nil {
		logger.Debug().
			Str("image", i.image).
			Msg("cannot start the agent container")
		return i.errorUpdate(ctx, instance, err)
	}

	logger.Debug().
		Str("image", i.image).
		Msg("agent container started")

	instance.State = autoscaler.StateRunning
	return i.servers.Update(ctx, instance)
}

func (i *installer) errorUpdate(ctx context.Context, server *autoscaler.Server, err error) error {
	if err != nil {
		server.State = autoscaler.StateError
		i.servers.Update(ctx, server)
	}
	return err
}

// helper function returns a new docker client.
func newDockerClient(server *autoscaler.Server) (docker.APIClient, error) {
	tlsCert, err := tls.X509KeyPair(server.TLSCert, server.TLSKey)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{tlsCert}}

	// https://github.com/moby/moby/issues/8943
	// https://github.com/prometheus/prometheus/pull/1658
	tlsConfig.ServerName = server.Name
	// tlsConfig.InsecureSkipVerify = true

	tlsConfig.RootCAs = x509.NewCertPool()
	tlsConfig.RootCAs.AppendCertsFromPEM(server.CACert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return docker.NewClient("https://"+server.Address+":2376", api.DefaultVersion, client, nil)
}
