// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud"
)

func TestOptions(t *testing.T) {
	v, err := New(
		WithComputeClient(&gophercloud.ServiceClient{}),
		WithNetworkClient(&gophercloud.ServiceClient{}),
		WithFloatingIpPool("ext-ips-1"),
		WithFlavor("053dc448-045b-4c15-a4a0-1908b6b9310d"),
		WithSecurityGroup("drone-ci"),
		WithSSHKey("drone-ci"),
		WithRegion("sto-01"),
		WithImage("0e9fe318-568f-417e-b2c1-f1218aa2712f"),
		WithMetadata(map[string]string{"foo": "bar", "baz": "qux"}),
		WithNetwork("c7d172c8-96e6-40ab-aaaa-4a555e247c73"),
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
	if got, want := p.flavor, "053dc448-045b-4c15-a4a0-1908b6b9310d"; got != want {
		t.Errorf("Want flavor %q, got %q", want, got)
	}
	if got, want := p.image, "0e9fe318-568f-417e-b2c1-f1218aa2712f"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.network, "c7d172c8-96e6-40ab-aaaa-4a555e247c73"; got != want {
		t.Errorf("Want network %q, got %q", want, got)
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
