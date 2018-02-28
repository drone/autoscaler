// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package metrics

import (
	"context"

	"github.com/drone/autoscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// ServerCreate provides metrics for servers created.
func ServerCreate(provider autoscaler.Provider) autoscaler.Provider {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_servers_created",
		Help: "Total number of servers created.",
	})
	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_servers_created_err",
		Help: "Total number of server creation errors.",
	})
	prometheus.MustRegister(counter)
	prometheus.MustRegister(errors)
	return &providerWrapCreate{
		Provider: provider,
		created:  counter,
		errors:   errors,
	}
}

// instruments the Provider to count server create events.
type providerWrapCreate struct {
	autoscaler.Provider
	created prometheus.Counter
	errors  prometheus.Counter
}

func (p *providerWrapCreate) Create(ctx context.Context, server *autoscaler.Server) error {
	err := p.Provider.Create(ctx, server)
	if err == nil {
		p.created.Add(1)
	} else {
		p.errors.Add(1)
	}
	return err
}
