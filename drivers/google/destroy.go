// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"net/http"

	"github.com/drone/autoscaler"
	"google.golang.org/api/googleapi"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	op, err := p.service.Instances.Delete(p.project, p.zone, instance.ID).Context(ctx).Do()
	if err != nil {
		// https://github.com/googleapis/google-api-go-client/blob/master/googleapi/googleapi.go#L135
		if gerr, ok := err.(*googleapi.Error); ok &&
			gerr.Code == http.StatusNotFound {
			return autoscaler.ErrInstanceNotFound
		}
		return err
	}
	return p.waitZoneOperation(ctx, op.Name)
}
