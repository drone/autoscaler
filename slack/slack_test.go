// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package slack

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/mocks"

	"github.com/bluele/slack"
	"github.com/golang/mock/gomock"
	"github.com/h2non/gock"
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

	gock.New("https://hooks.slack.com").
		Post("/services/XXX/YYY/ZZZ").
		JSON(createPayload).
		Reply(200)

	conf := config.Config{}
	conf.Slack.Webhook = "https://hooks.slack.com/services/XXX/YYY/ZZZ"
	conf.Slack.Create = true

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	slack := New(conf, store)
	err := slack.Update(noContext, server)
	if err != nil {
		t.Error(err)
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks not executed")
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

	gock.New("https://hooks.slack.com").
		Post("/services/XXX/YYY/ZZZ").
		Reply(200)

	conf := config.Config{}
	conf.Slack.Webhook = "https://hooks.slack.com/services/XXX/YYY/ZZZ"
	conf.Slack.Destroy = true

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	slack := New(conf, store)
	err := slack.Update(noContext, server)
	if err != nil {
		t.Error(err)
		return
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks not executed")
	}
}

func TestUpdateError(t *testing.T) {
	defer gock.Off()

	controller := gomock.NewController(t)
	defer controller.Finish()

	server := &autoscaler.Server{
		Name:   "this-is-a-test-message",
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
		Error:  "pc load letter",
		State:  autoscaler.StateError,
	}

	gock.New("https://hooks.slack.com").
		Post("/services/XXX/YYY/ZZZ").
		JSON(errorPayload).
		Reply(200)

	conf := config.Config{}
	conf.Slack.Webhook = "https://hooks.slack.com/services/XXX/YYY/ZZZ"
	conf.Slack.Error = true

	store := mocks.NewMockServerStore(controller)
	store.EXPECT().Update(gomock.Any(), server).Return(nil)

	slack := New(conf, store)
	err := slack.Update(noContext, server)
	if err != nil {
		t.Error(err)
	}

	if !gock.IsDone() {
		t.Errorf("Pending mocks not executed")
	}
}

var createPayload = slack.WebHookPostPayload{
	Text: "Provisioned server instance this-is-a-test-message",
	Attachments: []*slack.Attachment{
		{
			Color: "#00BFA5",
			Fields: []*slack.AttachmentField{
				{
					Title: "Name",
					Value: "this-is-a-test-message",
				},
				{
					Title: "Size",
					Value: "s-1vcpu-1gb",
				},
				{
					Title: "Region",
					Value: "nyc1",
				},
			},
		},
	},
}

var errorPayload = slack.WebHookPostPayload{
	Text: "Problem with server instance this-is-a-test-message",
	Attachments: []*slack.Attachment{
		{
			Color: "#F44336",
			Fields: []*slack.AttachmentField{
				{
					Title: "Name",
					Value: "this-is-a-test-message",
				},
				{
					Title: "Error",
					Value: "pc load letter",
				},
			},
		},
	},
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
