// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/mocks"

	docker "github.com/docker/docker/client"
	"github.com/golang/mock/gomock"
)

func TestCollect(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockServers := []*autoscaler.Server{
		{
			ID:      "i-1234567890abcdef0",
			Address: "1.2.3.4",
			State:   autoscaler.StateShutdown,
		},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().ContainerStop(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StateShutdown).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(nil)
	store.EXPECT().Update(gomock.Any(), mockServers[0]).Return(nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(nil)

	c := collector{
		servers:  store,
		provider: provider,
		client: func(*autoscaler.Server) (docker.APIClient, io.Closer, error) {
			return client, nil, nil
		},
	}
	err := c.Collect(mockctx)
	c.wg.Wait()

	if err != nil {
		t.Error(err)
	}
	if got, want := mockServers[0].State, autoscaler.StateStopped; got != want {
		t.Errorf("Want server state Stopped, got %v", got)
	}
}

func TestCollect_DockerStopError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockerr := errors.New("oh no")
	mockctx := context.Background()
	mockServers := []*autoscaler.Server{
		{
			ID:      "i-1234567890abcdef0",
			Address: "1.2.3.4",
			State:   autoscaler.StateShutdown,
		},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().ContainerStop(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockerr)

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StateShutdown).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(nil)
	store.EXPECT().Update(gomock.Any(), mockServers[0]).Return(nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(nil)

	c := collector{
		servers:  store,
		provider: provider,
		client: func(*autoscaler.Server) (docker.APIClient, io.Closer, error) {
			return client, nil, nil
		},
	}
	err := c.Collect(mockctx)
	c.wg.Wait()

	if err != nil {
		t.Error(err)
	}
	if got, want := mockServers[0].State, autoscaler.StateStopped; got != want {
		t.Errorf("Want server state Stopped, got %v", got)
	}
}

func TestCollect_ServerDestroyError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockerr := errors.New("mock error")
	mockServers := []*autoscaler.Server{
		{
			ID:      "i-1234567890abcdef0",
			Name:    "agent-807jVFwj",
			Address: "1.2.3.4",
			State:   autoscaler.StateShutdown,
		},
	}

	client := mocks.NewMockAPIClient(controller)
	client.EXPECT().ContainerStop(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StateShutdown).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(nil)
	store.EXPECT().Update(gomock.Any(), mockServers[0]).Return(nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(mockerr)

	c := collector{
		servers:  store,
		provider: provider,
		client: func(*autoscaler.Server) (docker.APIClient, io.Closer, error) {
			return client, nil, nil
		},
	}
	c.Collect(mockctx)
	c.wg.Wait()

	if got, want := mockServers[0].State, autoscaler.StateError; got != want {
		t.Errorf("Want server state Error, got %v", got)
	}
}

func TestCollect_ServerListError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockerr := errors.New("mock error")

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StateShutdown).Return(nil, mockerr)

	c := collector{servers: store}
	if got, want := c.Collect(mockctx), mockerr; got != want {
		t.Errorf("Want error getting server list")
	}
}

func TestCollect_ServerUpdateError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockerr := errors.New("mock error")
	mockServers := []*autoscaler.Server{
		{
			ID:      "i-1234567890abcdef0",
			Address: "1.2.3.4",
			State:   autoscaler.StateShutdown,
		},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StateShutdown).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(mockerr)

	c := collector{servers: store}
	if got, want := c.Collect(mockctx), mockerr; got != want {
		t.Errorf("Want error updating server")
	}
	if got, want := mockServers[0].State, autoscaler.StateStopping; got != want {
		t.Errorf("Want server state Stopping, got %v", got)
	}
}

func TestCollect_ServerNeverProvisioned(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockServer := &autoscaler.Server{
		ID:    "",
		State: autoscaler.StateShutdown,
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), mockServer).Return(nil).Times(1)

	c := collector{servers: store}
	if err := c.collect(mockctx, mockServer); err != nil {
		t.Error(err)
	}
	if got, want := mockServer.State, autoscaler.StateStopped; got != want {
		t.Errorf("Want server state Stopping, got %v", got)
	}
}
