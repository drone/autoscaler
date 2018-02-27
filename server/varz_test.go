// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drone/autoscaler/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/kr/pretty"
)

func TestHandleVarz(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockVarz := &varz{
		Paused: true,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/varz", nil)

	scaler := mocks.NewMockScaler(controller)
	scaler.EXPECT().Paused().Return(true)

	router := chi.NewRouter()
	router.Post("/varz", HandleVarz(scaler))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &varz{}, mockVarz
	json.NewDecoder(w.Body).Decode(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response body does match expected result")
		pretty.Ldiff(t, got, want)
	}
}
