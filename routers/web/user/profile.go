// Copyright 2015 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	"code.gitea.io/gitea/modules/base"
	shared_user "code.gitea.io/gitea/routers/web/shared/user"
	"code.gitea.io/gitea/services/context"
)

const (
	tplProfileBigAvatar base.TplName = "shared/user/profile_big_avatar"
)

// OwnerProfile render profile page for a user or a organization (aka, repo owner)
func OwnerProfile(ctx *context.Context) {
	userProfile(ctx)
}

func userProfile(ctx *context.Context) {
	ctx.Data["Title"] = ctx.ContextUser.DisplayName()
	ctx.Data["PageIsUserProfile"] = true

	// call PrepareContextForProfileBigAvatar later to avoid re-querying the NumFollowers & NumFollowing
	shared_user.PrepareContextForProfileBigAvatar(ctx)
	ctx.HTML(http.StatusOK, tplProfile)
}
