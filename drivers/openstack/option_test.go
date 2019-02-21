// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"github.com/gophercloud/gophercloud"
	"testing"
)

func TestOptions(t *testing.T) {
	v, err := New(
		WithComputeClient(&gophercloud.ServiceClient{}),
		WithFloatingIpPool("ext-ips-1"),
		WithFlavor("t1.medium"),
		WithSecurityGroup("drone-ci"),
		WithSSHKey("drone-ci"),
		WithRegion("sto-01"),
		WithImage("ubuntu-16.04-server-latest"),
		WithMetadata(map[string]string{"foo": "bar", "baz": "qux"}),
		WithSubnet("subnet-feedface"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)

	if got, want := p.pool, "ext-ips-1"; got != want {
		t.Errorf("Want pool %q, got %q", want, got)
	}
	if got, want := p.region, "sto-01"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.flavor, "t1.medium"; got != want {
		t.Errorf("Want flavor %q, got %q", want, got)
	}
	if got, want := p.image, "ubuntu-16.04-server-latest"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.subnet, "subnet-feedface"; got != want {
		t.Errorf("Want subnet %q, got %q", want, got)
	}
	if got, want := p.key, "drone-ci"; got != want {
		t.Errorf("Want key %q, got %q", want, got)
	}
	if got, want := len(p.metadata), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
	if got, want := p.metadata["foo"], "bar"; got != want {
		t.Errorf("Want foo=%q metadata, got foo=%q", want, got)
	}
	if got, want := p.metadata["baz"], "qux"; got != want {
		t.Errorf("Want baz=%q metadata, got baz=%q", want, got)
	}
}
