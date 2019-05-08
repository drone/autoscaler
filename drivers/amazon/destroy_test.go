// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"context"
	"os"
	"testing"

	"github.com/drone/autoscaler"

	"github.com/h2non/gock"
)

func TestDestroy(t *testing.T) {
	defer gock.Off()

	os.Setenv("AWS_ACCESS_KEY_ID", "your_access_key_id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "your_secret_access_key")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		Reply(200)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "i-1234567890abcdef0",
	}

	p := New(
		WithRegion("us-east-1"),
	).(*provider)
	p.retries = 1

	err := p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}

func TestDestroyDeleteError(t *testing.T) {
	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "i-1234567890abcdef0",
	}

	p := New(
		WithRegion("us-east-1"),
	).(*provider)
	p.retries = 1

	err := p.Destroy(mockContext, mockInstance)
	if err == nil {
		t.Errorf("Expect error returned from aws")
	}
}

func TestDestroyNotFound(t *testing.T) {
	defer gock.Off()

	os.Setenv("AWS_ACCESS_KEY_ID", "your_access_key_id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "your_secret_access_key")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		Reply(400).
		BodyString(`<Response><Errors><Error><Code>InvalidInstanceID.NotFound</Code><Message>The instance ID 'i-1a2b3c4d' does not exist</Message></Error></Errors><RequestID>ea966190-f9aa-478e-9ede-example</RequestID></Response>`)

	mockContext := context.TODO()
	mockInstance := &autoscaler.Instance{
		ID: "i-1234567890abcdef0",
	}

	p := New(
		WithRegion("us-east-1"),
	).(*provider)
	p.retries = 1

	err := p.Destroy(mockContext, mockInstance)
	if err == nil {
		t.Errorf("Expect error returned from aws")
	}
	if err != autoscaler.ErrInstanceNotFound {
		t.Errorf("Expect instance not found returned from aws")
	}
}
