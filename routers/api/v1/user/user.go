// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/routers/api/v1/utils"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/convert"
)

// Search search users
func Search(ctx *context.APIContext) {
	// swagger:operation GET /users/search user userSearch
	// ---
	// summary: Search for users
	// produces:
	// - application/json
	// parameters:
	// - name: q
	//   in: query
	//   description: keyword
	//   type: string
	// - name: uid
	//   in: query
	//   description: ID of the user to search for
	//   type: integer
	//   format: int64
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     description: "SearchResults of a successful search"
	//     schema:
	//       type: object
	//       properties:
	//         ok:
	//           type: boolean
	//         data:
	//           type: array
	//           items:
	//             "$ref": "#/definitions/User"

	listOptions := utils.GetListOptions(ctx)

	uid := ctx.FormInt64("uid")
	var users []*user_model.User
	var maxResults int64
	var err error

	switch uid {
	case user_model.GhostUserID:
		maxResults = 1
		users = []*user_model.User{user_model.NewGhostUser()}
	case user_model.ActionsUserID:
		maxResults = 1
		users = []*user_model.User{user_model.NewActionsUser()}
	default:
		users, maxResults, err = user_model.SearchUsers(ctx, &user_model.SearchUserOptions{
			Actor:       ctx.Doer,
			Keyword:     ctx.FormTrim("q"),
			UID:         uid,
			Type:        user_model.UserTypeIndividual,
			ListOptions: listOptions,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{
				"ok":    false,
				"error": err.Error(),
			})
			return
		}
	}

	ctx.SetLinkHeader(int(maxResults), listOptions.PageSize)
	ctx.SetTotalCountHeader(maxResults)

	ctx.JSON(http.StatusOK, map[string]any{
		"ok":   true,
		"data": convert.ToUsers(ctx, ctx.Doer, users),
	})
}

// GetInfo get user's information
func GetInfo(ctx *context.APIContext) {
	// swagger:operation GET /users/{username} user userGet
	// ---
	// summary: Get a user
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of user to get
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/User"
	//   "404":
	//     "$ref": "#/responses/notFound"

	if !user_model.IsUserVisibleToViewer(ctx, ctx.ContextUser, ctx.Doer) {
		// fake ErrUserNotExist error message to not leak information about existence
		ctx.NotFound("GetUserByName", user_model.ErrUserNotExist{Name: ctx.PathParam(":username")})
		return
	}
	ctx.JSON(http.StatusOK, convert.ToUser(ctx, ctx.ContextUser, ctx.Doer))
}

// GetAuthenticatedUser get current user's information
func GetAuthenticatedUser(ctx *context.APIContext) {
	// swagger:operation GET /user user userGetCurrent
	// ---
	// summary: Get the authenticated user
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/User"

	ctx.JSON(http.StatusOK, convert.ToUser(ctx, ctx.Doer, ctx.Doer))
}
