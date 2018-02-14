// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"github.com/drone/autoscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// ServerCount provides metrics for server counts.
func ServerCount(store autoscaler.ServerStore) autoscaler.ServerStore {
	prometheus.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "drone_server_count",
			Help: "Total number of active servers.",
		}, func() float64 {
			servers, _ := store.List(noContext)
			return float64(len(servers))
		}),
	)
	return store
}
