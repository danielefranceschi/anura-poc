// Copyright 2015 The Gogs Authors. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/util"
	"code.gitea.io/gitea/modules/web"
	webhook_module "code.gitea.io/gitea/modules/webhook"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/forms"
	webhook_service "code.gitea.io/gitea/services/webhook"
)

const (
	tplHooks        base.TplName = "repo/settings/webhook/base"
	tplHookNew      base.TplName = "repo/settings/webhook/new"
	tplOrgHookNew   base.TplName = "org/settings/hook_new"
	tplUserHookNew  base.TplName = "user/settings/hook_new"
	tplAdminHookNew base.TplName = "admin/hook_new"
)

// Webhooks render web hooks list page
func Webhooks(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.hooks")
	ctx.Data["PageIsSettingsHooks"] = true
	// ctx.Data["BaseLink"] = ctx.Repo.RepoLink + "/settings/hooks"
	// ctx.Data["BaseLinkNew"] = ctx.Repo.RepoLink + "/settings/hooks"
	ctx.Data["Description"] = ctx.Tr("repo.settings.hooks_desc", "https://docs.gitea.com/usage/webhooks")

	// ws, err := db.Find[webhook.Webhook](ctx, webhook.ListWebhookOptions{RepoID: ctx.Repo.Repository.ID})
	ws, err := db.Find[webhook.Webhook](ctx, webhook.ListWebhookOptions{})
	if err != nil {
		ctx.ServerError("GetWebhooksByRepoID", err)
		return
	}
	ctx.Data["Webhooks"] = ws

	ctx.HTML(http.StatusOK, tplHooks)
}

type ownerRepoCtx struct {
	OwnerID         int64
	RepoID          int64
	IsAdmin         bool
	IsSystemWebhook bool
	Link            string
	LinkNew         string
	NewTemplate     base.TplName
}

// getOwnerRepoCtx determines whether this is a repo, owner, or admin (both default and system) context.
func getOwnerRepoCtx(ctx *context.Context) (*ownerRepoCtx, error) {
	if ctx.Data["PageIsUserSettings"] == true {
		return &ownerRepoCtx{
			OwnerID:     ctx.Doer.ID,
			Link:        path.Join(setting.AppSubURL, "/user/settings/hooks"),
			LinkNew:     path.Join(setting.AppSubURL, "/user/settings/hooks"),
			NewTemplate: tplUserHookNew,
		}, nil
	}

	if ctx.Data["PageIsAdmin"] == true || ctx.Data["IsAdmin"] == true {
		return &ownerRepoCtx{
			IsAdmin:         true,
			IsSystemWebhook: ctx.PathParam(":configType") == "system-hooks",
			Link:            path.Join(setting.AppSubURL, "/admin/hooks"),
			LinkNew:         path.Join(setting.AppSubURL, "/admin/", ctx.PathParam(":configType")),
			NewTemplate:     tplAdminHookNew,
		}, nil
	}

	return nil, errors.New("unable to set OwnerRepo context: " + fmt.Sprintf("%v", ctx.Data))
}

func checkHookType(ctx *context.Context) string {
	hookType := strings.ToLower(ctx.PathParam(":type"))
	if !util.SliceContainsString(setting.Webhook.Types, hookType, true) {
		ctx.NotFound("checkHookType", nil)
		return ""
	}
	return hookType
}

// WebhooksNew render creating webhook page
func WebhooksNew(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.add_webhook")
	ctx.Data["Webhook"] = webhook.Webhook{HookEvent: &webhook_module.HookEvent{}}

	orCtx, err := getOwnerRepoCtx(ctx)
	if err != nil {
		ctx.ServerError("getOwnerRepoCtx", err)
		return
	}

	if orCtx.IsAdmin && orCtx.IsSystemWebhook {
		ctx.Data["PageIsAdminSystemHooks"] = true
		ctx.Data["PageIsAdminSystemHooksNew"] = true
	} else if orCtx.IsAdmin {
		ctx.Data["PageIsAdminDefaultHooks"] = true
		ctx.Data["PageIsAdminDefaultHooksNew"] = true
	} else {
		ctx.Data["PageIsSettingsHooks"] = true
		ctx.Data["PageIsSettingsHooksNew"] = true
	}

	hookType := checkHookType(ctx)
	ctx.Data["HookType"] = hookType
	if ctx.Written() {
		return
	}
	if hookType == "discord" {
		ctx.Data["DiscordHook"] = map[string]any{
			"Username": "Gitea",
		}
	}
	ctx.Data["BaseLink"] = orCtx.LinkNew
	ctx.Data["BaseLinkNew"] = orCtx.LinkNew

	ctx.HTML(http.StatusOK, orCtx.NewTemplate)
}

