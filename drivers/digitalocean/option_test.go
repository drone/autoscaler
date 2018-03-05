// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
		WithFingerprint("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithRegion("nyc3"),
		WithSize("s-8vcpu-32gb"),
		WithImage("ubuntu-16-04-x64"),
		WithTags("drone", "agent"),
	).(*provider)

	if got, want := p.image, "ubuntu-16-04-x64"; got != want {
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
