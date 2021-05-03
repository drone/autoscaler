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
		WithFirewall("f33e7128-f3e7-4229-b6cc-a4751381a104"),
		WithVpc("b5bc3aa3-374e-449f-a954-1ed2b6b5537f"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
		WithPrivateIP(false),
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
	if got, want := p.firewall, "f33e7128-f3e7-4229-b6cc-a4751381a104"; got != want {
		t.Errorf("Want firewall %q, got %q", want, got)
	}
	if got, want := p.vpc, "b5bc3aa3-374e-449f-a954-1ed2b6b5537f"; got != want {
		t.Errorf("Want vpc %q, got %q", want, got)
	}
	if got, want := p.privateIP, false; got != want {
		t.Errorf("Want %v privateIP, got %v", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
}
