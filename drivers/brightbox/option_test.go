// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package brightbox

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithApiURL("https://api.gb1.brightbox.com"),
		WithClientID("cli-xxxxx"),
		WithClientSecret("supersecret"),
		WithImage("img-xxxxx"),
		WithServerType("typ-xxxxx"),
		WithServerGroups([]string{"grp-aaaaa"}),
	).(*provider)

	if got, want := p.apiURL, "https://api.gb1.brightbox.com"; got != want {
		t.Errorf("Want API URL %q, got %q", want, got)
	}
	if got, want := p.clientID, "cli-xxxxx"; got != want {
		t.Errorf("Want client ID %q, got %q", want, got)
	}
	if got, want := p.clientSecret, "supersecret"; got != want {
		t.Errorf("Want client secret %q, got %q", want, got)
	}
	if got, want := p.image, "img-xxxxx"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.serverType, "typ-xxxxx"; got != want {
		t.Errorf("Want server type %q, got %q", want, got)
	}
	if got, want := p.serverGroups[0], "grp-aaaaa"; got != want {
		t.Errorf("Want server groups %q, got %q", want, got)
	}
}
