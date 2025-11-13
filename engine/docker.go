// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/drone/autoscaler"
)

type dockerClient interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, container string, options container.StartOptions) error
	ContainerStop(ctx context.Context, container string, options container.StopOptions) error
	Ping(ctx context.Context) (types.Ping, error)
}

var _ dockerClient = (docker.APIClient)(nil)

// clientFunc defines a builder funciton used to build and return
// the docker client from a Server. This is primarily used for
// mock unit testing.
type clientFunc func(*autoscaler.Server) (dockerClient, io.Closer, error)

// newDockerClient returns a new Docker client configured for the
// Server host and certificate chain.
func newDockerClient(server *autoscaler.Server) (dockerClient, io.Closer, error) {
	tlsCert, err := tls.X509KeyPair(server.TLSCert, server.TLSKey)
	if err != nil {
		return nil, nil, err
	}
	tlsConfig := &tls.Config{
		ServerName:   server.Name,
		Certificates: []tls.Certificate{tlsCert},
	}
	tlsConfig.RootCAs = x509.NewCertPool()
	tlsConfig.RootCAs.AppendCertsFromPEM(server.CACert)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	dockerClient, err := docker.NewClientWithOpts(
		docker.WithAPIVersionNegotiation(),
		docker.WithHTTPClient(client),
		docker.WithHost(fmt.Sprintf("https://%s:2376", server.Address)),
	)
	return dockerClient, dockerClient, err
}