// ParseHookEvent convert web form content to webhook.HookEvent
func ParseHookEvent(form forms.WebhookForm) *webhook_module.HookEvent {
	return &webhook_module.HookEvent{
		SendEverything: form.SendEverything(),
		ChooseEvents:   form.ChooseEvents(),
		HookEvents: webhook_module.HookEvents{
			Repository: form.Repository,
			Package:    form.Package,
		},
	}
}

type webhookParams struct {
	// Type should be imported from webhook package (webhook.XXX)
	Type string

	URL         string
	ContentType webhook.HookContentType
	Secret      string
	HTTPMethod  string
	WebhookForm forms.WebhookForm
	Meta        any
}

func createWebhook(ctx *context.Context, params webhookParams) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.add_webhook")
	ctx.Data["PageIsSettingsHooks"] = true
	ctx.Data["PageIsSettingsHooksNew"] = true
	ctx.Data["Webhook"] = webhook.Webhook{HookEvent: &webhook_module.HookEvent{}}
	ctx.Data["HookType"] = params.Type

	orCtx, err := getOwnerRepoCtx(ctx)
	if err != nil {
		ctx.ServerError("getOwnerRepoCtx", err)
		return
	}
	ctx.Data["BaseLink"] = orCtx.LinkNew

	if ctx.HasError() {
		ctx.HTML(http.StatusOK, orCtx.NewTemplate)
		return
	}

	var meta []byte
	if params.Meta != nil {
		meta, err = json.Marshal(params.Meta)
		if err != nil {
			ctx.ServerError("Marshal", err)
			return
		}
	}

	w := &webhook.Webhook{
		RepoID:          orCtx.RepoID,
		URL:             params.URL,
		HTTPMethod:      params.HTTPMethod,
		ContentType:     params.ContentType,
		Secret:          params.Secret,
		HookEvent:       ParseHookEvent(params.WebhookForm),
		IsActive:        params.WebhookForm.Active,
		Type:            params.Type,
		Meta:            string(meta),
		OwnerID:         orCtx.OwnerID,
		IsSystemWebhook: orCtx.IsSystemWebhook,
	}
	err = w.SetHeaderAuthorization(params.WebhookForm.AuthorizationHeader)
	if err != nil {
		ctx.ServerError("SetHeaderAuthorization", err)
		return
	}
	if err := w.UpdateEvent(); err != nil {
		ctx.ServerError("UpdateEvent", err)
		return
	} else if err := webhook.CreateWebhook(ctx, w); err != nil {
		ctx.ServerError("CreateWebhook", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("repo.settings.add_hook_success"))
	ctx.Redirect(orCtx.Link)
}

func editWebhook(ctx *context.Context, params webhookParams) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.update_webhook")
	ctx.Data["PageIsSettingsHooks"] = true
	ctx.Data["PageIsSettingsHooksEdit"] = true

	orCtx, w := checkWebhook(ctx)
	if ctx.Written() {
		return
	}
	ctx.Data["Webhook"] = w

	if ctx.HasError() {
		ctx.HTML(http.StatusOK, orCtx.NewTemplate)
		return
	}

	var meta []byte
	var err error
	if params.Meta != nil {
		meta, err = json.Marshal(params.Meta)
		if err != nil {
			ctx.ServerError("Marshal", err)
			return
		}
	}

	w.URL = params.URL
	w.ContentType = params.ContentType
	w.Secret = params.Secret
	w.HookEvent = ParseHookEvent(params.WebhookForm)
	w.IsActive = params.WebhookForm.Active
	w.HTTPMethod = params.HTTPMethod
	w.Meta = string(meta)

	err = w.SetHeaderAuthorization(params.WebhookForm.AuthorizationHeader)
	if err != nil {
		ctx.ServerError("SetHeaderAuthorization", err)
		return
	}

	if err := w.UpdateEvent(); err != nil {
		ctx.ServerError("UpdateEvent", err)
		return
	} else if err := webhook.UpdateWebhook(ctx, w); err != nil {
		ctx.ServerError("UpdateWebhook", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("repo.settings.update_hook_success"))
	ctx.Redirect(fmt.Sprintf("%s/%d", orCtx.Link, w.ID))
}

