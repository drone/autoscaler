// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package packet

import (
	"testing"

	"github.com/drone/autoscaler/drivers/internal/userdata"
)

func TestDefaults(t *testing.T) {
	n, e := New()
	if e != nil {
		t.Error(e)
		return
	}
	p := n.(*provider)
	if got, want := p.plan, "baremetal_0"; got != want {
		t.Errorf("Want plan %q, got %q", want, got)
	}
	if got, want := p.facility, "ewr1"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.billing, "hourly"; got != want {
		t.Errorf("Want billing %q, got %q", want, got)
	}
	if got, want := p.os, "ubuntu_16_04"; got != want {
		t.Errorf("Want os %q, got %q", want, got)
	}
	if p.userdata != userdata.Ubuntu {
		t.Errorf("Want ubuntu userdata template")
	}
}
