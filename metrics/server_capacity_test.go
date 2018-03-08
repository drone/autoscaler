// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"testing"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
)

func TestServerCapacity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// restore the default prometheus registerer
	// when the unit test is complete.
	snapshot := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = snapshot
		controller.Finish()
	}()

	// creates a blank registry
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry

	// x2 server count
	// x3 server capacity
	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1, Created: time.Now().Unix()},
		{Name: "server2", Capacity: 2, Created: time.Now().Unix()},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().ListState(gomock.Any(), autoscaler.StateRunning).Return(servers, nil)
	ServerCapacity(store)

	metrics, err := registry.Gather()
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := len(metrics), 1; want != got {
		t.Errorf("Expect registered metric")
		return
	}
	metric := metrics[0]
	if want, got := metric.GetName(), "drone_server_capacity"; want != got {
		t.Errorf("Expect metric name %s, got %s", want, got)
	}
	if want, got := metric.Metric[0].Gauge.GetValue(), float64(3); want != got {
		t.Errorf("Expect metric value %f, got %f", want, got)
	}
}
