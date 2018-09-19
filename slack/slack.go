// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package slack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"

	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
)

// New returns a new provider that is instrumented to send
// Slack notifications when server instances are provisioned
// or terminated.
func New(config config.Config, base autoscaler.ServerStore) autoscaler.ServerStore {
	return &notifier{
		ServerStore: base,
		client:      slack.NewWebHook(config.Slack.Webhook),
		create:      config.Slack.Create,
		destroy:     config.Slack.Destroy,
		error:       config.Slack.Error,
	}
}

type notifier struct {
	autoscaler.ServerStore
	client  *slack.WebHook
	channel string
	create  bool
	destroy bool
	error   bool
}

func (n *notifier) Update(ctx context.Context, server *autoscaler.Server) error {
	err := n.ServerStore.Update(ctx, server)
	switch {
	case server.State == autoscaler.StateRunning && n.create:
		n.notifyCreate(server)
	case server.State == autoscaler.StateStopped && n.destroy:
		n.notifyDestroy(server)
	case server.State == autoscaler.StateError && n.error:
		n.notifyError(server)
	}
	return err
}

func (n *notifier) notifyCreate(server *autoscaler.Server) error {
	opts := &slack.WebHookPostPayload{
		Text: fmt.Sprintf("Provisioned server instance %s", server.Name),
		Attachments: []*slack.Attachment{
			{
				Color: "#00BFA5",
				Fields: []*slack.AttachmentField{
					{
						Title: "Name",
						Value: server.Name,
						Short: false,
					},
					{
						Title: "Size",
						Value: server.Size,
						Short: false,
					},
					{
						Title: "Region",
						Value: server.Region,
						Short: false,
					},
				},
			},
		},
	}
	return n.client.PostMessage(opts)
}

func (n *notifier) notifyDestroy(server *autoscaler.Server) error {
	opts := &slack.WebHookPostPayload{
		Text: fmt.Sprintf("Terminated server instance %s", server.Name),
		Attachments: []*slack.Attachment{
			{
				Color: "#CFD8DC",
				Fields: []*slack.AttachmentField{
					{
						Title: "Name",
						Value: server.Name,
						Short: false,
					},
					{
						Title: "Size",
						Value: server.Size,
						Short: false,
					},
					{
						Title: "Region",
						Value: server.Region,
						Short: false,
					},
					{
						Title: "Uptime",
						Value: humanizeTime(server.Created),
						Short: false,
					},
				},
			},
		},
	}
	return n.client.PostMessage(opts)
}

func (n *notifier) notifyError(server *autoscaler.Server) error {
	opts := &slack.WebHookPostPayload{
		Text: fmt.Sprintf("Problem with server instance %s", server.Name),
		Attachments: []*slack.Attachment{
			{
				Color: "#F44336",
				Fields: []*slack.AttachmentField{
					{
						Title: "Name",
						Value: server.Name,
						Short: false,
					},
					{
						Title: "Error",
						Value: server.Error,
						Short: false,
					},
				},
			},
		},
	}
	return n.client.PostMessage(opts)
}

func humanizeTime(unix int64) string {
	d := time.Unix(unix, 0)
	s := humanize.RelTime(d, time.Now(), "", "")
	return strings.TrimSpace(s)
}
