// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"net/http"
	"testing"
)

func TestOptions(t *testing.T) {
	v, err := New(
		WithClient(http.DefaultClient),
		WithDiskSize(100),
		WithDiskType("local-ssd"),
		WithMachineImage("ubuntu-1604-lts"),
		WithMachineType("c3.large"),
		WithNetwork("global/defaults/foo"),
		WithPrivateIP(false),
		WithProject("my-project"),
		WithTags("drone", "agent"),
		WithZone("us-central1-f"),
		WithScopes("scope1,scope2"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)

	if got, want := p.diskSize, int64(100); got != want {
		t.Errorf("Want diskSize %d, got %d", want, got)
	}
	if got, want := p.diskType, "local-ssd"; got != want {
		t.Errorf("Want diskType %s, got %s", want, got)
	}
	if got, want := p.image, "ubuntu-1604-lts"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.network, "global/defaults/foo"; got != want {
		t.Errorf("Want network %q, got %q", want, got)
	}
	if got, want := p.privateIP, false; got != want {
		t.Errorf("Want %v privateIP, got %v", want, got)
	}
	if got, want := p.project, "my-project"; got != want {
		t.Errorf("Want project %q, got %q", want, got)
	}
	if got, want := p.size, "c3.large"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := len(p.tags), 2; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
	if got, want := p.zone, "us-central1-f"; got != want {
		t.Errorf("Want zone %q, got %q", want, got)
	}
	if got, want := len(p.scopes), 2; got != want {
		t.Errorf("Want %d scopes, got %d", want, got)
	}
}
