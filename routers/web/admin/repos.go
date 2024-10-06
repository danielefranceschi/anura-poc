// Copyright 2014 The Gogs Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/routers/web/explore"
	"code.gitea.io/gitea/services/context"
)

const (
	tplRepos base.TplName = "admin/repo/list"
)

// Repos show all the repositories
func Repos(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.repositories")
	ctx.Data["PageIsAdminRepositories"] = true

	explore.RenderRepoSearch(ctx, &explore.RepoSearchOptions{
		Private:          true,
		PageSize:         setting.UI.Admin.RepoPagingNum,
		TplName:          tplRepos,
		OnlyShowRelevant: false,
	})
}

// // DeleteRepo delete one repository
// func DeleteRepo(ctx *context.Context) {
// 	repo, err := repo_model.GetRepositoryByID(ctx, ctx.FormInt64("id"))
// 	if err != nil {
// 		ctx.ServerError("GetRepositoryByID", err)
// 		return
// 	}

// 	if ctx.Repo != nil && ctx.Repo.GitRepo != nil && ctx.Repo.Repository != nil && ctx.Repo.Repository.ID == repo.ID {
// 		ctx.Repo.GitRepo.Close()
// 	}

// 	if err := repo_service.DeleteRepository(ctx, ctx.Doer, repo, true); err != nil {
// 		ctx.ServerError("DeleteRepository", err)
// 		return
// 	}
// 	log.Trace("Repository deleted: %s", repo.FullName())

// 	ctx.Flash.Success(ctx.Tr("repo.settings.deletion_success"))
// 	ctx.JSONRedirect(setting.AppSubURL + "/admin/repos?page=" + url.QueryEscape(ctx.FormString("page")) + "&sort=" + url.QueryEscape(ctx.FormString("sort")))
// }
