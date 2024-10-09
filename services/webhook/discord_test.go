// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"context"
	"testing"

	webhook_model "code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/setting"
	webhook_module "code.gitea.io/gitea/modules/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscordPayload(t *testing.T) {
	dc := discordConvertor{}

	t.Run("Package", func(t *testing.T) {
		p := packageTestPayload()

		pl, err := dc.Package(p)
		require.NoError(t, err)

		assert.Len(t, pl.Embeds, 1)
		assert.Equal(t, "Package created: GiteaContainer:latest", pl.Embeds[0].Title)
		assert.Empty(t, pl.Embeds[0].Description)
		assert.Equal(t, "http://localhost:3000/user1/-/packages/container/GiteaContainer/latest", pl.Embeds[0].URL)
		assert.Equal(t, p.Sender.UserName, pl.Embeds[0].Author.Name)
		assert.Equal(t, setting.AppURL+p.Sender.UserName, pl.Embeds[0].Author.URL)
		assert.Equal(t, p.Sender.AvatarURL, pl.Embeds[0].Author.IconURL)
	})
}

func TestDiscordJSONPayload(t *testing.T) {
	p := packageTestPayload()
	data, err := p.JSONPayload()
	require.NoError(t, err)

	hook := &webhook_model.Webhook{
		RepoID:     3,
		IsActive:   true,
		Type:       webhook_module.DISCORD,
		URL:        "https://discord.example.com/",
		Meta:       `{}`,
		HTTPMethod: "POST",
	}
	task := &webhook_model.HookTask{
		HookID:         hook.ID,
		EventType:      webhook_module.HookEventPackage,
		PayloadContent: string(data),
		PayloadVersion: 2,
	}

	req, reqBody, err := newDiscordRequest(context.Background(), hook, task)
	require.NotNil(t, req)
	require.NotNil(t, reqBody)
	require.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https://discord.example.com/", req.URL.String())
	assert.Equal(t, "sha256=", req.Header.Get("X-Hub-Signature-256"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	var body DiscordPayload
	err = json.NewDecoder(req.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "[2020558](http://localhost:3000/test/repo/commit/2020558fe2e34debb818a514715839cabd25e778) commit message - user1\n[2020558](http://localhost:3000/test/repo/commit/2020558fe2e34debb818a514715839cabd25e778) commit message - user1", body.Embeds[0].Description)
}
