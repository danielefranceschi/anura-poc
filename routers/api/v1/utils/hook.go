// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/util"
	webhook_module "code.gitea.io/gitea/modules/webhook"
	"code.gitea.io/gitea/services/context"
	webhook_service "code.gitea.io/gitea/services/webhook"
)

// checkCreateHookOption check if a CreateHookOption form is valid. If invalid,
// write the appropriate error to `ctx`. Return whether the form is valid
func checkCreateHookOption(ctx *context.APIContext, form *api.CreateHookOption) bool {
	if !webhook_service.IsValidHookTaskType(form.Type) {
		ctx.Error(http.StatusUnprocessableEntity, "", fmt.Sprintf("Invalid hook type: %s", form.Type))
		return false
	}
	for _, name := range []string{"url", "content_type"} {
		if _, ok := form.Config[name]; !ok {
			ctx.Error(http.StatusUnprocessableEntity, "", "Missing config option: "+name)
			return false
		}
	}
	if !webhook.IsValidHookContentType(form.Config["content_type"]) {
		ctx.Error(http.StatusUnprocessableEntity, "", "Invalid content type")
		return false
	}
	return true
}

// AddSystemHook add a system hook
func AddSystemHook(ctx *context.APIContext, form *api.CreateHookOption) {
	hook, ok := addHook(ctx, form, 0, 0)
	if ok {
		h, err := webhook_service.ToHook(setting.AppSubURL+"/admin", hook)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, "convert.ToHook", err)
			return
		}
		ctx.JSON(http.StatusCreated, h)
	}
}

// toAPIHook converts the hook to its API representation.
// If there is an error, write to `ctx` accordingly. Return (hook, ok)
func toAPIHook(ctx *context.APIContext, repoLink string, hook *webhook.Webhook) (*api.Hook, bool) {
	apiHook, err := webhook_service.ToHook(repoLink, hook)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "ToHook", err)
		return nil, false
	}
	return apiHook, true
}

// addHook add the hook specified by `form`, `ownerID` and `repoID`. If there is
// an error, write to `ctx` accordingly. Return (webhook, ok)
func addHook(ctx *context.APIContext, form *api.CreateHookOption, ownerID, repoID int64) (*webhook.Webhook, bool) {
	var isSystemWebhook bool
	if !checkCreateHookOption(ctx, form) {
		return nil, false
	}

	if len(form.Events) == 0 {
		form.Events = []string{"push"}
	}
	if form.Config["is_system_webhook"] != "" {
		sw, err := strconv.ParseBool(form.Config["is_system_webhook"])
		if err != nil {
			ctx.Error(http.StatusUnprocessableEntity, "", "Invalid is_system_webhook value")
			return nil, false
		}
		isSystemWebhook = sw
	}
	w := &webhook.Webhook{
		OwnerID:         ownerID,
		RepoID:          repoID,
		URL:             form.Config["url"],
		ContentType:     webhook.ToHookContentType(form.Config["content_type"]),
		Secret:          form.Config["secret"],
		HTTPMethod:      "POST",
		IsSystemWebhook: isSystemWebhook,
		HookEvent: &webhook_module.HookEvent{
			ChooseEvents: true,
			HookEvents: webhook_module.HookEvents{
				Create: util.SliceContainsString(form.Events, string(webhook_module.HookEventCreate), true),
				Delete: util.SliceContainsString(form.Events, string(webhook_module.HookEventDelete), true),
			},
		},
		IsActive: form.Active,
		Type:     form.Type,
	}
	err := w.SetHeaderAuthorization(form.AuthorizationHeader)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "SetHeaderAuthorization", err)
		return nil, false
	}
	if w.Type == webhook_module.SLACK {
		channel, ok := form.Config["channel"]
		if !ok {
			ctx.Error(http.StatusUnprocessableEntity, "", "Missing config option: channel")
			return nil, false
		}
		channel = strings.TrimSpace(channel)

		if !webhook_service.IsValidSlackChannel(channel) {
			ctx.Error(http.StatusBadRequest, "", "Invalid slack channel name")
			return nil, false
		}

		meta, err := json.Marshal(&webhook_service.SlackMeta{
			Channel:  channel,
			Username: form.Config["username"],
			IconURL:  form.Config["icon_url"],
			Color:    form.Config["color"],
		})
		if err != nil {
			ctx.Error(http.StatusInternalServerError, "slack: JSON marshal failed", err)
			return nil, false
		}
		w.Meta = string(meta)
	}

	if err := w.UpdateEvent(); err != nil {
		ctx.Error(http.StatusInternalServerError, "UpdateEvent", err)
		return nil, false
	} else if err := webhook.CreateWebhook(ctx, w); err != nil {
		ctx.Error(http.StatusInternalServerError, "CreateWebhook", err)
		return nil, false
	}
	return w, true
}

// EditSystemHook edit system webhook `w` according to `form`. Writes to `ctx` accordingly
func EditSystemHook(ctx *context.APIContext, form *api.EditHookOption, hookID int64) {
	hook, err := webhook.GetSystemOrDefaultWebhook(ctx, hookID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetSystemOrDefaultWebhook", err)
		return
	}
	if !editHook(ctx, form, hook) {
		ctx.Error(http.StatusInternalServerError, "editHook", err)
		return
	}
	updated, err := webhook.GetSystemOrDefaultWebhook(ctx, hookID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetSystemOrDefaultWebhook", err)
		return
	}
	h, err := webhook_service.ToHook(setting.AppURL+"/admin", updated)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "convert.ToHook", err)
		return
	}
	ctx.JSON(http.StatusOK, h)
}

// editHook edit the webhook `w` according to `form`. If an error occurs, write
// to `ctx` accordingly and return the error. Return whether successful
func editHook(ctx *context.APIContext, form *api.EditHookOption, w *webhook.Webhook) bool {
	if form.Config != nil {
		if url, ok := form.Config["url"]; ok {
			w.URL = url
		}
		if ct, ok := form.Config["content_type"]; ok {
			if !webhook.IsValidHookContentType(ct) {
				ctx.Error(http.StatusUnprocessableEntity, "", "Invalid content type")
				return false
			}
			w.ContentType = webhook.ToHookContentType(ct)
		}

		if w.Type == webhook_module.SLACK {
			if channel, ok := form.Config["channel"]; ok {
				meta, err := json.Marshal(&webhook_service.SlackMeta{
					Channel:  channel,
					Username: form.Config["username"],
					IconURL:  form.Config["icon_url"],
					Color:    form.Config["color"],
				})
				if err != nil {
					ctx.Error(http.StatusInternalServerError, "slack: JSON marshal failed", err)
					return false
				}
				w.Meta = string(meta)
			}
		}
	}

	// Update events
	if len(form.Events) == 0 {
		form.Events = []string{"push"}
	}
	w.SendEverything = false
	w.ChooseEvents = true
	w.Create = util.SliceContainsString(form.Events, string(webhook_module.HookEventCreate), true)
	w.Delete = util.SliceContainsString(form.Events, string(webhook_module.HookEventDelete), true)

	err := w.SetHeaderAuthorization(form.AuthorizationHeader)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "SetHeaderAuthorization", err)
		return false
	}

	if err := w.UpdateEvent(); err != nil {
		ctx.Error(http.StatusInternalServerError, "UpdateEvent", err)
		return false
	}

	if form.Active != nil {
		w.IsActive = *form.Active
	}

	if err := webhook.UpdateWebhook(ctx, w); err != nil {
		ctx.Error(http.StatusInternalServerError, "UpdateWebhook", err)
		return false
	}
	return true
}
