// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import "testing"

func TestDefaults(t *testing.T) {
	p := New().(*provider)
	if got, want := p.image, "ubuntu-20.04"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.datacenter, ""; got != want {
		t.Errorf("Want datacenter %q, got %q", want, got)
	}
	if got, want := p.serverType, "cx11"; got != want {
		t.Errorf("Want server type %q, got %q", want, got)
	}
}
