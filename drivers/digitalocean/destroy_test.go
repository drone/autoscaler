// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/drone/autoscaler"

	"github.com/golang/mock/gomock"
	"github.com/h2non/gock"
)

func TestDestroy(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Delete("/v2/droplets/3164494").
		Reply(204)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	p := New(
		WithSSHKey("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	)

	err := p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}

func TestDestroyDeleteError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Delete("/v2/droplets/3164494").
		Reply(500)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	p := New(
		WithSSHKey("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	)

	err := p.Destroy(mockContext, mockInstance)
	if err == nil {
		t.Errorf("Expect error returned from digital ocean")
	} else if _, ok := err.(*godo.ErrorResponse); !ok {
		t.Errorf("Expect ErrorResponse digital ocean")
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
