// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package packet

import (
	"testing"

	"github.com/drone/autoscaler/drivers/internal/userdata"
)

func TestDefaults(t *testing.T) {
	p := New().(*provider)

	if got, want := p.plan, "baremetal_0"; got != want {
		t.Errorf("Want plan %q, got %q", want, got)
	}
	if got, want := p.facility, "ewr1"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.billing, "hourly"; got != want {
		t.Errorf("Want billing %q, got %q", want, got)
	}
	if got, want := p.os, "ubuntu_18_04"; got != want {
		t.Errorf("Want os %q, got %q", want, got)
	}
	if p.userdata != userdata.T {
		t.Errorf("Want default userdata template")
	}
}
