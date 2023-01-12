// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import "testing"

func TestDefaults(t *testing.T) {
	p := New().(*provider)
	if got, want := p.apiURL, "https://api.gb1.brightbox.com"; got != want {
		t.Errorf("Want API URL %q, got %q", want, got)
	}
	if got, want := p.image, "img-sm72o"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.serverType, "typ-1ni84"; got != want {
		t.Errorf("Want server type %q, got %q", want, got)
	}
}
