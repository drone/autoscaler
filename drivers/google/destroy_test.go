// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/h2non/gock"
	compute "google.golang.org/api/compute/v1"
)

func TestDestroy(t *testing.T) {
	defer gock.Off()

	gock.New("https://www.googleapis.com").
		Delete("/compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance").
		Reply(200).
		BodyString(`{ "operation": { "name": "operation-name" }}`)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "my-instance",
	}

	p := New(
		WithZone("us-central1-a"),
		WithProject("my-project"),
	).(*provider)
	p.service, _ = compute.New(http.DefaultClient)

	err := p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}
