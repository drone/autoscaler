// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package slack

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/h2non/gock"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/mocks"

	"github.com/golang/mock/gomock"
)

var noContext = context.TODO()

func TestHumanizeTime(t *testing.T) {
	unix := time.Now().Add(time.Minute * 60 * -1).Unix()
	text := humanizeTime(unix)
	if got, want := text, "1 hour"; got != want {
		t.Errorf("Want humanized time %s, got %s", want, got)
	}
}

func TestUpdateRunning(t *testing.T) {
	defer gock.Off()

	controller := gomock.NewController(t)
	defer controller.Finish()

	server := &autoscaler.Server{
		Name:   "this-is-a-test-message",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
		State:  autoscaler.StateRunning,
	}

	// TODO: verify the contents of the Slack payload.

	gock.New("https://hooks.slack.com").
		Get("/services/XXX/YYY/ZZZ").
		Reply(200)

	conf := config.Config{}
	conf.Slack.Webhook = "https://hooks.slack.com/services/XXX/YYY/ZZZ"

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	slack := New(conf, store)
	err := slack.Update(noContext, server)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateStopped(t *testing.T) {
	defer gock.Off()

	controller := gomock.NewController(t)
	defer controller.Finish()

	server := &autoscaler.Server{
		Name:   "this-is-a-test-message",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
		State:  autoscaler.StateStopped,
	}

	// TODO: verify the contents of the Slack payload.

	gock.New("https://hooks.slack.com").
		Get("/services/XXX/YYY/ZZZ").
		Reply(200)

	conf := config.Config{}
	conf.Slack.Webhook = "https://hooks.slack.com/services/XXX/YYY/ZZZ"

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	slack := New(conf, store)
	err := slack.Update(noContext, server)
	if err != nil {
		t.Error(err)
	}
}

// This is an integration test that will send a real
// message to a Slack channel using a webhook defined
// in the TEST_SLACK_WEBHOOK environment variable.
func TestIntegration(t *testing.T) {
	webhook := os.Getenv("TEST_SLACK_WEBHOOK")
	if webhook == "" {
		t.Skipf("Skip Slack integration test. No webhook provided.")
		return
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	server := &autoscaler.Server{
		Name:     "i-123789331",
		Address:  "1.2.3.4",
		Region:   "nyc1",
		Size:     "s-1vcpu-1gb",
		Capacity: 2,
		State:    autoscaler.StateRunning,
	}

	conf := config.Config{}
	conf.Slack.Webhook = webhook

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	slack := New(conf, store)
	err := slack.Update(noContext, server)
	if err != nil {
		t.Error(err)
	}
}