// DiscordHooksNewPost response for creating Discord webhook
func DiscordHooksNewPost(ctx *context.Context) {
	createWebhook(ctx, discordHookParams(ctx))
}

// DiscordHooksEditPost response for editing Discord webhook
func DiscordHooksEditPost(ctx *context.Context) {
	editWebhook(ctx, discordHookParams(ctx))
}

func discordHookParams(ctx *context.Context) webhookParams {
	form := web.GetForm(ctx).(*forms.NewDiscordHookForm)

	return webhookParams{
		Type:        webhook_module.DISCORD,
		URL:         form.PayloadURL,
		ContentType: webhook.ContentTypeJSON,
		WebhookForm: form.WebhookForm,
		Meta: &webhook_service.DiscordMeta{
			Username: form.Username,
			IconURL:  form.IconURL,
		},
	}
}

// MSTeamsHooksNewPost response for creating MSTeams webhook
func MSTeamsHooksNewPost(ctx *context.Context) {
	createWebhook(ctx, mSTeamsHookParams(ctx))
}

// MSTeamsHooksEditPost response for editing MSTeams webhook
func MSTeamsHooksEditPost(ctx *context.Context) {
	editWebhook(ctx, mSTeamsHookParams(ctx))
}

func mSTeamsHookParams(ctx *context.Context) webhookParams {
	form := web.GetForm(ctx).(*forms.NewMSTeamsHookForm)

	return webhookParams{
		Type:        webhook_module.MSTEAMS,
		URL:         form.PayloadURL,
		ContentType: webhook.ContentTypeJSON,
		WebhookForm: form.WebhookForm,
	}
}

// SlackHooksNewPost response for creating Slack webhook
func SlackHooksNewPost(ctx *context.Context) {
	createWebhook(ctx, slackHookParams(ctx))
}

// SlackHooksEditPost response for editing Slack webhook
func SlackHooksEditPost(ctx *context.Context) {
	editWebhook(ctx, slackHookParams(ctx))
}

func slackHookParams(ctx *context.Context) webhookParams {
	form := web.GetForm(ctx).(*forms.NewSlackHookForm)

	return webhookParams{
		Type:        webhook_module.SLACK,
		URL:         form.PayloadURL,
		ContentType: webhook.ContentTypeJSON,
		WebhookForm: form.WebhookForm,
		Meta: &webhook_service.SlackMeta{
			Channel:  strings.TrimSpace(form.Channel),
			Username: form.Username,
			IconURL:  form.IconURL,
			Color:    form.Color,
		},
	}
}

func checkWebhook(ctx *context.Context) (*ownerRepoCtx, *webhook.Webhook) {
	orCtx, err := getOwnerRepoCtx(ctx)
	if err != nil {
		ctx.ServerError("getOwnerRepoCtx", err)
		return nil, nil
	}
	ctx.Data["BaseLink"] = orCtx.Link
	ctx.Data["BaseLinkNew"] = orCtx.LinkNew

	var w *webhook.Webhook
	if orCtx.RepoID > 0 {
		w, err = webhook.GetWebhookByRepoID(ctx, orCtx.RepoID, ctx.PathParamInt64(":id"))
	} else if orCtx.OwnerID > 0 {
		w, err = webhook.GetWebhookByOwnerID(ctx, orCtx.OwnerID, ctx.PathParamInt64(":id"))
	} else if orCtx.IsAdmin {
		w, err = webhook.GetSystemOrDefaultWebhook(ctx, ctx.PathParamInt64(":id"))
	}
	if err != nil || w == nil {
		if webhook.IsErrWebhookNotExist(err) {
			ctx.NotFound("GetWebhookByID", nil)
		} else {
			ctx.ServerError("GetWebhookByID", err)
		}
		return nil, nil
	}

	ctx.Data["HookType"] = w.Type
	switch w.Type {
	case webhook_module.SLACK:
		ctx.Data["SlackHook"] = webhook_service.GetSlackHook(w)
	case webhook_module.DISCORD:
		ctx.Data["DiscordHook"] = webhook_service.GetDiscordHook(w)
	}

	ctx.Data["History"], err = w.History(ctx, 1)
	if err != nil {
		ctx.ServerError("History", err)
	}
	return orCtx, w
}

