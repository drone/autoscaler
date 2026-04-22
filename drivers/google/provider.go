// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/avast/retry-go"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/drone/autoscaler/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

var (
	defaultTags = []string{
		"allow-docker",
	}

	defaultScopes = []string{
		"https://www.googleapis.com/auth/devstorage.read_only",
		"https://www.googleapis.com/auth/logging.write",
		"https://www.googleapis.com/auth/monitoring.write",
		"https://www.googleapis.com/auth/trace.append",
	}
)

// provider implements a Google Cloud Platform provider.
type provider struct {
	init sync.Once

	diskSize            int64
	diskType            string
	image               string
	labels              map[string]string
	network             string
	subnetwork          string
	stackType           string
	project             string
	privateIP           bool
	scopes              []string
	serviceAccountEmail string
	size                string
	tags                []string
	zones               []string
	userdata            *template.Template
	userdataKey         string

	rateLimiter *rate.Limiter

	service *compute.Service
}

// New returns a new Google Cloud Platform provider.
func New(opts ...Option) (autoscaler.Provider, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.diskSize == 0 {
		p.diskSize = 50
	}
	if p.diskType == "" {
		p.diskType = "pd-standard"
	}
	if len(p.zones) == 0 {
		p.zones = []string{"us-central1-a"}
	}
	if p.size == "" {
		p.size = "n1-standard-1"
	}
	if p.image == "" {
		p.image = "ubuntu-os-cloud/global/images/ubuntu-2004-focal-v20220712"
	}
	if p.network == "" {
		p.network = "global/networks/default"
	}
	if p.stackType == "" {
		p.stackType = "IPV4_ONLY"
	}
	if p.userdata == nil {
		p.userdata = userdata.T
	}
	if p.userdataKey == "" {
		p.userdataKey = "user-data"
	}
	if len(p.tags) == 0 {
		p.tags = defaultTags
	}
	if len(p.scopes) == 0 {
		p.scopes = defaultScopes
	}
	if p.serviceAccountEmail == "" {
		p.serviceAccountEmail = "default"
	}

	if p.rateLimiter == nil {
		// If unspecified, set to the max read rate limit for the API 25/s
		// Source: https://cloud.google.com/compute/docs/api-rate-limits
		p.rateLimiter = rate.NewLimiter(rate.Every(time.Second/25), 1)
	}

	if p.service == nil {
		client, err := google.DefaultClient(oauth2.NoContext, compute.ComputeScope)
		if err != nil {
			return nil, err
		}
		p.service, err = compute.New(client)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *provider) waitZoneOperation(ctx context.Context, name string, zone string) error {
	var op *compute.Operation
	for {
		if err := p.rateLimiter.Wait(ctx); err != nil {
			return err
		}
		err := retry.Do(
			func() error {
				var err error
				op, err = p.service.ZoneOperations.Get(p.project, zone, name).Context(ctx).Do()
				if err != nil {
					if gerr, ok := err.(*googleapi.Error); ok &&
						gerr.Code == http.StatusNotFound {
						return retry.Unrecoverable(autoscaler.ErrInstanceNotFound)
					}
					if isTransientError(err) {
						return err
					}
					return retry.Unrecoverable(err)
				}
				if op == nil {
					return errors.New("zone operation get returned nil response")
				}
				return nil
			},
			retry.Attempts(5),
			retry.Context(ctx),
			retry.LastErrorOnly(true),
			retry.MaxDelay(time.Second*5),
			retry.OnRetry(func(n uint, err error) {
				logger.FromContext(ctx).
					WithField("attempt", n+1).
					WithField("operation", name).
					WithField("zone", zone).
					WithError(err).
					Debugln("retrying zone operation poll")
			}),
		)
		if err != nil {
			return err
		}
		if op == nil {
			return errors.New("zone operation response was nil after retry")
		}

		if op.Error != nil && len(op.Error.Errors) > 0 {
			return errors.New(op.Error.Errors[0].Message)
		}
		if op.Status == "DONE" {
			return nil
		}
	}
}

func (p *provider) waitGlobalOperation(ctx context.Context, name string) error {
	var op *compute.Operation
	for {
		if err := p.rateLimiter.Wait(ctx); err != nil {
			return err
		}
		err := retry.Do(
			func() error {
				var err error
				op, err = p.service.GlobalOperations.Get(p.project, name).Context(ctx).Do()
				if err != nil {
					if isTransientError(err) {
						return err
					}
					return retry.Unrecoverable(err)
				}
				if op == nil {
					return errors.New("global operation get returned nil response")
				}
				return nil
			},
			retry.Attempts(5),
			retry.Context(ctx),
			retry.LastErrorOnly(true),
			retry.MaxDelay(time.Second*5),
			retry.OnRetry(func(n uint, err error) {
				logger.FromContext(ctx).
					WithField("attempt", n+1).
					WithField("operation", name).
					WithError(err).
					Debugln("retrying global operation poll")
			}),
		)
		if err != nil {
			return err
		}
		if op == nil {
			return errors.New("global operation response was nil after retry")
		}

		if op.Error != nil && len(op.Error.Errors) > 0 {
			return errors.New(op.Error.Errors[0].Message)
		}
		if op.Status == "DONE" {
			return nil
		}
	}
}
