// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/drone/autoscaler/drivers/internal/userdata"
)

func TestDefaults(t *testing.T) {
	v, err := New(
		WithClient(http.DefaultClient),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)

	if got, want := p.diskSize, int64(50); got != want {
		t.Errorf("Want diskSize %d, got %d", want, got)
	}
	if got, want := p.diskType, "pd-standard"; got != want {
		t.Errorf("Want diskType %s, got %s", want, got)
	}
	if got, want := p.image, "ubuntu-os-cloud/global/images/ubuntu-2004-focal-v20220110"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.network, "global/networks/default"; got != want {
		t.Errorf("Want network %q, got %q", want, got)
	}
	if !reflect.DeepEqual(p.scopes, defaultScopes) {
		t.Errorf("Want default scopes")
	}
	if got, want := p.size, "n1-standard-1"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if !reflect.DeepEqual(p.tags, defaultTags) {
		t.Errorf("Want default tags")
	}
	if p.userdata != userdata.T {
		t.Errorf("Want default userdata template")
	}
	if p.userdataKey != "user-data" {
		t.Errorf("Want default userdata key")
	}
	if got, want := p.zones, []string{"us-central1-a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Want region %q, got %q", want, got)
	}
}
