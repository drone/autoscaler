// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drone/autoscaler/mocks"
	"github.com/drone/drone-go/drone"
	"github.com/golang/mock/gomock"
	"github.com/kr/pretty"
)

func TestHandleQueueList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/queue", nil)

	// x4 running builds
	// x3 pending builds
	mockBuilds := []*drone.Activity{
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusRunning},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
		{Status: drone.StatusPending},
	}

	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(mockBuilds, nil)

	HandleQueueList(client).ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := []*drone.Activity{}, mockBuilds
	json.NewDecoder(w.Body).Decode(&got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response body does match expected result")
		pretty.Ldiff(t, got, want)
	}
}

func TestHandleQueueListErr(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/queue", nil)

	err := errors.New("pc load letter")
	client := mocks.NewMockClient(controller)
	client.EXPECT().BuildQueue().Return(nil, err)

	HandleQueueList(client).ServeHTTP(w, r)

	if got, want := w.Code, 500; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}
