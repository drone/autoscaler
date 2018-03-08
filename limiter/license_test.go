// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package limiter

import (
	"encoding/json"
	"testing"
	"time"

	"crypto/rand"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

func TestParseVerify(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	a := &License{
		Pro: "stripe",
		Sub: "cus_CS7Nimer2KxGNp",
		Lim: 25,
		Exp: time.Now().UTC(),
	}

	data, _ := json.Marshal(a)
	token, err := paseto.NewV2().Sign(private, data)
	if err != nil {
		t.Errorf(token)
	}

	b, err := ParseVerify(token, public)
	if err != nil {
		t.Error(err)
	}

	if want, got := a.Sub, b.Sub; want != got {
		t.Errorf("Want Sub %s, got %s", want, got)
	}
	if want, got := a.Pro, b.Pro; want != got {
		t.Errorf("Want Pro %s, got %s", want, got)
	}
	if want, got := a.Lim, b.Lim; want != got {
		t.Errorf("Want Lim %d, got %d", want, got)
	}
	if want, got := a.Exp, b.Exp; want != got {
		t.Errorf("Want Exp %s, got %s", want, got)
	}
}

func TestExpired(t *testing.T) {
	tests := []struct {
		Exp     time.Time
		expired bool
	}{
		// zero value indicates no time limit
		{
			Exp:     time.Time{},
			expired: false,
		},
		// one hour in the future
		{
			Exp:     time.Now().Add(time.Hour),
			expired: false,
		},
		// one hour in the past
		{
			Exp:     time.Now().Add(-1 * time.Hour),
			expired: true,
		},
	}

	for _, test := range tests {
		l := License{Exp: test.Exp}
		if got, want := l.Expired(), test.expired; got != want {
			t.Errorf("Want expired %v, got %v for %s", want, got, l.Exp)
		}
	}
}
