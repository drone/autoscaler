// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package scaleway

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/drone/autoscaler"
	"github.com/rs/zerolog/log"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	p.init.Do(func() {
		p.setup(ctx)
	})

	api := instance.NewAPI(p.client)

	req := &instance.CreateServerRequest{
		Name:              opts.Name,
		DynamicIPRequired: scw.BoolPtr(p.dynamicIP),
		CommercialType:    p.size,
		Image:             p.image,
		Tags:              p.tags,
		SecurityGroup:     p.securityGroup,
	}

	logger := log.Ctx(ctx).With().
		Str("datacenter", string(p.zone)).
		Str("image", req.Image).
		Str("size", req.CommercialType).
		Str("name", req.Name).
		Logger()

	logger.Info().
		Msg("instance create")

	resp, err := api.CreateServer(req)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot create instance")
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = p.userdata.Execute(buf, &opts)
	if err != nil {
		return nil, err
	}

	err = api.SetServerUserData(&instance.SetServerUserDataRequest{
		Zone:     req.Zone,
		ServerID: resp.Server.ID,
		Key:      "cloud-init",
		Content:  buf,
	})

	if err != nil {
		return nil, err
	}

	logger.Debug().
		Str("name", req.Name).
		Msg("powering instance on")

	server, err := serverPowerAction(api, ctx, instance.ServerActionPoweron, resp.Server.ID)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("cannot power on instance")
		return nil, err
	}
	if server.State != instance.ServerStateRunning {
		return nil, errors.New("instance in invalid state: " + string(server.State))
	}

	logger.Info().
		Str("name", req.Name).
		Msg("instance created")

	ip := server.PublicIP
	if ip == nil {
		return nil, errors.New("server not assigned ip")
	}

	return &autoscaler.Instance{
		Provider: autoscaler.ProviderScaleway,
		ID:       server.ID,
		Name:     server.Name,
		Address:  ip.Address.String(),
		Region:   string(req.Zone),
		Image:    req.Image,
		Size:     req.CommercialType,
	}, nil
}

func serverPowerAction(api *instance.API, ctx context.Context, action instance.ServerAction, serverID string) (*instance.Server, error) {

	saReq := &instance.ServerActionRequest{
		ServerID: serverID,
		Action:   action,
	}

	gsReq := &instance.GetServerRequest{
		ServerID: serverID,
	}

	terminal := map[instance.ServerState]struct{}{
		instance.ServerStateStopped:        {},
		instance.ServerStateStoppedInPlace: {},
		instance.ServerStateLocked:         {},
		instance.ServerStateRunning:        {},
	}

	// Call to power the server on
	_, err := api.ServerAction(saReq, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var complete bool
	var server *instance.Server

	// Wait for context end, or poll every 3 seconds for
	// server status, until it is powered on
	for !complete {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(3 * time.Second):
			gsResp, err := api.GetServer(gsReq, scw.WithContext(ctx))
			if err != nil {
				return server, err
			}
			if _, complete = terminal[gsResp.Server.State]; complete {
				server = gsResp.Server
				break
			}
		}
	}
	if server == nil {
		return nil, errors.New("server is nil")
	}

	return server, nil
}
