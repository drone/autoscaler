// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package openstack

import (
	"context"
	"github.com/drone/autoscaler"
	"github.com/h2non/gock"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
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

	fipResp1 := helperLoad(t, "fipresp1.json")
	gock.New("http://ops.my.cloud").
		Post("/compute/v2.1/os-floating-ips").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(fipResp1))

	imageListResp := helperLoad(t, "imagelistresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/compute/v2.1/images/detail").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(imageListResp))

	flavorListResp1 := helperLoad(t, "flavorlistresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/compute/v2.1/flavors/detail").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(flavorListResp1))

	serverCreateResp1 := helperLoad(t, "servercreateresp1.json")
	gock.New("http://ops.my.cloud").
		Post("/compute/v2.1/servers").
		MatchHeader("X-Auth-Token", authToken).
		Reply(202).
		SetHeader("Content-Type", "application/json").
		BodyString(string(serverCreateResp1))

	serverStatusResp1 := helperLoad(t, "serverstatusresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/compute/v2.1/servers/56046f6d-3184-495b-938b-baa450db970d").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(serverStatusResp1))

	associateResp1 := helperLoad(t, "associateresp1.json")
	gock.New("http://ops.my.cloud").
		Post("/compute/v2.1/servers/56046f6d-3184-495b-938b-baa450db970d/action").
		MatchHeader("X-Auth-Token", authToken).
		BodyString(string("{\"addFloatingIp\":{\"address\":\"172.24.4.5\"}}")).
		Reply(202).
		SetHeader("Content-Type", "application/json").
		BodyString(string(associateResp1))

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
	p.init.Do(func() {}) // prevent init function

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent-RjISb5v1"})
	if err != nil {
		t.Error(err)
	}

	if !gock.IsDone() {
		t.Error("Not all expected http requests completed")
	}
	t.Run("Instance Attributes", testInstance(instance))
}

func TestAuthFail(t *testing.T) {
	defer gock.Off()
	setupEnv(t)

	err := os.Setenv("OS_PASSWORD", "BAADF00D")
	if err != nil {
		t.Error("Unable to set OS_PASSWORD")
	}
	authResp1 := helperLoad(t, "authresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/identity").
		Reply(300).
		SetHeader("Content-Type", "application/json").
		BodyString(string(authResp1))

	gock.New("http://ops.my.cloud").
		Post("/identity/v3/auth/tokens").
		Reply(401)

	_, err = New(
		WithRegion("RegionOne"),
		WithFlavor("m1.small"),
		WithImage("ubuntu-16.04-server-latest"),
		WithFloatingIpPool("public"),
		WithSSHKey("drone-ci-key"),
	)

	if err == nil {
		t.Error("Expected authentication error from OpenStack")
	}

	if !gock.IsDone() {
		t.Error("Not all expected http requests completed")
	}
}

func TestCreateFail(t *testing.T) {
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

	fipResp1 := helperLoad(t, "fipresp1.json")
	gock.New("http://ops.my.cloud").
		Post("/compute/v2.1/os-floating-ips").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(fipResp1))

	imageListResp := helperLoad(t, "imagelistresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/compute/v2.1/images/detail").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(imageListResp))

	flavorListResp1 := helperLoad(t, "flavorlistresp1.json")
	gock.New("http://ops.my.cloud").
		Get("/compute/v2.1/flavors/detail").
		MatchHeader("X-Auth-Token", authToken).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(string(flavorListResp1))

	gock.New("http://ops.my.cloud").
		Post("/compute/v2.1/servers").
		MatchHeader("X-Auth-Token", authToken).
		Reply(500)

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
	p.init.Do(func() {}) // prevent init function

	_, err = p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent-RjISb5v1"})
	if err == nil {
		t.Error("Expected error creating instance")
	}

	if !gock.IsDone() {
		t.Error("Not all expected http requests completed")
	}
}

func setupEnv(t *testing.T) {
	err := os.Setenv("OS_AUTH_URL", "http://ops.my.cloud/identity")
	if err != nil {
		t.Error("Unable to set OS_AUTH_URL")
	}
	err = os.Setenv("OS_USERNAME", "admin")
	if err != nil {
		t.Error("Unable to set OS_USERNAME")
	}
	err = os.Setenv("OS_PASSWORD", "admin")
	if err != nil {
		t.Error("Unable to set OS_USERNAME")
	}
	err = os.Setenv("OS_DOMAIN_NAME", "demo")
	if err != nil {
		t.Error("Unable to set OS_DOMAIN_NAME")
	}
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if want, got := instance.Address, "172.24.4.5"; got != want {
			t.Errorf("Want instance IP %q, got %q", want, got)
		}
		if want, got := instance.Image, "ubuntu-16.04-server-latest"; got != want {
			t.Errorf("Want instance ID %q, got %q", want, got)
		}
		if want, got := instance.ID, "56046f6d-3184-495b-938b-baa450db970d"; got != want {
			t.Errorf("Want instance ID %q, got %q", want, got)
		}
		if want, got := instance.Name, "agent-RjISb5v1"; got != want {
			t.Errorf("Want instance Name %q, got %q", want, got)
		}
		if want, got := instance.Provider, autoscaler.ProviderOpenStack; got != want {
			t.Errorf("Want OpenStack Provider type")
		}
		if want, got := instance.Region, "RegionOne"; got != want {
			t.Errorf("Want instance Region %q, got %q", want, got)
		}
		if want, got := instance.Size, "m1.small"; got != want {
			t.Errorf("Want instance Size %q, got %q", want, got)
		}
	}
}