// WebHooksEdit render editing web hook page
func WebHooksEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.settings.update_webhook")
	ctx.Data["PageIsSettingsHooks"] = true
	ctx.Data["PageIsSettingsHooksEdit"] = true

	orCtx, w := checkWebhook(ctx)
	if ctx.Written() {
		return
	}
	ctx.Data["Webhook"] = w

	ctx.HTML(http.StatusOK, orCtx.NewTemplate)
}

// // TestWebhook test if web hook is work fine
// func TestWebhook(ctx *context.Context) {
// 	hookID := ctx.PathParamInt64(":id")
// 	w, err := webhook.GetWebhookByRepoID(ctx, ctx.Repo.Repository.ID, hookID)
// 	if err != nil {
// 		ctx.Flash.Error("GetWebhookByRepoID: " + err.Error())
// 		ctx.Status(http.StatusInternalServerError)
// 		return
// 	}

// 	// Grab latest commit or fake one if it's empty repository.
// 	commit := ctx.Repo.Commit
// 	if commit == nil {
// 		ghost := user_model.NewGhostUser()
// 		objectFormat := git.ObjectFormatFromName(ctx.Repo.Repository.ObjectFormatName)
// 		commit = &git.Commit{
// 			ID:            objectFormat.EmptyObjectID(),
// 			Author:        ghost.NewGitSig(),
// 			Committer:     ghost.NewGitSig(),
// 			CommitMessage: "This is a fake commit",
// 		}
// 	}

// 	apiUser := convert.ToUserWithAccessMode(ctx, ctx.Doer, perm.AccessModeNone)

// 	apiCommit := &api.PayloadCommit{
// 		ID:      commit.ID.String(),
// 		Message: commit.Message(),
// 		URL:     ctx.Repo.Repository.HTMLURL() + "/commit/" + url.PathEscape(commit.ID.String()),
// 		Author: &api.PayloadUser{
// 			Name:  commit.Author.Name,
// 			Email: commit.Author.Email,
// 		},
// 		Committer: &api.PayloadUser{
// 			Name:  commit.Committer.Name,
// 			Email: commit.Committer.Email,
// 		},
// 	}

// 	commitID := commit.ID.String()
// 	p := &api.PushPayload{
// 		Ref:          git.BranchPrefix + ctx.Repo.Repository.DefaultBranch,
// 		Before:       commitID,
// 		After:        commitID,
// 		CompareURL:   setting.AppURL + ctx.Repo.Repository.ComposeCompareURL(commitID, commitID),
// 		Commits:      []*api.PayloadCommit{apiCommit},
// 		TotalCommits: 1,
// 		HeadCommit:   apiCommit,
// 		Repo:         convert.ToRepo(ctx, ctx.Repo.Repository, access_model.Permission{AccessMode: perm.AccessModeNone}),
// 		Pusher:       apiUser,
// 		Sender:       apiUser,
// 	}
// 	if err := webhook_service.PrepareWebhook(ctx, w, webhook_module.HookEventPush, p); err != nil {
// 		ctx.Flash.Error("PrepareWebhook: " + err.Error())
// 		ctx.Status(http.StatusInternalServerError)
// 	} else {
// 		ctx.Flash.Info(ctx.Tr("repo.settings.webhook.delivery.success"))
// 		ctx.Status(http.StatusOK)
// 	}
// }

// ReplayWebhook replays a webhook
func ReplayWebhook(ctx *context.Context) {
	hookTaskUUID := ctx.PathParam(":uuid")

	orCtx, w := checkWebhook(ctx)
	if ctx.Written() {
		return
	}

	if err := webhook_service.ReplayHookTask(ctx, w, hookTaskUUID); err != nil {
		if webhook.IsErrHookTaskNotExist(err) {
			ctx.NotFound("ReplayHookTask", nil)
		} else {
			ctx.ServerError("ReplayHookTask", err)
		}
		return
	}

	ctx.Flash.Success(ctx.Tr("repo.settings.webhook.delivery.success"))
	ctx.Redirect(fmt.Sprintf("%s/%d", orCtx.Link, w.ID))
}

// DeleteWebhook delete a webhook
func DeleteWebhook(ctx *context.Context) {
	// if err := webhook.DeleteWebhookByRepoID(ctx, ctx.Repo.Repository.ID, ctx.FormInt64("id")); err != nil {
	// 	ctx.Flash.Error("DeleteWebhookByRepoID: " + err.Error())
	// } else {
	// 	ctx.Flash.Success(ctx.Tr("repo.settings.webhook_deletion_success"))
	// }

	// ctx.JSONRedirect(ctx.Repo.RepoLink + "/settings/hooks")
}
