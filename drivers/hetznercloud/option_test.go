// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package hetznercloud

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithImage("ubuntu-17.04"),
		WithDatacenter("fsn1-dc8"),
		WithServerType("cx20"),
		WithSSHKey(23234),
	).(*provider)

	if got, want := p.image, "ubuntu-17.04"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.datacenter, "fsn1-dc8"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.serverType, "cx20"; got != want {
		t.Errorf("Want serverType %q, got %q", want, got)
	}
	if got, want := p.key, 23234; got != want {
		t.Errorf("Want key %d, got %d", want, got)
	}
}
