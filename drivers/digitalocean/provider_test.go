// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package digitalocean

import "testing"

func TestDefaults(t *testing.T) {
	p := New().(*provider)
	if got, want := p.image, "docker-18-04"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.region, "nyc1"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.size, "s-2vcpu-4gb"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := p.key, ""; got != want {
		t.Errorf("Want key %q, got %q", want, got)
	}
	if got, want := p.token, ""; got != want {
		t.Errorf("Want token %q, got %q", want, got)
	}
	if got, want := len(p.tags), 0; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
}
