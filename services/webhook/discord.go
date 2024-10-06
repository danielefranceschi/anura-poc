// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	webhook_model "code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
)

type (
	// DiscordEmbedFooter for Embed Footer Structure.
	DiscordEmbedFooter struct {
		Text string `json:"text"`
	}

	// DiscordEmbedAuthor for Embed Author Structure
	DiscordEmbedAuthor struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		IconURL string `json:"icon_url"`
	}

	// DiscordEmbedField for Embed Field Structure
	DiscordEmbedField struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// DiscordEmbed is for Embed Structure
	DiscordEmbed struct {
		Title       string              `json:"title"`
		Description string              `json:"description"`
		URL         string              `json:"url"`
		Color       int                 `json:"color"`
		Footer      DiscordEmbedFooter  `json:"footer"`
		Author      DiscordEmbedAuthor  `json:"author"`
		Fields      []DiscordEmbedField `json:"fields"`
	}

	// DiscordPayload represents
	DiscordPayload struct {
		Wait      bool           `json:"wait"`
		Content   string         `json:"content"`
		Username  string         `json:"username"`
		AvatarURL string         `json:"avatar_url,omitempty"`
		TTS       bool           `json:"tts"`
		Embeds    []DiscordEmbed `json:"embeds"`
	}

	// DiscordMeta contains the discord metadata
	DiscordMeta struct {
		Username string `json:"username"`
		IconURL  string `json:"icon_url"`
	}
)

// GetDiscordHook returns discord metadata
func GetDiscordHook(w *webhook_model.Webhook) *DiscordMeta {
	s := &DiscordMeta{}
	if err := json.Unmarshal([]byte(w.Meta), s); err != nil {
		log.Error("webhook.GetDiscordHook(%d): %v", w.ID, err)
	}
	return s
}

func color(clr string) int {
	if clr != "" {
		clr = strings.TrimLeft(clr, "#")
		if s, err := strconv.ParseInt(clr, 16, 32); err == nil {
			return int(s)
		}
	}

	return 0
}

var (
	greenColor       = color("1ac600")
	greenColorLight  = color("bfe5bf")
	yellowColor      = color("ffd930")
	greyColor        = color("4f545c")
	purpleColor      = color("7289da")
	orangeColor      = color("eb6420")
	orangeColorLight = color("e68d60")
	redColor         = color("ff3232")
)

type discordConvertor struct {
	Username  string
	AvatarURL string
}

// Repository implements PayloadConvertor Repository method
// func (d discordConvertor) Repository(p *api.RepositoryPayload) (DiscordPayload, error) {
// 	var title, url string
// 	var color int
// 	switch p.Action {
// 	case api.HookRepoCreated:
// 		title = fmt.Sprintf("[%s] Repository created", p.Repository.FullName)
// 		url = p.Repository.HTMLURL
// 		color = greenColor
// 	case api.HookRepoDeleted:
// 		title = fmt.Sprintf("[%s] Repository deleted", p.Repository.FullName)
// 		color = redColor
// 	}

// 	return d.createPayload(p.Sender, title, "", url, color), nil
// }

func (d discordConvertor) Package(p *api.PackagePayload) (DiscordPayload, error) {
	text, color := getPackagePayloadInfo(p, noneLinkFormatter, false)

	return d.createPayload(p.Sender, text, "", p.Package.HTMLURL, color), nil
}

func newDiscordRequest(_ context.Context, w *webhook_model.Webhook, t *webhook_model.HookTask) (*http.Request, []byte, error) {
	meta := &DiscordMeta{}
	if err := json.Unmarshal([]byte(w.Meta), meta); err != nil {
		return nil, nil, fmt.Errorf("newDiscordRequest meta json: %w", err)
	}
	var pc payloadConvertor[DiscordPayload] = discordConvertor{
		Username:  meta.Username,
		AvatarURL: meta.IconURL,
	}
	return newJSONRequest(pc, w, t, true)
}

func (d discordConvertor) createPayload(s *api.User, title, text, url string, color int) DiscordPayload {
	return DiscordPayload{
		Username:  d.Username,
		AvatarURL: d.AvatarURL,
		Embeds: []DiscordEmbed{
			{
				Title:       title,
				Description: text,
				URL:         url,
				Color:       color,
				Author: DiscordEmbedAuthor{
					Name:    s.UserName,
					URL:     setting.AppURL + s.UserName,
					IconURL: s.AvatarURL,
				},
			},
		},
	}
}
