// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http/httptest"
	"testing"

	"github.com/drone/autoscaler/mocks"
	"github.com/golang/mock/gomock"
)

func TestHandleEnginePause(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/pause", nil)

	e := mocks.NewMockEngine(controller)
	e.EXPECT().Pause()

	HandleEnginePause(e).ServeHTTP(w, r)

	if got, want := w.Code, 204; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func TestHandleEngineResume(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/resume", nil)

	e := mocks.NewMockEngine(controller)
	e.EXPECT().Resume()

	HandleEngineResume(e).ServeHTTP(w, r)

	if got, want := w.Code, 204; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}
