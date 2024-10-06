// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package explore

import (
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/services/context"
)

const (
	// tplExploreRepos explore repositories page template
	tplExploreRepos        base.TplName = "explore/repos"
	relevantReposOnlyParam string       = "only_show_relevant"
)

// RepoSearchOptions when calling search repositories
type RepoSearchOptions struct {
	OwnerID          int64
	Private          bool
	Restricted       bool
	PageSize         int
	OnlyShowRelevant bool
	TplName          base.TplName
}

// RenderRepoSearch render repositories search page
// This function is also used to render the Admin Repository Management page.
func RenderRepoSearch(ctx *context.Context, opts *RepoSearchOptions) {
	// // Sitemap index for sitemap paths
	// page := int(ctx.PathParamInt64("idx"))
	// isSitemap := ctx.PathParam("idx") != ""
	// if page <= 1 {
	// 	page = ctx.FormInt("page")
	// }

	// if page <= 0 {
	// 	page = 1
	// }

	// if isSitemap {
	// 	opts.PageSize = setting.UI.SitemapPagingNum
	// }

	// var (
	// 	repos   []*repo_model.Repository
	// 	count   int64
	// 	err     error
	// 	orderBy db.SearchOrderBy
	// )

	// sortOrder := ctx.FormString("sort")
	// if sortOrder == "" {
	// 	sortOrder = setting.UI.ExploreDefaultSort
	// }

	// if order, ok := repo_model.OrderByFlatMap[sortOrder]; ok {
	// 	orderBy = order
	// } else {
	// 	sortOrder = "recentupdate"
	// 	orderBy = db.SearchOrderByRecentUpdated
	// }
	// ctx.Data["SortType"] = sortOrder

	// keyword := ctx.FormTrim("q")

	// ctx.Data["OnlyShowRelevant"] = opts.OnlyShowRelevant

	// private := ctx.FormOptionalBool("private")
	// ctx.Data["IsPrivate"] = private

	// repos, count, err = repo_model.SearchRepository(ctx, &repo_model.SearchRepoOptions{
	// 	ListOptions: db.ListOptions{
	// 		Page:     page,
	// 		PageSize: opts.PageSize,
	// 	},
	// 	Actor:              ctx.Doer,
	// 	OrderBy:            orderBy,
	// 	Private:            opts.Private,
	// 	Keyword:            keyword,
	// 	OwnerID:            opts.OwnerID,
	// 	AllPublic:          true,
	// 	AllLimited:         true,
	// 	IncludeDescription: setting.UI.SearchRepoDescription,
	// 	OnlyShowRelevant:   opts.OnlyShowRelevant,
	// 	IsPrivate:          private,
	// })
	// if err != nil {
	// 	ctx.ServerError("SearchRepository", err)
	// 	return
	// }
	// if isSitemap {
	// 	m := sitemap.NewSitemap()
	// 	for _, item := range repos {
	// 		m.Add(sitemap.URL{URL: item.HTMLURL(), LastMod: item.UpdatedUnix.AsTimePtr()})
	// 	}
	// 	ctx.Resp.Header().Set("Content-Type", "text/xml")
	// 	if _, err := m.WriteTo(ctx.Resp); err != nil {
	// 		log.Error("Failed writing sitemap: %v", err)
	// 	}
	// 	return
	// }

	// ctx.Data["Keyword"] = keyword
	// ctx.Data["Total"] = count
	// ctx.Data["Repos"] = repos

	// pager := context.NewPagination(int(count), opts.PageSize, page, 5)
	// pager.SetDefaultParams(ctx)
	// if private.Has() {
	// 	pager.AddParamString("private", fmt.Sprint(private.Value()))
	// }
	// ctx.Data["Page"] = pager

	// ctx.HTML(http.StatusOK, opts.TplName)
}

// Repos render explore repositories page
func Repos(ctx *context.Context) {
	ctx.Data["UsersIsDisabled"] = setting.Service.Explore.DisableUsersPage
	ctx.Data["Title"] = ctx.Tr("explore")
	ctx.Data["PageIsExplore"] = true
	ctx.Data["PageIsExploreRepositories"] = true

	var ownerID int64
	if ctx.Doer != nil && !ctx.Doer.IsAdmin {
		ownerID = ctx.Doer.ID
	}

	onlyShowRelevant := setting.UI.OnlyShowRelevantRepos

	_ = ctx.Req.ParseForm() // parse the form first, to prepare the ctx.Req.Form field
	if len(ctx.Req.Form[relevantReposOnlyParam]) != 0 {
		onlyShowRelevant = ctx.FormBool(relevantReposOnlyParam)
	}

	RenderRepoSearch(ctx, &RepoSearchOptions{
		PageSize:         setting.UI.ExplorePagingNum,
		OwnerID:          ownerID,
		Private:          ctx.Doer != nil,
		TplName:          tplExploreRepos,
		OnlyShowRelevant: onlyShowRelevant,
	})
}
