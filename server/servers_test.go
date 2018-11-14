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

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/kr/pretty"
)

func TestHandleServerList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/servers", nil)

	servers := []*autoscaler.Server{
		{Name: "server1", Capacity: 1},
		{Name: "server2", Capacity: 1},
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(servers, nil)

	HandleServerList(store).ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := []*autoscaler.Server{}, servers
	json.NewDecoder(w.Body).Decode(&got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response body does match expected result")
		pretty.Ldiff(t, got, want)
	}
}

func TestHandleServerListErr(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/servers", nil)

	err := errors.New("not found")
	store := mocks.NewMockServerStore(controller)
	store.EXPECT().List(gomock.Any()).Return(nil, err)

	HandleServerList(store).ServeHTTP(w, r)

	if got, want := w.Code, 500; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestHandleServerFind(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/servers/server1", nil)

	server := &autoscaler.Server{Name: "server1", Capacity: 1}
	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), "server1").Return(server, nil)

	router := chi.NewRouter()
	router.Get("/api/servers/{name}", HandleServerFind(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	got, want := &autoscaler.Server{}, server
	json.NewDecoder(w.Body).Decode(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response body does match expected result")
		pretty.Ldiff(t, got, want)
	}
}

func TestHandleServerFindErr(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/servers/server1", nil)

	err := errors.New("not found")
	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), "server1").Return(nil, err)

	router := chi.NewRouter()
	router.Get("/api/servers/{name}", HandleServerFind(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 404; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestHandleServerCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/servers", nil)

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	HandleServerCreate(store, config.Config{}).ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func TestHandleServerCreateFailure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/servers", nil)

	err := errors.New("oops")
	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Create(gomock.Any(), gomock.Any()).Return(err)

	h := HandleServerCreate(store, config.Config{})
	h.ServeHTTP(w, r)

	if got, want := w.Code, 500; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestHandleServerDelete(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/servers/i-5203422c", nil)

	server := &autoscaler.Server{
		Name:   "i-5203422c",
		Image:  "docker-16-04",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), server.Name).Return(server, nil)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	router := chi.NewRouter()
	router.Delete("/api/servers/{name}", HandleServerDelete(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
	if got, want := server.State, autoscaler.StateShutdown; got != want {
		t.Errorf("Want server state Shutdown, got %s", got)
	}
}

func TestHandleServerDeleteNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/servers/i-5203422c", nil)

	err := errors.New("not found")

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), "i-5203422c").Return(nil, err)

	router := chi.NewRouter()
	router.Delete("/api/servers/{name}", HandleServerDelete(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 404; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestHandleServerDeleteFailure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/servers/i-5203422c", nil)

	server := &autoscaler.Server{
		Name:   "i-5203422c",
		Image:  "docker-16-04",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
	}

	err := errors.New("bad request")

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), server.Name).Return(server, nil)
	store.EXPECT().Update(gomock.Any(), server).Return(err)

	router := chi.NewRouter()
	router.Delete("/api/servers/{name}", HandleServerDelete(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 500; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestHandleServerDeleteErrorState(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/servers/i-5203422c", nil)

	server := &autoscaler.Server{
		ID:     "",
		State:  autoscaler.StateError,
		Name:   "i-5203422c",
		Image:  "docker-16-04",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), server.Name).Return(server, nil)
	store.EXPECT().Delete(gomock.Any(), server).Return(nil)

	router := chi.NewRouter()
	router.Delete("/api/servers/{name}", HandleServerDelete(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 204; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func TestHandleServerForceDeleteErrorState(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/servers/i-5203422c?force=true", nil)

	server := &autoscaler.Server{
		ID:     "i-5203422c",
		State:  autoscaler.StateError,
		Name:   "i-5203422c",
		Image:  "docker-16-04",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
	}

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Find(gomock.Any(), server.Name).Return(server, nil)
	store.EXPECT().Delete(gomock.Any(), server).Return(nil)

	router := chi.NewRouter()
	router.Delete("/api/servers/{name}", HandleServerDelete(store))
	router.ServeHTTP(w, r)

	if got, want := w.Code, 204; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}
