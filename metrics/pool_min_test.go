// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"testing"

	"github.com/drone/autoscaler/config"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMinPool(t *testing.T) {
	// restore the default prometheus registerer
	// when the unit test is complete.
	snapshot := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = snapshot
	}()

	// creates a blank registry
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry

	conf := config.Config{}
	conf.Pool.Min = 5
	MinPool(conf)

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
	if want, got := metric.GetName(), "drone_server_min_pool"; want != got {
		t.Errorf("Expect metric name %s, got %s", want, got)
	}
	if want, got := metric.Metric[0].Gauge.GetValue(), float64(5); want != got {
		t.Errorf("Expect metric value %f, got %f", want, got)
	}
}
