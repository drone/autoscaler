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
func New(config config.Config, base autoscaler.Provider) autoscaler.Provider {
	return &notifier{
		Provider: base,
		client:   slack.NewWebHook(config.Slack.Webhook),
	}
}

type notifier struct {
	autoscaler.Provider
	client  *slack.WebHook
	channel string
}

func (n *notifier) Create(ctx context.Context, opts *autoscaler.ServerOpts) (*autoscaler.Server, error) {
	server, err := n.Provider.Create(ctx, opts)
	if err == nil {
		n.notifyCreate(server)
	}
	return server, err
}

func (n *notifier) Destroy(ctx context.Context, server *autoscaler.Server) error {
	err := n.Provider.Destroy(ctx, server)
	if err == nil {
		n.notifyDestroy(server)
	}
	return err
}

func (n *notifier) notifyCreate(server *autoscaler.Server) error {
	opts := &slack.WebHookPostPayload{
		Text: fmt.Sprintf("Provisioned server instance %s", server.Name),
		Attachments: []*slack.Attachment{
			{
				Color: "good",
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
				Color: "danger",
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

func humanizeTime(unix int64) string {
	d := time.Unix(unix, 0)
	s := humanize.RelTime(d, time.Now(), "", "")
	return strings.TrimSpace(s)
}
