// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"context"
	"testing"

	webhook_model "code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/json"
	webhook_module "code.gitea.io/gitea/modules/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlackPayload(t *testing.T) {
	sc := slackConvertor{}

	t.Run("Package", func(t *testing.T) {
		p := packageTestPayload()

		pl, err := sc.Package(p)
		require.NoError(t, err)

		assert.Equal(t, "Package created: <http://localhost:3000/user1/-/packages/container/GiteaContainer/latest|GiteaContainer:latest> by <https://try.gitea.io/user1|user1>", pl.Text)
	})
}

func TestSlackJSONPayload(t *testing.T) {
	p := packageTestPayload()
	data, err := p.JSONPayload()
	require.NoError(t, err)

	hook := &webhook_model.Webhook{
		RepoID:     3,
		IsActive:   true,
		Type:       webhook_module.SLACK,
		URL:        "https://slack.example.com/",
		Meta:       `{}`,
		HTTPMethod: "POST",
	}
	task := &webhook_model.HookTask{
		HookID:         hook.ID,
		EventType:      webhook_module.HookEventPackage,
		PayloadContent: string(data),
		PayloadVersion: 2,
	}

	req, reqBody, err := newSlackRequest(context.Background(), hook, task)
	require.NotNil(t, req)
	require.NotNil(t, reqBody)
	require.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https://slack.example.com/", req.URL.String())
	assert.Equal(t, "sha256=", req.Header.Get("X-Hub-Signature-256"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	var body SlackPayload
	err = json.NewDecoder(req.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "[<http://localhost:3000/test/repo|test/repo>:<http://localhost:3000/test/repo/src/branch/test|test>] 2 new commits pushed by user1", body.Text)
}

func TestIsValidSlackChannel(t *testing.T) {
	tt := []struct {
		channelName string
		expected    bool
	}{
		{"gitea", true},
		{"#gitea", true},
		{"  ", false},
		{"#", false},
		{" #", false},
		{"gitea   ", false},
		{"  gitea", false},
	}

	for _, v := range tt {
		assert.Equal(t, v.expected, IsValidSlackChannel(v.channelName))
	}
}
