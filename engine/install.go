// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/logger"
	"github.com/drone/autoscaler/metrics"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type installer struct {
	wg sync.WaitGroup

	os               string
	arch             string
	image            string
	secret           string
	volumes          []string
	ports            []string
	host             string
	proto            string
	envs             []string
	keepaliveTime    time.Duration
	keepaliveTimeout time.Duration
	runner           config.Runner
	labels           map[string]string

	checkInterval time.Duration
	checkDeadline time.Duration

	gcEnabled  bool
	gcDebug    bool
	gcImage    string
	gcIgnore   []string
	gcInterval time.Duration
	gcCache    string

	watchtowerEnabled  bool
	watchtowerImage    string
	watchtowerInterval int
	watchtowerTimeout  time.Duration

	servers autoscaler.ServerStore
	metrics metrics.Collector
	client  clientFunc
}

func (i *installer) Install(ctx context.Context) error {
	logger := logger.FromContext(ctx)

	servers, err := i.servers.ListState(ctx, autoscaler.StateCreated)
	if err != nil {
		return err
	}

	for _, server := range servers {
		server.State = autoscaler.StateStaging
		err = i.servers.Update(ctx, server)
		if err != nil {
			logger.WithError(err).
				WithField("server", server.Name).
				WithField("state", "staging").
				Errorln("failed to update server state")
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
	logger := logger.FromContext(ctx).
		WithField("ip", instance.Address).
		WithField("name", instance.Name)

	client, closer, err := i.client(instance)
	if closer != nil {
		defer closer.Close()
	}
	if err != nil {
		i.metrics.IncrServerInitError()
		logger.WithError(err).
			Errorln("cannot create docker client")
		return i.errorUpdate(ctx, instance, err)
	}

	logger.WithField("name", instance.Name).
		Debugln("check docker connectivity")

	timeout, cancel := context.WithTimeout(ctx, i.checkDeadline)
	defer cancel()

	start := time.Now()
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-timeout.Done():
			i.metrics.IncrServerInitError()
			logger.WithField("name", instance.Name).
				Debugln("connection timeout")

			return i.errorUpdate(ctx, instance, timeout.Err())
		case <-time.After(interval):
			interval = i.checkInterval

			logger.WithField("name", instance.Name).
				Debugln("connecting to docker")

			_, err := client.ContainerList(ctx, types.ContainerListOptions{})
			if err != nil {
				logger.
					WithField("error", err.Error()).
					WithField("name", instance.Name).
					Debugf("cannot connect, retry in %v", interval)
				continue
			}
			break poller
		}
	}

	// track elapsed time to establish a connection
	i.metrics.TrackServerInitTime(start)

	logger.WithField("image", i.image).
		Debugln("pull docker image")

	start = time.Now()
	rc, err := client.ImagePull(ctx, i.image, types.ImagePullOptions{})
	if err != nil {
		i.metrics.IncrServerSetupError()
		logger.WithError(err).
			WithField("image", i.image).
			Errorln("cannot pull docker image")
		return i.errorUpdate(ctx, instance, err)
	}
	io.Copy(ioutil.Discard, rc)
	rc.Close()

	logger.WithField("image", i.image).
		Debugln("create agent container")

	envs := append(i.envs,
		fmt.Sprintf("DRONE_RPC_HOST=%s", i.host),
		fmt.Sprintf("DRONE_RPC_PROTO=%s", i.proto),
		fmt.Sprintf("DRONE_RPC_SERVER=%s://%s", i.proto, i.host),
		fmt.Sprintf("DRONE_RPC_SECRET=%s", i.secret),
		fmt.Sprintf("DRONE_RUNNER_CAPACITY=%v", instance.Capacity),
		fmt.Sprintf("DRONE_RUNNER_NAME=%s", instance.Name),
	)

	if s := i.runner.Volumes; s != "" {
		envs = append(envs, fmt.Sprintf("DRONE_RUNNER_VOLUMES=%s", s))
	}
	if s := i.runner.Devices; s != "" {
		envs = append(envs, fmt.Sprintf("DRONE_RUNNER_DEVICES=%s", s))
	}
	if s := i.runner.EnvFile; s != "" {
		envs = append(envs, fmt.Sprintf("DRONE_RUNNER_ENV_FILE=%s", s))
	}
	if s := i.runner.Privileged; s != "" {
		envs = append(envs, fmt.Sprintf("DRONE_RUNNER_PRIVILEGED_IMAGES=%s", s))
	}

	if len(i.labels) > 0 {
		var stringLabels []string

		for key, val := range i.labels {
			stringLabels = append(stringLabels, fmt.Sprintf("%s:%s", key, val))
		}

		envs = append(envs,
			fmt.Sprintf("DRONE_RUNNER_LABELS=%s", strings.Join(stringLabels, ",")),
		)
	}

	var mounts []mount.Mount
	volumes := i.volumes
	switch i.os {
	case "windows":
		mounts = append(mounts, mount.Mount{
			Source: `\\.\pipe\docker_engine`,
			Target: `\\.\pipe\docker_engine`,
			Type:   mount.TypeNamedPipe,
		})
	default:
		volumes = append(volumes,
			"/var/run/docker.sock:/var/run/docker.sock",
		)

		// if memory serves me correctly, we need to explicitly
		// set this to nil to ensure the json representation
		// of this value is null. but I could be wrong in which
		// case this can be removed. ‾\_(ツ)_/‾
		mounts = nil
	}

	exposedPorts, portBindings, err := nat.ParsePortSpecs(i.ports)
	if err != nil {
		i.metrics.IncrServerInitError()
		logger.WithError(err).Errorln("could not create port binding")
		return i.errorUpdate(ctx, instance, err)
	}

	res, err := client.ContainerCreate(ctx,
		&container.Config{
			Image:        i.image,
			AttachStdout: true,
			AttachStderr: true,
			Env:          envs,
			Volumes:      toVol(volumes),
			ExposedPorts: exposedPorts,
			Labels: map[string]string{
				"com.centurylinklabs.watchtower.enable":      "true",
				"com.centurylinklabs.watchtower.stop-signal": "SIGHUP",
				"io.drone.agent.name":                        instance.Name,
				"io.drone.agent.zone":                        instance.Region,
				"io.drone.agent.size":                        instance.Size,
				"io.drone.agent.instance":                    instance.ID,
				"io.drone.agent.capacity":                    fmt.Sprint(instance.Capacity),
			},
		},
		&container.HostConfig{
			Binds:        volumes,
			Mounts:       mounts,
			PortBindings: portBindings,
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		}, nil, "agent")

	if err != nil {
		i.metrics.IncrServerSetupError()
		logger.WithField("image", i.image).
			WithError(err).
			Errorln("cannot create agent container")
		return i.errorUpdate(ctx, instance, err)
	}

	logger.WithField("image", i.image).
		Debugln("start the agent container")

	err = client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	if err != nil {
		i.metrics.IncrServerSetupError()
		logger.WithField("image", i.image).
			Debugln("cannot start the agent container")
		return i.errorUpdate(ctx, instance, err)
	}

	logger.WithField("image", i.image).
		Debugln("agent container started")

	if i.gcEnabled {
		logger.WithField("image", i.image).
			Debugln("setup the garbage collector")
		err = i.setupGarbageCollector(ctx, client)
		if err != nil {
			logger.WithError(err).
				WithField("image", i.image).
				Warnln("cannot setup the garbage collector")
		}
	}

	if i.watchtowerEnabled {
		logger.WithField("image", i.image).
			Debugln("setup watchtower")
		err = i.setupWatchtower(ctx, client)
		if err != nil {
			logger.WithError(err).
				WithField("image", i.image).
				Warnln("cannot setup watchtwoer")
		}
	}

	// track elapsed time to install software.
	i.metrics.TrackServerSetupTime(start)

	instance.State = autoscaler.StateRunning
	err = i.servers.Update(ctx, instance)
	if err != nil {
		i.metrics.IncrServerSetupError()
		logger.WithError(err).
			WithField("server", instance.Name).
			WithField("state", "running").
			Errorln("failed to update server state")
		return err
	}

	return nil
}

