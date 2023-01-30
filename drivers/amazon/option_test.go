// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package amazon

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithDeviceName("/dev/sda2"),
		WithImage("ami-0aab355e1bfa1e72e"),
		WithPrivateIP(true),
		WithRegion("us-west-2"),
		WithRetries(10),
		WithSecurityGroup("sg-770eabe1"),
		WithSize("t3.2xlarge"),
		WithSSHKey("id_rsa"),
		WithSubnets([]string{"subnet-0b32177f"}),
		WithTags(map[string]string{"foo": "bar", "baz": "qux"}),
		WithVolumeSize(64),
		WithVolumeType("io1"),
	).(*provider)

	if got, want := p.deviceName, "/dev/sda2"; got != want {
		t.Errorf("Want device name %q, got %q", want, got)
	}
	if got, want := p.image, "ami-0aab355e1bfa1e72e"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.region, "us-west-2"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.size, "t3.2xlarge"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := p.key, "id_rsa"; got != want {
		t.Errorf("Want key %q, got %q", want, got)
	}
	if got, want := p.groups[0], "sg-770eabe1"; got != want {
		t.Errorf("Want security groups %q, got %q", want, got)
	}
	if got, want := p.subnets, []string{"subnet-0b32177f"}; len(got) != 1 || got[0] != want[0] {
		t.Errorf("Want subnet %q, got %q", want, got)
	}
	if got, want := p.retries, 10; got != want {
		t.Errorf("Want %d retries, got %d", want, got)
	}
	if got, want := p.privateIP, true; got != want {
		t.Errorf("Want %v privateIP, got %v", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
	if got, want := p.volumeSize, int64(64); got != want {
		t.Errorf("Want volume size %d, got %d", want, got)
	}
	if got, want := p.volumeType, "io1"; got != want {
		t.Errorf("Want volume type %q, got %q", want, got)
	}
}
