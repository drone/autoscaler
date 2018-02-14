// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"github.com/drone/autoscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// ServerCapacity provides metrics for server capacity count.
func ServerCapacity(store autoscaler.ServerStore) autoscaler.ServerStore {
	prometheus.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "drone_server_capacity",
			Help: "Total capacity of active servers.",
		}, func() float64 {
			var capacity int
			servers, _ := store.List(noContext)
			for _, server := range servers {
				capacity += server.Capacity
			}
			return float64(capacity)
		}),
	)
	return store
}
