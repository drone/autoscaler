// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"errors"
	"net/http"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// Destroy deletes instance. destroyTimeout bounds the whole call - both
// asking Google to delete and waiting for the delete operation to finish:
// the reaper passes the engine's root ctx, which has no deadline of its own,
// so without this an operation that never reaches a terminal state would
// poll forever in waitOperation.
//
// Any error leaves the server in the engine's error state for the reaper to
// retry, including a timeout while waiting on the operation - the delete may
// well have succeeded, in which case the retry's 404 maps to
// ErrInstanceNotFound and is treated as done.
func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	ctx, cancel := context.WithTimeout(ctx, p.destroyTimeout)
	defer cancel()

	// An instance's Region is actually a Zone in the google provider
	var op *compute.Operation
	err := doWithRetry(ctx, "instances.delete", p.maxTransientDelay, func() error {
		var err error
		op, err = p.service.Instances.Delete(p.project, instance.Region, instance.ID).Context(ctx).Do()
		if err != nil {
			// https://github.com/googleapis/google-api-go-client/blob/master/googleapi/googleapi.go#L135
			if gerr, ok := err.(*googleapi.Error); ok &&
				gerr.Code == http.StatusNotFound {
				return autoscaler.ErrInstanceNotFound
			}
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = p.waitZoneOperation(ctx, op.Name, instance.Region)
	if err != nil && !errors.Is(err, autoscaler.ErrInstanceNotFound) {
		entry := logger.FromContext(ctx).
			WithField("zone", instance.Region).
			WithField("name", instance.Name).
			WithField("operation", op.Name).
			WithError(err)

		var opErr *operationError
		if errors.As(err, &opErr) {
			entry.WithField("code", opErr.Code).
				Errorln("delete operation failed")
		} else {
			entry.Warnln("delete was accepted but unable to confirm the instance was deleted")
		}
	}
	return err
}
