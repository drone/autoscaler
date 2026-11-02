// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
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

	// createTimeout bounds one Create attempt: the insert, its operation poll
	// and the follow-up get.
	createTimeout time.Duration

	// destroyTimeout bounds Destroy: asking Google to delete and waiting for
	// the delete operation to finish.
	destroyTimeout time.Duration

	// maxTransientDelay caps the exponential+jitter backoff doWithRetry
	// applies to ordinary transient errors (not rate limits, which use the
	// server's Retry-After instead).
	maxTransientDelay time.Duration

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
	if p.createTimeout <= 0 {
		p.createTimeout = 3 * time.Minute
	}
	if p.destroyTimeout == 0 {
		p.destroyTimeout = 15 * time.Minute
	}
	if p.maxTransientDelay == 0 {
		p.maxTransientDelay = 5 * time.Second
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

// operationError preserves Google's structured Code/Message operation-error
// fields instead of collapsing them to a plain string.
type operationError struct {
	Code    string
	Message string
}

func (e *operationError) Error() string {
	return e.Message
}

// waitOperation polls a compute operation until it is DONE, returning its
// error (if any) only then, or until ctx ends. get performs one poll; it is
// shared by waitZoneOperation and waitGlobalOperation, which differ only in
// the API call they make.
func (p *provider) waitOperation(ctx context.Context, get func() (*compute.Operation, error)) error {
	for {
		op, err := get()
		if err != nil {
			var rl *retryAfterError
			if errors.As(err, &rl) {
				if err := sleepOrDone(ctx, rl.retryAfter); err != nil {
					return err
				}
				continue
			}
			return err
		}

		// Only act on the outcome once the operation is DONE: an operation
		// can report its error while still RUNNING. A failed instance insert
		// (e.g. a stockout) does, while Compute Engine tears the instance
		// down, and the instance name stays reserved until then.
		if op.Status == "DONE" {
			if op.Error != nil && len(op.Error.Errors) > 0 {
				return &operationError{
					Code:    op.Error.Errors[0].Code,
					Message: op.Error.Errors[0].Message,
				}
			}
			return nil
		}

		// Keep polling cadence at ~1s while still honoring cancellation.
		if err := sleepOrDone(ctx, time.Second); err != nil {
			return err
		}
	}
}

func (p *provider) waitZoneOperation(ctx context.Context, name string, zone string) error {
	return p.waitOperation(ctx, func() (*compute.Operation, error) {
		var op *compute.Operation
		err := doWithRetry(ctx, fmt.Sprintf("zone operations.get %s", name), p.maxTransientDelay, func() error {
			if err := p.rateLimiter.Wait(ctx); err != nil {
				return err
			}
			var err error
			op, err = p.service.ZoneOperations.Get(p.project, zone, name).Context(ctx).Do()
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok &&
					gerr.Code == http.StatusNotFound {
					return autoscaler.ErrInstanceNotFound
				}
				return err
			}
			if op == nil {
				return errors.New("zone operation get returned nil response")
			}
			return nil
		})
		return op, err
	})
}

func (p *provider) waitGlobalOperation(ctx context.Context, name string) error {
	return p.waitOperation(ctx, func() (*compute.Operation, error) {
		var op *compute.Operation
		err := doWithRetry(ctx, fmt.Sprintf("global operations.get %s", name), p.maxTransientDelay, func() error {
			if err := p.rateLimiter.Wait(ctx); err != nil {
				return err
			}
			var err error
			op, err = p.service.GlobalOperations.Get(p.project, name).Context(ctx).Do()
			if err != nil {
				return err
			}
			if op == nil {
				return errors.New("global operation get returned nil response")
			}
			return nil
		})
		return op, err
	})
}
