// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drone/autoscaler/config"
	"github.com/drone/drone-go/drone"

	"github.com/h2non/gock"
)

func TestAuthorize(t *testing.T) {
	defer gock.Off()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer NTE2M2MwMWRlYToxNGM3MWEyYTIx")

	user := &drone.User{
		Login: "octocat",
		Admin: true,
	}

	c := config.Config{}
	c.Server.Host = "company.drone.com"
	c.Server.Proto = "https"

	gock.New("https://company.drone.com").
		Get("/api/user").
		MatchHeader("Authorization", "Bearer NTE2M2MwMWRlYToxNGM3MWEyYTIx").
		Reply(200).
		JSON(user)

	CheckDrone(c)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}),
	).ServeHTTP(w, r)

	if got, want := w.Code, http.StatusTeapot; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}
}

func TestAuthorizeMissingToken(t *testing.T) {
	defer gock.Off()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	c := config.Config{}
	c.Server.Host = "company.drone.com"
	c.Server.Proto = "https"

	CheckDrone(c)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("Expect access to handler is restricted")
		}),
	).ServeHTTP(w, r)

	if got, want := w.Code, 401; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, errInvalidToken.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestAuthorizeNotFound(t *testing.T) {
	defer gock.Off()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer NTE2M2MwMWRlYToxNGM3MWEyYTIx")

	c := config.Config{}
	c.Server.Host = "company.drone.com"
	c.Server.Proto = "https"

	gock.New("https://company.drone.com").
		Get("/api/user").
		MatchHeader("Authorization", "Bearer NTE2M2MwMWRlYToxNGM3MWEyYTIx").
		Reply(404)

	CheckDrone(c)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("Expect access to handler is restricted")
		}),
	).ServeHTTP(w, r)

	if got, want := w.Code, 401; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, errUnauthorized.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestAuthorizeNonAdmin(t *testing.T) {
	defer gock.Off()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer NTE2M2MwMWRlYToxNGM3MWEyYTIx")

	user := &drone.User{
		Login: "octocat",
		Admin: false,
	}

	c := config.Config{}
	c.Server.Host = "company.drone.com"
	c.Server.Proto = "https"

	gock.New("https://company.drone.com").
		Get("/api/user").
		MatchHeader("Authorization", "Bearer NTE2M2MwMWRlYToxNGM3MWEyYTIx").
		Reply(200).
		JSON(user)

	CheckDrone(c)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("Expect access to handler is restricted")
		}),
	).ServeHTTP(w, r)

	if got, want := w.Code, 403; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	errjson := &Error{}
	json.NewDecoder(w.Body).Decode(errjson)
	if got, want := errjson.Message, errForbidden.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}
