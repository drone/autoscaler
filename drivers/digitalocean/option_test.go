// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package digitalocean

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithImage("ubuntu-18-04-x64"),
		WithRegion("nyc3"),
		WithSize("s-8vcpu-32gb"),
		WithSSHKey("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithTags("drone", "agent"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	if got, want := p.image, "ubuntu-18-04-x64"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.region, "nyc3"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.size, "s-8vcpu-32gb"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := p.key, "58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"; got != want {
		t.Errorf("Want key %q, got %q", want, got)
	}
	if got, want := p.token, "77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"; got != want {
		t.Errorf("Want token %q, got %q", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
}
