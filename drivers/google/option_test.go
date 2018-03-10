// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package google

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithImage("ubuntu-1604-lts"),
		WithProject("my-project"),
		WithSSHKey("id_rsa"),
		WithSize("c3.large"),
		WithTags("drone", "agent"),
		WithZone("us-central1-f"),
	).(*provider)

	if got, want := p.image, "ubuntu-1604-lts"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.project, "my-project"; got != want {
		t.Errorf("Want project %q, got %q", want, got)
	}
	if got, want := p.key, "id_rsa"; got != want {
		t.Errorf("Want ssh key %q, got %q", want, got)
	}
	if got, want := p.size, "c3.large"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
	if got, want := p.zone, "us-central1-f"; got != want {
		t.Errorf("Want zone %q, got %q", want, got)
	}
}
