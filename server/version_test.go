// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/kr/pretty"
)

func TestHandleVersion(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/version", nil)

	mockVersion := &versionInfo{
		Source:  "github.com/octocat/hello-world",
		Version: "1.0.0",
		Commit:  "ad2aec",
	}

	router := chi.NewRouter()
	router.Get("/version", HandleVersion(mockVersion.Source, mockVersion.Version, mockVersion.Commit))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &versionInfo{}, mockVersion
	json.NewDecoder(w.Body).Decode(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response body does match expected result")
		pretty.Ldiff(t, got, want)
	}
}
