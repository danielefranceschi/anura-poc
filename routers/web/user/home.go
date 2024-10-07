// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"
	"regexp"
	"strings"

	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/services/context"
)

const (
	tplDashboard base.TplName = "user/dashboard/dashboard"
	tplProfile   base.TplName = "user/profile"
)

// getDashboardContextUser finds out which context user dashboard is being viewed as .
func getDashboardContextUser(ctx *context.Context) *user_model.User {
	ctxUser := ctx.Doer
	ctx.Data["ContextUser"] = ctxUser
	return ctxUser
}

// Dashboard render the dashboard page
func Dashboard(ctx *context.Context) {
	ctxUser := getDashboardContextUser(ctx)
	if ctx.Written() {
		return
	}

	var (
		date = ctx.FormString("date")
		page = ctx.FormInt("page")
	)

	// Make sure page number is at least 1. Will be posted to ctx.Data.
	if page <= 1 {
		page = 1
	}

	ctx.Data["Title"] = ctxUser.DisplayName() + " - " + ctx.Locale.TrString("dashboard")
	ctx.Data["PageIsDashboard"] = true
	ctx.Data["PageIsNews"] = true
	ctx.Data["Date"] = date

	var uid int64
	if ctxUser != nil {
		uid = ctxUser.ID
	}

	ctx.PageData["dashboardRepoList"] = map[string]any{
		"searchLimit": setting.UI.User.RepoPagingNum,
		"uid":         uid,
	}

	pager := context.NewPagination(10, setting.UI.FeedPagingNum, page, 5) // TODO
	pager.AddParamString("date", date)
	ctx.Data["Page"] = pager

	ctx.HTML(http.StatusOK, tplDashboard)
}

// Regexp for repos query
var issueReposQueryPattern = regexp.MustCompile(`^\[\d+(,\d+)*,?\]$`)

func UsernameSubRoute(ctx *context.Context) {
	// WORKAROUND to support usernames with "." in it
	// https://github.com/go-chi/chi/issues/781
	username := ctx.PathParam("username")
	reloadParam := func(suffix string) (success bool) {
		ctx.SetPathParam("username", strings.TrimSuffix(username, suffix))
		context.UserAssignmentWeb()(ctx)
		return !ctx.Written()
	}
	switch {
	case strings.HasSuffix(username, ".png"):
		if reloadParam(".png") {
			AvatarByUserName(ctx)
		}
	default:
		context.UserAssignmentWeb()(ctx)
		if !ctx.Written() {
			OwnerProfile(ctx)
		}
	}
}
