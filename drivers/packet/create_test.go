// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package packet

import (
	"context"
	"testing"
	"time"

	"github.com/drone/autoscaler"
	"github.com/h2non/gock"
	"github.com/packethost/packngo"
)

func TestCreate(t *testing.T) {
	defer gock.Off()

	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Post(createDevice).
		Reply(200).
		BodyString(respCreate)

	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Get(getDevice).
		MatchParam("include", "facility").
		Reply(200).
		Delay(10 * time.Second).
		BodyString(respCreate)

	instance, err := prov.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: prov.os})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Run("Attributes", testInstance(instance))
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, instanceID; got != want {
			t.Errorf("Want ID %v, got %v", want, got)
		}
		if got, want := instance.Image, prov.os; got != want {
			t.Errorf("Want Image %v, got %v", want, got)
		}
		if got, want := instance.Name, prov.os; got != want {
			t.Errorf("Want Name %v, got %v", want, got)
		}
		if got, want := instance.Region, prov.facility; got != want {
			t.Errorf("Want Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderPacket; got != want {
			t.Errorf("Want Provider %v, got %v", want, got)
		}
	}
}

func TestCreate_Timeout(t *testing.T) {
	defer gock.Off()

	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Post(getDevice).
		Reply(200).
		BodyString(respCreateInactive)
	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Get(getDevice).
		Reply(200).
		BodyString(respCreateInactive)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := prov.Create(ctx, autoscaler.InstanceCreateOpts{Name: prov.os}); err != context.Canceled {
		t.Errorf("Expected error creating a device")
	}
}

func TestCreate_Erro(t *testing.T) {
	defer gock.Off()

	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Post(getDevice).
		Reply(400)

	_, err := prov.Create(context.Background(), autoscaler.InstanceCreateOpts{Name: prov.os})
	if err == nil {
		t.Errorf("Expect error returned when creatiung the device")
	} else if _, ok := err.(*packngo.ErrorResponse); !ok {
		t.Errorf("Expect error to be of type  ErrorResponse")
	}
}

func TestCreate_WaitToBecomeActive(t *testing.T) {
	defer gock.Off()

	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Post(getDevice).
		Reply(200).
		BodyString(respCreateInactive)
	gock.New(baseURL).
		MatchHeader("X-Auth-Token", apiKey).
		Get(getDevice).
		Reply(200).
		BodyString(respCreateInactive)

	wait := make(chan struct{})
	go func() {
		prov.Create(context.Background(), autoscaler.InstanceCreateOpts{Name: prov.os})
		close(wait)
	}()

	select {
	case <-wait:
		t.Errorf("Expected device creation to block when device is not set to active")
	case <-time.After(50 * time.Millisecond):
	}
}
