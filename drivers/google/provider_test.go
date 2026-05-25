// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package google

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/h2non/gock"
)

func TestDefaults(t *testing.T) {
	v, err := New(
		WithClient(http.DefaultClient),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)

	if got, want := p.diskSize, int64(50); got != want {
		t.Errorf("Want diskSize %d, got %d", want, got)
	}
	if got, want := p.diskType, "pd-standard"; got != want {
		t.Errorf("Want diskType %s, got %s", want, got)
	}
	if got, want := p.image, "ubuntu-os-cloud/global/images/ubuntu-2004-focal-v20220712"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.network, "global/networks/default"; got != want {
		t.Errorf("Want network %q, got %q", want, got)
	}
	if !reflect.DeepEqual(p.scopes, defaultScopes) {
		t.Errorf("Want default scopes")
	}
	if got, want := p.size, "n1-standard-1"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if !reflect.DeepEqual(p.tags, defaultTags) {
		t.Errorf("Want default tags")
	}
	if p.userdata != userdata.T {
		t.Errorf("Want default userdata template")
	}
	if p.userdataKey != "user-data" {
		t.Errorf("Want default userdata key")
	}
	if got, want := p.zones, []string{"us-central1-a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Want region %q, got %q", want, got)
	}
}

func TestWaitZoneOperationSuccess(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitZoneOperation(context.Background(), "op-123", "us-central1-a")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWaitZoneOperationTransientErrorRetry(t *testing.T) {
	defer gock.Off()

	// First call returns 503 Service Unavailable (transient error)
	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Reply(503).
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    503,
				"message": "Service Unavailable",
			},
		})

	// Subsequent calls succeed
	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Times(5).
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitZoneOperation(context.Background(), "op-123", "us-central1-a")
	if err != nil {
		t.Errorf("Expected no error after retry, got %v", err)
	}
}

func TestWaitZoneOperationNonTransientError(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Reply(400).
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    400,
				"message": "Bad Request",
			},
		})

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitZoneOperation(context.Background(), "op-123", "us-central1-a")
	if err == nil {
		t.Errorf("Expected error for non-transient error, got nil")
	}
}

func TestWaitZoneOperationNotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Reply(404).
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": "Not Found",
			},
		})

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitZoneOperation(context.Background(), "op-123", "us-central1-a")
	if err != autoscaler.ErrInstanceNotFound {
		t.Errorf("Expected ErrInstanceNotFound, got %v", err)
	}
}

func TestWaitGlobalOperationSuccess(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/operations/op-123").
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitGlobalOperation(context.Background(), "op-123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWaitGlobalOperationTransientErrorRetry(t *testing.T) {
	defer gock.Off()

	// First call returns 503 Service Unavailable (transient error)
	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/operations/op-123").
		Reply(503).
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    503,
				"message": "Service Unavailable",
			},
		})

	// Subsequent calls succeed
	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/operations/op-123").
		Times(5).
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitGlobalOperation(context.Background(), "op-123")
	if err != nil {
		t.Errorf("Expected no error after retry, got %v", err)
	}
}

func TestWaitZoneOperationRateLimitRetry(t *testing.T) {
	defer gock.Off()

	// First call returns 429 with Retry-After: 0 (retry immediately)
	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Reply(429).
		AddHeader("Retry-After", "0").
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    429,
				"message": "Too Many Requests",
			},
		})

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/zones/us-central1-a/operations/op-123").
		Reply(200).
		BodyString(`{ "status": "DONE" }`)

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
		WithRateLimit(1000),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitZoneOperation(context.Background(), "op-123", "us-central1-a")
	if err != nil {
		t.Errorf("Expected no error after rate limit retry, got %v", err)
	}
}

func TestWaitGlobalOperationNonTransientError(t *testing.T) {
	defer gock.Off()

	gock.New("https://compute.googleapis.com").
		Get("/compute/v1/projects/my-project/global/operations/op-123").
		Reply(400).
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    400,
				"message": "Bad Request",
			},
		})

	v, err := New(
		WithClient(http.DefaultClient),
		WithProject("my-project"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	err = p.waitGlobalOperation(context.Background(), "op-123")
	if err == nil {
		t.Errorf("Expected error for non-transient error, got nil")
	}
}
