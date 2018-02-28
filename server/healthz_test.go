// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandleHealthz(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/healthz", nil)

	HandleHealthz().ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}
