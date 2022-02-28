// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package server

import (
	"net/http/httptest"
	"testing"
)

const expectedGetHeaderTextPlain = "text/plain; version=0.0.4; charset=utf-8"

func TestHandleMetrics(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer correct-horse-batter-staple")

	HandleMetrics("correct-horse-batter-staple").ServeHTTP(w, r)

	if got, want := w.Code, 200; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	if got, want := w.HeaderMap.Get("Content-Type"), expectedGetHeaderTextPlain; got != want {
		t.Errorf("Want prometheus header %q, got %q", want, got)
	}
}

func TestHandleMetricsUnprotected(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	HandleMetrics("").ServeHTTP(w, r)

	if got, want := w.Code, 200; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	if got, want := w.HeaderMap.Get("Content-Type"), expectedGetHeaderTextPlain; got != want {
		t.Errorf("Want prometheus header %q, got %q", want, got)
	}
}

func TestHandleMetricsMissingToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	HandleMetrics("correct-horse-batter-staple").ServeHTTP(w, r)

	if got, want := w.Code, 401; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}
}

func TestHandleMetricsInvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "correct-horse-batter-staple")

	HandleMetrics("correct-horse-batter-staple").ServeHTTP(w, r)

	if got, want := w.Code, 401; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}
}
