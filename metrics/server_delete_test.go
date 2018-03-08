// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"errors"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
)

func TestServerDelete(t *testing.T) {
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

	instance := &autoscaler.Instance{Name: "server1"}

	provider := mocks.NewMockProvider(controller)
	provider.EXPECT().Destroy(noContext, instance).Times(3).Return(nil)
	provider.EXPECT().Destroy(noContext, instance).Return(errors.New("error"))

	providerInst := ServerDelete(provider)
	for i := 0; i < 3; i++ {
		err := providerInst.Destroy(noContext, instance)
		if err != nil {
			t.Error(err)
		}
	}
	err := providerInst.Destroy(noContext, instance)
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
	if got, want := metrics[0].GetName(), "drone_servers_deleted"; want != got {
		t.Errorf("Expect metric name %s, got %s", want, got)
	}
	if got, want := metrics[0].Metric[0].Counter.GetValue(), float64(3); want != got {
		t.Errorf("Expect metric value %f, got %f", want, got)
	}
	if got, want := metrics[1].GetName(), "drone_servers_deleted_err"; want != got {
		t.Errorf("Expect metric name %s, got %s", want, got)
	}
	if got, want := metrics[1].Metric[0].Counter.GetValue(), float64(1); want != got {
		t.Errorf("Expect metric value %f, got %f", want, got)
	}
}