func (i *installer) setupWatchtower(ctx context.Context, client docker.APIClient) error {
	vols := []string{"/var/run/docker.sock:/var/run/docker.sock"}
	res, err := client.ContainerCreate(ctx,
		&container.Config{
			Image:        i.watchtowerImage,
			AttachStdout: true,
			AttachStderr: true,
			Volumes:      toVol(vols),
			Env: []string{
				fmt.Sprintf("WATCHTOWER_POLL_INTERVAL=%d", i.watchtowerInterval),
				fmt.Sprintf("WATCHTOWER_TIMEOUT=%s", i.watchtowerTimeout),
				fmt.Sprintf("WATCHTOWER_CLEANUP=true"),
				fmt.Sprintf("WATCHTOWER_LABEL_ENABLE=true"),
			},
		},
		&container.HostConfig{
			Binds: vols,
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		}, nil, "watchtower")
	if err != nil {
		return err
	}
	return client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
}

func (i *installer) setupGarbageCollector(ctx context.Context, client docker.APIClient) error {
	logger := logger.FromContext(ctx)
	vols := []string{"/var/run/docker.sock:/var/run/docker.sock"}
	envs := []string{
		fmt.Sprintf("GC_CACHE=%s", i.gcCache),
		fmt.Sprintf("GC_DEBUG=%v", i.gcDebug),
		fmt.Sprintf("GC_INTERVAL=%s", i.gcInterval),
	}
	if len(i.gcIgnore) > 0 {
		envs = append(envs,
			fmt.Sprintf("GC_IGNORE=%s", strings.Join(i.gcIgnore, ",")),
		)
	}

	logger.WithField("image", i.gcImage).
		Debugln("pull gc image")

	rc, err := client.ImagePull(ctx, i.gcImage, types.ImagePullOptions{})
	if err != nil {
		logger.WithError(err).
			WithField("image", i.gcImage).
			Errorln("cannot pull gc image")
		return err
	}
	io.Copy(ioutil.Discard, rc)
	rc.Close()

	res, err := client.ContainerCreate(ctx,
		&container.Config{
			Image:        i.gcImage,
			AttachStdout: true,
			AttachStderr: true,
			Volumes:      toVol(vols),
			Env:          envs,
			Labels: map[string]string{
				"com.centurylinklabs.watchtower.enable": "true",
			},
		},
		&container.HostConfig{
			Binds: vols,
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		}, nil, "drone-gc")
	if err != nil {
		return err
	}
	return client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
}

