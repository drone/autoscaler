// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithImage("ubuntu-17.04"),
		WithDatacenter("fsn1-dc8"),
		WithServerType("cx20"),
		WithSSHKey(23234),
		WithFirewall(12345),
		WithNetwork(23456),
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
	if got, want := p.firewall, 12345; got != want {
		t.Errorf("Want firewall %d, got %d", want, got)
	}
	if got, want := p.network, 23456; got != want {
		t.Errorf("Want network %d, got %d", want, got)
	}
}
