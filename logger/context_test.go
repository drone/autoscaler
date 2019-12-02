// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package logger

import (
	"context"
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	entry := Discard()

	ctx := WithContext(context.Background(), entry)
	got := FromContext(ctx)

	if got != entry {
		t.Errorf("Expected Logger from context")
	}
}

func TestEmptyContext(t *testing.T) {
	got := FromContext(context.Background())
	if got == nil {
		t.Errorf("Expected Logger from context")
	}
	if _, ok := got.(*discard); !ok {
		t.Errorf("Expected discard Logger from context")
	}
}

func TestRequest(t *testing.T) {
	entry := Discard()

	ctx := WithContext(context.Background(), entry)
	req := new(http.Request)
	req = req.WithContext(ctx)

	got := FromRequest(req)

	if got != entry {
		t.Errorf("Expected Logger from http.Request")
	}
}
