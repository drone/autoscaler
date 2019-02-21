// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"github.com/gophercloud/gophercloud"
	"testing"
)

func TestDefaults(t *testing.T) {
	v, err := New(
		WithComputeClient(&gophercloud.ServiceClient{}),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	// Add tests if we set some actual defaults in the future.
	_ = p
}
