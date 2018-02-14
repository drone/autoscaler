// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scaler

import (
	"context"
	"testing"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/mocks"
	"github.com/drone/drone-go/drone"
	"github.com/golang/mock/gomock"
)

// This test verifies that if the server capacity is
// >= the pending count, and the server capacity is
// <= the pool minimum size, no actions are taken.
func TestScale_Noop(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 2},
		{Name: "server2", Capacity: 2},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return([]*drone.Activity{
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
	}, nil)

	config := config.Config{}
	config.Pool.Min = 2
	scaler := Scaler{
		Client:  client,
		Servers: store,
		Config:  config,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if the server capacity is
// < than the pending count, and the server capacity is
// >= the pool maximum, no actions are taken.
func TestScale_MaxCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x4 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1},
		{Name: "server2", Capacity: 1},
		{Name: "server3", Capacity: 1},
		{Name: "server4", Capacity: 1},
	}

	// x4 running builds
	// x3 pending builds
	builds := []*drone.Activity{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(builds, nil)

	config := config.Config{}
	config.Pool.Min = 2
	config.Pool.Max = 4
	config.Agent.Concurrency = 2

	scaler := Scaler{
		Client:  client,
		Servers: store,
		Config:  config,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if the server capacity is
// less than the pending count, and the server capacity is
// < the pool maximum, additional servers are provisioned.
func TestScale_MoreCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1},
		{Name: "server2", Capacity: 1},
	}

	// x2 running builds
	// x3 pending builds
	builds := []*drone.Activity{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
	}

	// x2 mock servers provisioned
	server1 := &autoscaler.Server{Name: "i-5203422c"}
	server2 := &autoscaler.Server{Name: "i-4421485g"}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().Create(gomock.Any(), server1).Return(nil)
	store.EXPECT().Create(gomock.Any(), server2).Return(nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(builds, nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Create(gomock.Any(), gomock.Any()).Return(server1, nil)
	provider.EXPECT().Create(gomock.Any(), gomock.Any()).Return(server2, nil)

	config := config.Config{}
	config.Pool.Min = 2
	config.Pool.Max = 4
	config.Agent.Concurrency = 2

	scaler := Scaler{
		Client:   client,
		Servers:  store,
		Config:   config,
		Provider: provider,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if that no servers are
// destroyed if there is excess capacity and the
// the server count <= the min pool size.
func TestScale_MinPool(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1},
		{Name: "server2", Capacity: 1},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Activity{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(builds, nil)

	provider := mocks.NewMockProvider(controller)

	config := config.Config{}
	config.Pool.Min = 2
	config.Pool.Max = 4
	config.Agent.Concurrency = 2

	scaler := Scaler{
		Client:   client,
		Servers:  store,
		Config:   config,
		Provider: provider,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if that no servers are
// destroyed if no idle servers exist.
func TestScale_NoIdle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x3 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1},
		{Name: "server2", Capacity: 1},
		{Name: "server3", Capacity: 1},
	}

	// x3 running builds
	// x0 pending builds
	builds := []*drone.Activity{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(builds, nil)
	client.EXPECT().BuildQueue().Return(builds, nil)
	client.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(&drone.Build{Procs: []*drone.Proc{{Machine: "server1"}}}, nil)
	client.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(&drone.Build{Procs: []*drone.Proc{{Machine: "server2"}}}, nil)
	client.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(&drone.Build{Procs: []*drone.Proc{{Machine: "server3"}}}, nil)

	provider := mocks.NewMockProvider(controller)

	config := config.Config{}
	config.Pool.Min = 2
	config.Pool.Max = 4
	config.Agent.Concurrency = 2

	scaler := Scaler{
		Client:   client,
		Servers:  store,
		Config:   config,
		Provider: provider,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that idle servers are not
// garbage collected until the min-age is reached.
func TestScale_MinAge(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, Created: time.Now().Unix()},
		{Name: "server2", Capacity: 1, Created: time.Now().Unix()},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Activity{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(builds, nil)
	client.EXPECT().BuildQueue().Return(builds, nil)

	provider := mocks.NewMockProvider(controller)

	config := config.Config{}
	config.Pool.Min = 1
	config.Pool.Max = 4
	config.Pool.MinAge = time.Hour
	config.Agent.Concurrency = 2

	scaler := Scaler{
		Client:   client,
		Servers:  store,
		Config:   config,
		Provider: provider,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that idle servers are
// garbage collected while preserving the minimum
// pool size.
func TestScale_DestroyIdle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x3 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, Created: 1},
		{Name: "server2", Capacity: 1, Created: 2},
		{Name: "server3", Capacity: 1, Created: 3},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Activity{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().Delete(gomock.Any(), servers[1]).Return(nil)
	store.EXPECT().Delete(gomock.Any(), servers[2]).Return(nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(builds, nil)
	client.EXPECT().BuildQueue().Return(builds, nil)

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Destroy(gomock.Any(), servers[1]).Return(nil)
	provider.EXPECT().Destroy(gomock.Any(), servers[2]).Return(nil)

	config := config.Config{}
	config.Pool.Min = 1
	config.Pool.Max = 4
	config.Agent.Concurrency = 2

	scaler := Scaler{
		Client:   client,
		Servers:  store,
		Config:   config,
		Provider: provider,
	}

	err := scaler.Scale(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

func TestListBusy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := mocks.NewMockClient(controller)
	client.EXPECT().Build("octocat", "hello-world", 1).Return(&drone.Build{
		Procs: []*drone.Proc{
			{PID: 1, Machine: "machine1"},
			{PID: 2, Machine: "machine2"},
		},
	}, nil)
	client.EXPECT().BuildQueue().Return([]*drone.Activity{
		{Status: drone.StatusPending},
		{Status: drone.StatusRunning, Owner: "octocat", Name: "hello-world", Number: 1},
	}, nil)

	scaler := Scaler{Client: client}
	busy, err := scaler.listBusy(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(busy), 2; got != want {
		t.Errorf("Want busy server count %d, got %d", want, got)
	}
	if _, ok := busy["machine1"]; !ok {
		t.Errorf("Expected server not in busy list")
	}
	if _, ok := busy["machine2"]; !ok {
		t.Errorf("Expected server not in busy list")
	}
}

func TestCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 4},
		{Name: "server2", Capacity: 3},
		{Name: "server3", Capacity: 2},
		{Name: "server4", Capacity: 1},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	scaler := Scaler{Servers: store}
	capacity, count, err := scaler.capacity(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := capacity, 10; got != want {
		t.Errorf("Want capacity count %d, got %d", want, got)
	}
	if got, want := count, 4; got != want {
		t.Errorf("Want server count %d, got %d", want, got)
	}
}

func TestCount(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return([]*drone.Activity{
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
	}, nil)

	scaler := Scaler{Client: client}
	pending, running, err := scaler.count(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := pending, 3; got != want {
		t.Errorf("Want pending count %d, got %d", want, got)
	}
	if got, want := running, 2; got != want {
		t.Errorf("Want running count %d, got %d", want, got)
	}
}

func TestRequiredCapacity(t *testing.T) {
	tests := []struct {
		pending,
		available,
		concurrency,
		want int
	}{
		{0, 2, 2, 0}, // no pending builds
		{2, 2, 2, 0}, // use 2 of 2 existing
		{1, 2, 2, 0}, // use 1 of 2 existing
		{4, 2, 2, 1}, // want 2 servers
		{4, 2, 1, 2}, // want 2 servers
		{5, 2, 2, 2}, // want 2 servers (round-up)
	}
	for _, test := range tests {
		capacity := requiredCapacity(
			test.pending,
			test.available,
			test.concurrency,
		)
		if got, want := capacity, test.want; got != want {
			t.Errorf("Got capacity %d, want %d", got, want)
		}
	}
}
