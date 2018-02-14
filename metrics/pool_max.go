// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"github.com/drone/autoscaler/config"
	"github.com/prometheus/client_golang/prometheus"
)

// MaxPool exposes the min pool metric.
func MaxPool(config config.Config) {
	prometheus.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "drone_server_max_pool",
			Help: "Maximum number of active servers.",
		}, func() float64 {
			return float64(config.Pool.Max)
		}),
	)
}
