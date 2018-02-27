// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http/httptest"
	"testing"

	"github.com/drone/autoscaler/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
)

func TestHandleScalerPause(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/pause", nil)

	scaler := mocks.NewMockScaler(controller)
	scaler.EXPECT().Pause()

	router := chi.NewRouter()
	router.Post("/api/pause", HandleScalerPause(scaler))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 204; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func TestHandleScalerResume(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/resume", nil)

	scaler := mocks.NewMockScaler(controller)
	scaler.EXPECT().Resume()

	router := chi.NewRouter()
	router.Post("/api/resume", HandleScalerResume(scaler))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 204; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}
