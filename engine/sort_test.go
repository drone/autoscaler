// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import (
	"sort"
	"testing"

	"github.com/drone/autoscaler"
)

func TestSortByCreated(t *testing.T) {
	servers := []*autoscaler.Server{
		{Created: 4, Name: "fourth"},
		{Created: 2, Name: "second"},
		{Created: 3, Name: "third"},
		{Created: 5, Name: "fifth"},
		{Created: 1, Name: "first"},
	}

	sort.Sort(byCreated(servers))

	for i, server := range servers {
		if server.Created != int64(i+1) {
			t.Errorf("Invalid sort order %d for %q", i, server.Name)
		}
	}
}