func (i *installer) errorUpdate(ctx context.Context, server *autoscaler.Server, err error) error {
	if err != nil {
		server.State = autoscaler.StateError
		server.Error = err.Error()
		xerr := i.servers.Update(ctx, server)
		if xerr != nil {
			logger.FromContext(ctx).
				WithError(xerr).
				WithField("server", server.Name).
				WithField("state", "error").
				Errorln("failed to update server state")
		}
	}
	return err
}

// helper function that converts a slice of volume paths to a set of
// unique volume names.
func toVol(paths []string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, path := range paths {
		parts, err := splitVolumeParts(path)
		if err != nil {
			continue
		}
		if len(parts) < 2 {
			continue
		}
		set[parts[1]] = struct{}{}
	}
	return set
}

// helper function that split volume path
func splitVolumeParts(volumeParts string) ([]string, error) {
	pattern := `^((?:[\w]\:)?[^\:]*)\:((?:[\w]\:)?[^\:]*)(?:\:([rwom]*))?`
	r, err := regexp.Compile(pattern)
	if err != nil {
		return []string{}, err
	}
	if r.MatchString(volumeParts) {
		results := r.FindStringSubmatch(volumeParts)[1:]
		cleanResults := []string{}
		for _, item := range results {
			if item != "" {
				cleanResults = append(cleanResults, item)
			}
		}
		return cleanResults, nil
	} else {
		return strings.Split(volumeParts, ":"), nil
	}
}
