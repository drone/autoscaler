// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func helperLoad(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

const authToken = "gAAAAABb1tQPtYVBv68airR0dgKC2vXpkLNfEHx0w1EL89dOOjKrtdYHR7IZrDd4VjwZapC5Sri4CndpPscw-nHoh0VQsrvFjtuvT6M64RdrrOljmJbvP0o7PbV713-Pi8OpRIfunvsQFnEQ2DxDH56QC6fsLEcF14VtogOQwTRBod0SkeOCpi4"