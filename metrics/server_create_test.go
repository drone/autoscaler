// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"errors"
	"testing"

	"github.com/kr/pretty"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
)

func TestServerCreate(t *testing.T) {
	controller := gomock.NewController(t)

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

	opts := &autoscaler.ServerOpts{}
	server := &autoscaler.Server{Name: "server1", Capacity: 1}

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Create(gomock.Any(), opts).Times(3).Return(server, nil)
	provider.EXPECT().Create(gomock.Any(), opts).Return(nil, errors.New("error"))

	providerInst := ServerCreate(provider)
	for i := 0; i < 3; i++ {
		result, err := providerInst.Create(noContext, opts)
		if err != nil {
			t.Error(err)
		}
		if got, want := result, server; got != want {
			t.Errorf("Expect server returned from provider, CALL %d", i)
			pretty.Ldiff(t, got, want)
		}
	}
	_, err := providerInst.Create(noContext, opts)
	if err == nil {
		t.Errorf("Expect error returned from provider")
	}

	metrics, err := registry.Gather()
	if err != nil {
		t.Error(err)
		return
	}
	if want, got := len(metrics), 2; want != got {
		t.Errorf("Expect registered metric")
		return
	}
	if got, want := metrics[0].GetName(), "drone_servers_created"; want != got {
		t.Errorf("Expect metric name %s, got %s", want, got)
	}
	if got, want := metrics[0].Metric[0].Counter.GetValue(), float64(3); want != got {
		t.Errorf("Expect metric value %f, got %f", want, got)
	}
	if got, want := metrics[1].GetName(), "drone_servers_created_err"; want != got {
		t.Errorf("Expect metric name %s, got %s", want, got)
	}
	if got, want := metrics[1].Metric[0].Counter.GetValue(), float64(1); want != got {
		t.Errorf("Expect metric value %f, got %f", want, got)
	}
}
