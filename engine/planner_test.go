// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

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
func TestPlan_Noop(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 2, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 2, State: autoscaler.StateRunning},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return([]*drone.Stage{
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
	}, nil)

	p := planner{
		cap:     2,
		min:     2,
		max:     10,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if that no servers are
// destroyed if there is excess capacity and the
// the server count <= the min pool size.
func TestPlan_MinBufferCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Stage{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		buffer:  1,
		min:     2,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if the server capacity minus buffer is
// less than the pending count, and the server capacity is
// >= the pool maximum, no actions are taken.
func TestPlan_MaxBufferCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x4 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server3", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server4", Capacity: 1, State: autoscaler.StateRunning},
	}

	// x3 running builds
	// x1 pending builds
	builds := []*drone.Stage{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     1,
		buffer:  2,
		min:     2,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if the server capacity minus buffer is
// less than the pending count, and the server capacity is
// < the pool maximum, additional servers are provisioned.
func TestPlan_MoreBufferCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x4 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 2, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 2, State: autoscaler.StateRunning},
	}

	// x2 running builds
	// x1 pending builds
	builds := []*drone.Stage{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		buffer:  2,
		min:     2,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if the server capacity is
// < than the pending count, and the server capacity is
// >= the pool maximum, no actions are taken.
func TestPlan_MaxCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x4 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server3", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server4", Capacity: 1, State: autoscaler.StateRunning},
	}

	// x4 running builds
	// x3 pending builds
	builds := []*drone.Stage{
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
	client.EXPECT().Queue().Return(builds, nil)

	config := config.Config{}
	config.Pool.Min = 2
	config.Pool.Max = 4
	config.Agent.Concurrency = 2

	p := planner{
		cap:     2,
		min:     2,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if the server capacity is
// less than the pending count, and the server capacity is
// < the pool maximum, additional servers are provisioned.
func TestPlan_MoreCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning},
	}

	// x2 running builds
	// x3 pending builds
	builds := []*drone.Stage{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending}, // ignore, would exceed max pool size
		{Status: drone.StatusPending}, // ignore, would exceed max pool size
		{Status: drone.StatusPending}, // ignore, would exceed max pool size
		{Status: drone.StatusPending}, // ignore, would exceed max pool size
		{Status: drone.StatusPending}, // ignore, would exceed max pool size
		{Status: drone.StatusPending}, // ignore, would exceed max pool size
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	store.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     2,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that if that no servers are
// destroyed if there is excess capacity and the
// the server count <= the min pool size.
func TestPlan_MinPool(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Stage{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     2,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that no servers are
// destroyed if no idle servers exist.
func TestPlan_NoIdle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 2, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 2, State: autoscaler.StateRunning},
	}

	// x2 running builds
	// x0 pending builds
	builds := []*drone.Stage{
		{Status: drone.StatusRunning, Machine: "server1"},
		{Status: drone.StatusRunning, Machine: "server2"},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     1,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.Background())
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
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning, Created: time.Now().Unix()},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning, Created: time.Now().Unix()},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Stage{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers, nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     1,
		max:     4,
		ttu:     time.Hour,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// This test verifies that idle servers are not
// garbage collected until the min-idle is reached.
func TestScale_MinIdle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x2 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, State: autoscaler.StateRunning, Created: time.Now().Unix(), LastBusy: time.Now().Add(time.Minute * 10 * -1).Unix()},
		{Name: "server2", Capacity: 1, State: autoscaler.StateRunning, Created: time.Now().Unix(), LastBusy: time.Now().Add(time.Minute * 60 * -1).Unix()},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Stage{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers, nil)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers, nil)
	// we should expect a call to shut down only server2 since the last busy
	// time of 1 hour ago satisfies the MinIdle (tti) of 30 minutes.
	store.EXPECT().Update(gomock.Any(), servers[1]).Return(nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)
	client.EXPECT().Queue().Return(builds, nil)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     1,
		max:     4,
		ttu:     0,
		tti:     time.Minute * 30,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

