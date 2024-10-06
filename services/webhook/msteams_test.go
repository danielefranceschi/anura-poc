// Copyright 2021 The Gitea Authors. All rights reserved.
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

func TestMSTeamsPayload(t *testing.T) {
	mc := msteamsConvertor{}

	t.Run("Package", func(t *testing.T) {
		p := packageTestPayload()

		pl, err := mc.Package(p)
		require.NoError(t, err)

		assert.Equal(t, "Package created: GiteaContainer:latest", pl.Title)
		assert.Equal(t, "Package created: GiteaContainer:latest", pl.Summary)
		assert.Len(t, pl.Sections, 1)
		assert.Equal(t, "user1", pl.Sections[0].ActivitySubtitle)
		assert.Empty(t, pl.Sections[0].Text)
		assert.Len(t, pl.Sections[0].Facts, 1)
		for _, fact := range pl.Sections[0].Facts {
			if fact.Name == "Package:" {
				assert.Equal(t, p.Package.Name, fact.Value)
			} else {
				t.Fail()
			}
		}
		assert.Len(t, pl.PotentialAction, 1)
		assert.Len(t, pl.PotentialAction[0].Targets, 1)
		assert.Equal(t, "http://localhost:3000/user1/-/packages/container/GiteaContainer/latest", pl.PotentialAction[0].Targets[0].URI)
	})

}

func TestMSTeamsJSONPayload(t *testing.T) {
	p := packageTestPayload()
	data, err := p.JSONPayload()
	require.NoError(t, err)

	hook := &webhook_model.Webhook{
		RepoID:     3,
		IsActive:   true,
		Type:       webhook_module.MSTEAMS,
		URL:        "https://msteams.example.com/",
		Meta:       ``,
		HTTPMethod: "POST",
	}
	task := &webhook_model.HookTask{
		HookID:         hook.ID,
		EventType:      webhook_module.HookEventPackage,
		PayloadContent: string(data),
		PayloadVersion: 2,
	}

	req, reqBody, err := newMSTeamsRequest(context.Background(), hook, task)
	require.NotNil(t, req)
	require.NotNil(t, reqBody)
	require.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https://msteams.example.com/", req.URL.String())
	assert.Equal(t, "sha256=", req.Header.Get("X-Hub-Signature-256"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	var body MSTeamsPayload
	err = json.NewDecoder(req.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "[test/repo:test] 2 new commits", body.Summary)
}
