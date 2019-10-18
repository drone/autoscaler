// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/mocks"

	"github.com/golang/mock/gomock"
)

func TestAllocate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockInstance := &autoscaler.Instance{}
	mockServers := []*autoscaler.Server{
		{State: autoscaler.StatePending},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StatePending).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(nil)
	store.EXPECT().Update(gomock.Any(), mockServers[0]).Return(nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockInstance, nil)

	a := allocator{servers: store, provider: provider}
	err := a.Allocate(mockctx)
	a.wg.Wait()

	if err != nil {
		t.Error(err)
	}
	if got, want := mockServers[0].State, autoscaler.StateCreated; got != want {
		t.Errorf("Want server state Created, got %v", got)
	}
}

func TestAllocate_ServerCreateError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockerr := errors.New("mock error")
	mockServers := []*autoscaler.Server{
		{State: autoscaler.StatePending},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StatePending).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(nil)
	store.EXPECT().Update(gomock.Any(), mockServers[0]).Return(nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, mockerr)

	a := allocator{servers: store, provider: provider}
	a.Allocate(mockctx)
	a.wg.Wait()

	if got, want := mockServers[0].State, autoscaler.StateError; got != want {
		t.Errorf("Want server state Error, got %v", got)
	}
}

func TestAllocate_ServerListError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockerr := errors.New("mock error")

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StatePending).Return(nil, mockerr)

	a := allocator{servers: store}
	if got, want := a.Allocate(mockctx), mockerr; got != want {
		t.Errorf("Want error getting server list")
	}
}

func TestAllocate_ServerUpdateError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockctx := context.Background()
	mockerr := errors.New("mock error")
	mockServers := []*autoscaler.Server{
		{State: autoscaler.StatePending},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(mockctx, autoscaler.StatePending).Return(mockServers, nil)
	store.EXPECT().Update(mockctx, mockServers[0]).Return(mockerr)

	a := allocator{servers: store}
	if got, want := a.Allocate(mockctx), mockerr; got != want {
		t.Errorf("Want error updating server")
	}
	if got, want := mockServers[0].State, autoscaler.StateCreating; got != want {
		t.Errorf("Want server state Staging, got %v", got)
	}
}
