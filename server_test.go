// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import (
	"sort"
	"strings"
	"testing"
)

func TestByCreated(t *testing.T) {
	servers := []*Server{
		{Created: 4, Name: "fourth"},
		{Created: 2, Name: "second"},
		{Created: 3, Name: "third"},
		{Created: 5, Name: "fifth"},
		{Created: 1, Name: "first"},
	}

	sort.Sort(ByCreated(servers))

	for i, server := range servers {
		if server.Created != int64(i+1) {
			t.Errorf("Invalid sort order %d for %q", i, server.Name)
		}
	}
}

func TestServerOpts(t *testing.T) {
	opts := NewServerOpts("agent", 4)
	if got, want := opts.Capacity, 4; got != want {
		t.Errorf("Want capacity %d, got %d", want, got)
	}
	if !strings.HasPrefix(opts.Name, "agent-") {
		t.Errorf("Want server name prefixed with %s, got %s", "agent-", opts.Name)
	}
	if got, want := len(opts.Name), len("agent-")+5; got != want {
		t.Errorf("Want server name length %d, got %d", want, got)
	}
}
