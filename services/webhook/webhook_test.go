// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	webhook_model "code.gitea.io/gitea/models/webhook"

	"github.com/stretchr/testify/assert"
)

func TestWebhook_GetSlackHook(t *testing.T) {
	w := &webhook_model.Webhook{
		Meta: `{"channel": "foo", "username": "username", "color": "blue"}`,
	}
	slackHook := GetSlackHook(w)
	assert.Equal(t, *slackHook, SlackMeta{
		Channel:  "foo",
		Username: "username",
		Color:    "blue",
	})
}

// func TestPrepareWebhooks(t *testing.T) {
// 	assert.NoError(t, unittest.PrepareTestDatabase())

// 	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
// 	hookTasks := []*webhook_model.HookTask{
// 		{HookID: 1, EventType: webhook_module.HookEventPush},
// 	}
// 	for _, hookTask := range hookTasks {
// 		unittest.AssertNotExistsBean(t, hookTask)
// 	}
// 	assert.NoError(t, PrepareWebhooks(db.DefaultContext, EventSource{Repository: repo}, webhook_module.HookEventPush, &api.PushPayload{Commits: []*api.PayloadCommit{{}}}))
// 	for _, hookTask := range hookTasks {
// 		unittest.AssertExistsAndLoadBean(t, hookTask)
// 	}
// }
