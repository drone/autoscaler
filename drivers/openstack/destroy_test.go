// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"context"
	"github.com/drone/autoscaler"
	"github.com/h2non/gock"
	"testing"
)

func TestDestroy(t *testing.T) {
	defer gock.Off()
	setupEnv(t)

	authResp1 := helperLoad(t, "authresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/identity").
		Reply(300).
		SetHeader("Content-Type", "application/json").
		BodyString(string(authResp1))

	tokenResp1 := helperLoad(t, "tokenresp1.json")
	gock.New("http://ops.my.cloud").
		Post("/identity/v3/auth/tokens").
		Reply(201).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Subject-Token", authToken).
		BodyString(string(tokenResp1))

	gock.New("http://ops.my.cloud").
		MatchHeader("X-Auth-Token", authToken).
		Post("/compute/v2.1/servers/56046f6d-3184-495b-938b-baa450db970d/action").
		Reply(202)

	fipResp1 := helperLoad(t, "fipresp1.json")
	gock.New("http://ops.my.cloud").
		MatchHeader("X-Auth-Token", authToken).
		Get("/compute/v2.1/os-floating-ips").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(fipResp1))

	gock.New("http://ops.my.cloud").
		MatchHeader("X-Auth-Token", authToken).
		Delete("/compute/v2.1/servers/56046f6d-3184-495b-938b-baa450db970d").
		Reply(204)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID:      "56046f6d-3184-495b-938b-baa450db970d",
		Address: "172.24.4.5",
	}

	v, err := New(
		WithRegion("RegionOne"),
		WithFlavor("m1.small"),
		WithImage("ubuntu-16.04-server-latest"),
		WithFloatingIpPool("public"),
		WithSSHKey("drone-ci-key"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {}) //

	err = p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}

	if !gock.IsDone() {
		t.Error("Not all expected http requests completed")
	}
}
