// Copyright 2014 The Gogs Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	webhook_model "code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/log"
	api "code.gitea.io/gitea/modules/structs"
)

// SlackMeta contains the slack metadata
type SlackMeta struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
	Color    string `json:"color"`
}

// GetSlackHook returns slack metadata
func GetSlackHook(w *webhook_model.Webhook) *SlackMeta {
	s := &SlackMeta{}
	if err := json.Unmarshal([]byte(w.Meta), s); err != nil {
		log.Error("webhook.GetSlackHook(%d): %v", w.ID, err)
	}
	return s
}

// SlackPayload contains the information about the slack channel
type SlackPayload struct {
	Channel     string            `json:"channel"`
	Text        string            `json:"text"`
	Username    string            `json:"username"`
	IconURL     string            `json:"icon_url"`
	UnfurlLinks int               `json:"unfurl_links"`
	LinkNames   int               `json:"link_names"`
	Attachments []SlackAttachment `json:"attachments"`
}

// SlackAttachment contains the slack message
type SlackAttachment struct {
	Fallback  string `json:"fallback"`
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
}

// SlackTextFormatter replaces &, <, > with HTML characters
// see: https://api.slack.com/docs/formatting
func SlackTextFormatter(s string) string {
	// replace & < >
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// SlackShortTextFormatter replaces &, <, > with HTML characters
func SlackShortTextFormatter(s string) string {
	s = strings.Split(s, "\n")[0]
	// replace & < >
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// SlackLinkFormatter creates a link compatible with slack
func SlackLinkFormatter(url, text string) string {
	return fmt.Sprintf("<%s|%s>", url, SlackTextFormatter(text))
}

func (s slackConvertor) Package(p *api.PackagePayload) (SlackPayload, error) {
	text, _ := getPackagePayloadInfo(p, SlackLinkFormatter, true)

	return s.createPayload(text, nil), nil
}

func (s slackConvertor) createPayload(text string, attachments []SlackAttachment) SlackPayload {
	return SlackPayload{
		Channel:     s.Channel,
		Text:        text,
		Username:    s.Username,
		IconURL:     s.IconURL,
		Attachments: attachments,
	}
}

type slackConvertor struct {
	Channel  string
	Username string
	IconURL  string
	Color    string
}

func newSlackRequest(_ context.Context, w *webhook_model.Webhook, t *webhook_model.HookTask) (*http.Request, []byte, error) {
	meta := &SlackMeta{}
	if err := json.Unmarshal([]byte(w.Meta), meta); err != nil {
		return nil, nil, fmt.Errorf("newSlackRequest meta json: %w", err)
	}
	var pc payloadConvertor[SlackPayload] = slackConvertor{
		Channel:  meta.Channel,
		Username: meta.Username,
		IconURL:  meta.IconURL,
		Color:    meta.Color,
	}
	return newJSONRequest(pc, w, t, true)
}

var slackChannel = regexp.MustCompile(`^#?[a-z0-9_-]{1,80}$`)

// IsValidSlackChannel validates a channel name conforms to what slack expects:
// https://api.slack.com/methods/conversations.rename#naming
// Conversation names can only contain lowercase letters, numbers, hyphens, and underscores, and must be 80 characters or less.
// Gitea accepts if it starts with a #.
func IsValidSlackChannel(name string) bool {
	return slackChannel.MatchString(name)
}
