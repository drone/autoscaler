// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import "testing"

func TestDefaults(t *testing.T) {
	p := New().(*provider)
	if got, want := p.image, "ubuntu-16.04"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.datacenter, "nbg1-dc3"; got != want {
		t.Errorf("Want datacenter %q, got %q", want, got)
	}
	if got, want := p.serverType, "cx11"; got != want {
		t.Errorf("Want server type %q, got %q", want, got)
	}
}
