// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package metrics

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// ServerDelete provides metrics for servers deleted.
func ServerDelete(provider autoscaler.Provider) autoscaler.Provider {
	created := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_servers_deleted",
		Help: "Total number of servers deleted.",
	})
	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_servers_deleted_err",
		Help: "Total number of server deletion errors.",
	})
	prometheus.MustRegister(created)
	prometheus.MustRegister(errors)
	return &providerWrapDestroy{
		Provider: provider,
		created:  created,
		errors:   errors,
	}
}

// instruments the Provider to count server destroy events.
type providerWrapDestroy struct {
	autoscaler.Provider
	created prometheus.Counter
	errors  prometheus.Counter
}

func (p *providerWrapDestroy) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	err := p.Provider.Destroy(ctx, instance)
	if err == nil {
		p.created.Add(1)
	} else {
		p.errors.Add(1)
	}
	return err
}
