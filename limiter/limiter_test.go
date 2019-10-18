// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/mocks"
	"github.com/golang/mock/gomock"
)

func TestLimiter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x4 server count
	servers := []*autoscaler.Server{
		{Name: "server1", State: autoscaler.StateRunning},
		{Name: "server2", State: autoscaler.StateRunning},
		{Name: "server3", State: autoscaler.StateRunning},
		{Name: "server4", State: autoscaler.StateRunning},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)
	store.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	limiter := &limiter{store, &License{
		Lim: 5,
	}}

	err := limiter.Create(context.TODO(), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestLimiter_ErrExpired(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mocks.NewMockServerStore(controller)
	limiter := &limiter{store, &License{
		Lim: 5,
		Exp: time.Now().Add(time.Hour * -1),
	}}

	err := limiter.Create(context.TODO(), nil)
	if err != errLicenseExpired {
		t.Errorf("Want license expried error, got %s", err)
	}
}

func TestLimiter_ErrExceeded(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// x4 server count
	servers := []*autoscaler.Server{
		{Name: "server1", State: autoscaler.StateRunning},
		{Name: "server2", State: autoscaler.StateRunning},
		{Name: "server3", State: autoscaler.StateRunning},
		{Name: "server4", State: autoscaler.StateRunning},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	limiter := &limiter{store, &License{
		Lim: 4,
		Exp: time.Now().Add(time.Hour),
	}}

	err := limiter.Create(context.TODO(), nil)
	if err != errServerLimitExceeded {
		t.Errorf("Want limit exceeded error, got %s", err)
	}
}
