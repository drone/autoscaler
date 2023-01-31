// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/h2non/gock"
)

func TestDestroy(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Delete("/compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance").
		Reply(200).
		BodyString(`{ "name": "operation-name" }`)

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/operation-name").
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID:     "my-instance",
		Region: "us-central1-a",
	}

	p, err := New(
		WithClient(http.DefaultClient),
		WithZones("us-central1-a"),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}

func TestDestroy_Error(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Delete("/compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance").
		Reply(404)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID:     "my-instance",
		Region: "us-central1-a",
	}

	p, err := New(
		WithClient(http.DefaultClient),
		WithZones("us-central1-a"),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = p.Destroy(mockContext, mockInstance)
	if err == nil {
		t.Errorf("Expect error deleting server")
	}
}