func TestPlan_ShutdownIdle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x3 capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 2, Created: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 2, Created: 2, State: autoscaler.StateRunning},
		{Name: "server3", Capacity: 2, Created: 3, State: autoscaler.StateRunning},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Stage{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers, nil)
	store.EXPECT().Update(gomock.Any(), servers[2]).Return(nil)
	store.EXPECT().Update(gomock.Any(), servers[1]).Return(nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     1,
		max:     4,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

func TestPlan_ExcludePendingWhenTerminating(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	servers := []*autoscaler.Server{
		// x3 capacity
		{Name: "server1", Capacity: 2, Created: 1, State: autoscaler.StateRunning},
		{Name: "server2", Capacity: 2, Created: 2, State: autoscaler.StateRunning},
		{Name: "server3", Capacity: 2, Created: 3, State: autoscaler.StateRunning},

		// x3 pending / staging / starting
		{Name: "server4", Capacity: 2, Created: 4, State: autoscaler.StateCreating},
		{Name: "server5", Capacity: 2, Created: 5, State: autoscaler.StateCreated},
		{Name: "server6", Capacity: 2, Created: 6, State: autoscaler.StateStaging},
	}

	// x0 running builds
	// x0 pending builds
	builds := []*drone.Stage{}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers[:3], nil)

	client := mocks.NewMockClient(controller)
	client.EXPECT().Queue().Return(builds, nil)

	p := planner{
		cap:     2,
		min:     3,
		max:     10,
		client:  client,
		servers: store,
	}

	err := p.Plan(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

// func TestListBusy(t *testing.T) {
// 	controller := gomock.NewController(t)
// 	defer controller.Finish()

// 	client := mocks.NewMockClient(controller)
// 	client.EXPECT().Build("octocat", "hello-world", 1).Return(&drone.Build{
// 		Procs: []*drone.Proc{
// 			{PID: 1, Machine: "machine1"},
// 			{PID: 2, Machine: "machine2"},
// 		},
// 	}, nil)
// 	client.EXPECT().BuildQueue().Return([]*drone.Activity{
// 		{Status: drone.StatusPending},
// 		{Status: drone.StatusRunning, Owner: "octocat", Name: "hello-world", Number: 1},
// 	}, nil)

// 	scaler := Scaler{Client: client}
// 	busy, err := scaler.listBusy(context.TODO())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if got, want := len(busy), 2; got != want {
// 		t.Errorf("Want busy server count %d, got %d", want, got)
// 	}
// 	if _, ok := busy["machine1"]; !ok {
// 		t.Errorf("Expected server not in busy list")
// 	}
// 	if _, ok := busy["machine2"]; !ok {
// 		t.Errorf("Expected server not in busy list")
// 	}
// }

// func TestCapacity(t *testing.T) {
// 	controller := gomock.NewController(t)
// 	defer controller.Finish()

// 	servers := []*autoscaler.Server{
// 		{Name: "server1", Capacity: 4},
// 		{Name: "server2", Capacity: 3},
// 		{Name: "server3", Capacity: 2},
// 		{Name: "server4", Capacity: 1},
// 	}

// 	store := mocks.NewMockServerStore(controller)
// 	store.EXPECT().List(gomock.Any()).Return(servers, nil)

// 	scaler := Scaler{Servers: store}
// 	capacity, count, err := scaler.capacity(context.TODO())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if got, want := capacity, 10; got != want {
// 		t.Errorf("Want capacity count %d, got %d", want, got)
// 	}
// 	if got, want := count, 4; got != want {
// 		t.Errorf("Want server count %d, got %d", want, got)
// 	}
// }

// func TestCount(t *testing.T) {
// 	controller := gomock.NewController(t)
// 	defer controller.Finish()

// 	client := mocks.NewMockClient(controller)
// 	client.EXPECT().BuildQueue().Return([]*drone.Activity{
// 		{Status: drone.StatusPending},
// 		{Status: drone.StatusPending},
// 		{Status: drone.StatusPending},
// 		{Status: drone.StatusRunning},
// 		{Status: drone.StatusRunning},
// 	}, nil)

// 	scaler := Scaler{Client: client}
// 	pending, running, err := scaler.count(context.TODO())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if got, want := pending, 3; got != want {
// 		t.Errorf("Want pending count %d, got %d", want, got)
// 	}
// 	if got, want := running, 2; got != want {
// 		t.Errorf("Want running count %d, got %d", want, got)
// 	}
// }

func TestMatch(t *testing.T) {
	tests := []struct {
		match   bool
		os      string
		arch    string
		version string
		kernel  string
		labels  map[string]string
		stage   *drone.Stage
	}{
		{
			match: true,
			os:    "linux",
			arch:  "amd64",
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
			},
		},
		{
			match: false,
			os:    "linux",
			arch:  "amd64",
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "arm",
			},
		},
		{
			match: false,
			os:    "linux",
			arch:  "amd64",
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
				Labels: map[string]string{
					"region": "us-west-2",
				},
			},
		},
		{
			match: false,
			os:    "linux",
			arch:  "amd64",
			labels: map[string]string{
				"region": "us-west-2",
			},
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
			},
		},
		{
			match: true,
			os:    "linux",
			arch:  "amd64",
			labels: map[string]string{
				"region": "us-west-2",
			},
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
				Labels: map[string]string{
					"region": "us-west-2",
				},
			},
		},
		{
			match: true,
			os:    "linux",
			arch:  "amd64",
			labels: map[string]string{
				"region": "us-west-2",
				"mem":    "high",
			},
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
				Labels: map[string]string{
					"region": "us-west-2",
					"mem":    "high",
				},
			},
		},
		{
			match: false,
			os:    "linux",
			arch:  "amd64",
			labels: map[string]string{
				"region": "us-east-2",
			},
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
				Labels: map[string]string{
					"region": "us-west-2",
				},
			},
		},
		{
			match: false,
			os:    "linux",
			arch:  "amd64",
			labels: map[string]string{
				"region": "us-east-2",
				"mem":    "high",
			},
			stage: &drone.Stage{
				OS:   "linux",
				Arch: "amd64",
				Labels: map[string]string{
					"region": "us-west-2",
				},
			},
		},
	}
	for _, test := range tests {
		p := &planner{
			os:      test.os,
			arch:    test.arch,
			version: test.version,
			kernel:  test.kernel,
			labels:  test.labels,
		}
		if p.match(test.stage) != test.match {
			t.Fail()
			return
		}
	}
}
