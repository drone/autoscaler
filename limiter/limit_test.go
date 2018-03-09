// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package limiter

import (
	"os"
	"testing"

	"github.com/drone/autoscaler/mocks"
	"github.com/golang/mock/gomock"
)

func TestLimit_EmptyToken(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mocks.NewMockServerStore(controller)

	l, ok := Limit(store, "").(*limiter)
	if !ok {
		t.Errorf("Want limiter type")
	}

	if got, want := l.license.Lim, 10; got != want {
		t.Errorf("Want limit %d, got %d", want, got)
	}
	if got, want := l.ServerStore, store; got != want {
		t.Errorf("Want store wrapped by limiter")
	}
}

func TestLimit(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Skip()
	}
	l, ok := Limit(nil, token).(*limiter)
	if !ok {
		t.Errorf("Want limiter type")
	}
	if l == nil || l.license == nil {
		t.Errorf("Want license parsed")
	}
	if l.license.Lim == 10 {
		t.Errorf("Want license, got trial license")
	}
}
