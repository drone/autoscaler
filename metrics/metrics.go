// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var noContext = context.Background()

// Collector defines a metrics collector.
type Collector interface {
	// TrackServerCreateTime registers the elapsed time it takes
	// to provision a server instance.
	TrackServerCreateTime(start time.Time)

	// TrackServerInitTime registers the elapsed time it takes
	// for a server instance to initialize and begin accepting
	// network connections.
	TrackServerInitTime(start time.Time)

	// TrackServerSetupTime registers the elapsed time it takes
	// to install software (i.e. docker, runners) on the server.
	TrackServerSetupTime(start time.Time)

	// IncrServerCreateError keeps a count of errors encountered
	// when provisioning servers.
	IncrServerCreateError()

	// IncrServerInitError keeps a count of errors encountered
	// when initializing and establishing networking connections
	// with servers.
	IncrServerInitError()

	// IncrServerSetupError keeps a count of errors encountered
	// when installing software on servers.
	IncrServerSetupError()
}

// Prometheus is a Prometheus metrics collector.
type Prometheus struct {
	trackServerCreateTime prometheus.Histogram
	trackServerInitTime   prometheus.Histogram
	trackServerSetupTime  prometheus.Histogram
	countServerCreateErr  prometheus.Counter
	countServerInitErr    prometheus.Counter
	countServerSetupErr   prometheus.Counter
}

// New returns a new Prometheus metrics provider.
func New() *Prometheus {
	p := new(Prometheus)
	p.trackServerCreateTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "drone_server_create_time_seconds",
		Help:    "Elapsed time creating a server.",
		Buckets: []float64{60, 150, 300, 600, 900, 1200},
	})
	p.trackServerInitTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "drone_server_boot_time_seconds",
		Help:    "Elapsed time initializing a server.",
		Buckets: []float64{60, 150, 300, 600, 900, 1200},
	})
	p.trackServerSetupTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "drone_server_install_time_seconds",
		Help:    "Elapsed time installing software on a server.",
		Buckets: []float64{60, 150, 300, 600, 900, 1200},
	})
	p.countServerCreateErr = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_server_create_errors_total",
		Help: "Total number of errors initializing a server.",
	})
	p.countServerInitErr = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_server_boot_errors_total",
		Help: "Total number of errors initializing a server.",
	})
	p.countServerSetupErr = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "drone_server_install_errors_total",
		Help: "Total number of errors installing software on a server.",
	})
	prometheus.MustRegister(p.trackServerCreateTime)
	prometheus.MustRegister(p.trackServerInitTime)
	prometheus.MustRegister(p.trackServerSetupTime)
	prometheus.MustRegister(p.countServerCreateErr)
	prometheus.MustRegister(p.countServerInitErr)
	prometheus.MustRegister(p.countServerSetupErr)
	return p
}

// TrackServerCreateTime registers the elapsed time it takes
// to provision a server instance.
func (m *Prometheus) TrackServerCreateTime(start time.Time) {
	m.trackServerCreateTime.Observe(
		time.Now().Sub(start).Round(time.Second).Seconds(),
	)
}

// TrackServerInitTime registers the elapsed time it takes
// for a server instance to initialize and begin accepting
// network connections.
func (m *Prometheus) TrackServerInitTime(start time.Time) {
	m.trackServerInitTime.Observe(
		time.Now().Sub(start).Round(time.Second).Seconds(),
	)
}

// TrackServerSetupTime registers the elapsed time it takes
// to install software (i.e. docker, runners) on the server.
func (m *Prometheus) TrackServerSetupTime(start time.Time) {
	m.trackServerSetupTime.Observe(
		time.Now().Sub(start).Round(time.Second).Seconds(),
	)
}

// IncrServerCreateError keeps a count of errors encountered
// when provisioning servers.
func (m *Prometheus) IncrServerCreateError() {
	m.countServerCreateErr.Inc()
}

// IncrServerInitError keeps a count of errors encountered
// when initializing and establishing networking connections
// with servers.
func (m *Prometheus) IncrServerInitError() {
	m.countServerInitErr.Inc()
}

// IncrServerSetupError keeps a count of errors encountered
// when installing software on servers.
func (m *Prometheus) IncrServerSetupError() {
	m.countServerSetupErr.Inc()
}

// NopCollector provides a no-op metrics collector.
type NopCollector struct{}

// TrackServerCreateTime registers the elapsed time it takes
// to provision a server instance.
func (*NopCollector) TrackServerCreateTime(start time.Time) {}

// TrackServerInitTime registers the elapsed time it takes
// for a server instance to initialize and begin accepting
// network connections.
func (*NopCollector) TrackServerInitTime(start time.Time) {}

// TrackServerSetupTime registers the elapsed time it takes
// to install software (i.e. docker, runners) on the server.
func (*NopCollector) TrackServerSetupTime(start time.Time) {}

// IncrServerCreateError keeps a count of errors encountered
// when provisioning servers.
func (*NopCollector) IncrServerCreateError() {}

// IncrServerInitError keeps a count of errors encountered
// when initializing and establishing networking connections
// with servers.
func (*NopCollector) IncrServerInitError() {}

// IncrServerSetupError keeps a count of errors encountered
// when installing software on servers.
func (*NopCollector) IncrServerSetupError() {}
