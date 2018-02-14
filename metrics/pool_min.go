// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"github.com/drone/autoscaler/config"
	"github.com/prometheus/client_golang/prometheus"
)

// MinPool exposes the min pool metric.
func MinPool(config config.Config) {
	prometheus.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "drone_server_min_pool",
			Help: "Minimum number of active servers.",
		}, func() float64 {
			return float64(config.Pool.Min)
		}),
	)
}
