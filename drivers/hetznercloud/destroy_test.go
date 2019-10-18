// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import (
	"context"
	"strconv"
	"testing"

	"github.com/drone/autoscaler"

	"github.com/h2non/gock"
)

func TestDestroy(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Delete("/v1/servers/3164494").
		Reply(200)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	p := New(
		WithToken("LRK9DAWQ1ZAEFSrCNEEzLCUwhYX1U3g7wMg4dTlkkDC96fyDuyJ39nVbVjCKSDfj"),
	)
	err := p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}

func TestDestroyDeleteError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.hetzner.cloud").
		Delete("/v1/servers/3164494").
		Reply(500)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	p := New(
		WithToken("LRK9DAWQ1ZAEFSrCNEEzLCUwhYX1U3g7wMg4dTlkkDC96fyDuyJ39nVbVjCKSDfj"),
	)
	err := p.Destroy(mockContext, mockInstance)
	if err == nil {
		t.Errorf("Expect error returned from hetzner cloud")
	}
}

func TestDestroyInvalidInput(t *testing.T) {
	i := &autoscaler.Instance{}
	p := provider{}
	err := p.Destroy(context.TODO(), i)
	if _, ok := err.(*strconv.NumError); !ok {
		t.Errorf("Expected invalid or missing ID error")
	}
}
