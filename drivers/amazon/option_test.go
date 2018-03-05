// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithImage("ami-66506c1c"),
		WithRegion("us-west-2"),
		WithSecurityGroup("sg-770eabe1"),
		WithSize("t2.2xlarge"),
		WithSSHKey("id_rsa"),
		WithSubnet("subnet-0b32177f"),
		WithTags("drone", "agent"),
	).(*provider)

	if got, want := p.image, "ami-66506c1c"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.region, "us-west-2"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.size, "t2.2xlarge"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := p.key, "id_rsa"; got != want {
		t.Errorf("Want key %q, got %q", want, got)
	}
	if got, want := p.groups[0], "sg-770eabe1"; got != want {
		t.Errorf("Want security groups %q, got %q", want, got)
	}
	if got, want := p.subnet, "subnet-0b32177f"; got != want {
		t.Errorf("Want subnet %q, got %q", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
}
