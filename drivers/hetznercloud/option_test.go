// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package hetznercloud

import "testing"

var (
	mockLabels = map[string]string{
		"label": "test",
	}
)

func TestOptions(t *testing.T) {
	p := New(
		WithImage("ubuntu-17.04"),
		WithDatacenter("fsn1-dc8"),
		WithServerType("cx20"),
		WithSSHKey(23234),
		WithLabels(mockLabels),
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
	for mk, mv := range mockLabels {
		if func() bool {
			for k, v := range p.labels {
				if k == mk && v == mv {
					return false
				}
			}
			return true
		}() {
			t.Errorf("Want %s, got %s", mockLabels, p.labels)
		}
	}
}
