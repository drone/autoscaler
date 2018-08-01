// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package packet

import "testing"

func TestOptions(t *testing.T) {
	pp := New(
		WithAPIKey("my_authentication_token"),
		WithFacility("sjc1"),
		WithOS("ubuntu_16_10"),
		WithPlan("baremetal_1"),
		WithProject("my_project"),
		WithSSHKey("id_rsa"),
		WithTags("drone", "agent"),
	)
	p := pp.(*provider)

	if got, want := p.apikey, "my_authentication_token"; got != want {
		t.Errorf("Want api key %q, got %q", want, got)
	}
	if got, want := p.facility, "sjc1"; got != want {
		t.Errorf("Want facility %q, got %q", want, got)
	}
	if got, want := p.os, "ubuntu_16_10"; got != want {
		t.Errorf("Want os %q, got %q", want, got)
	}
	if got, want := p.plan, "baremetal_1"; got != want {
		t.Errorf("Want os %q, got %q", want, got)
	}
	if got, want := p.project, "my_project"; got != want {
		t.Errorf("Want os %q, got %q", want, got)
	}
	if got, want := p.sshkey, "id_rsa"; got != want {
		t.Errorf("Want os %q, got %q", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
}
