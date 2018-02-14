// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import "testing"

func TestGenerateSha(t *testing.T) {
	got, want := generateSha("correct-horse-batter-stample"), "y2z/A6MK2bl1XJnCRmTSRq2TTC4="
	if got != want {
		t.Errorf("Want sha %s, got %s", want, got)
	}
}
