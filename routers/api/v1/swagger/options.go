// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "code.gitea.io/gitea/modules/structs"
)

// not actually a response, just a hack to get go-swagger to include definitions
// of the various XYZOption structs

// parameterBodies
// swagger:response parameterBodies
type swaggerParameterBodies struct {

	// in:body
	CreateEmailOption api.CreateEmailOption
	// in:body
	DeleteEmailOption api.DeleteEmailOption

	// in:body
	CreateHookOption api.CreateHookOption
	// in:body
	EditHookOption api.EditHookOption

	// in:body
	RenameUserOption api.RenameUserOption

	// in:body
	MarkupOption api.MarkupOption
	// in:body
	MarkdownOption api.MarkdownOption

	// in:body
	CreateUserOption api.CreateUserOption

	// in:body
	EditUserOption api.EditUserOption

	// in:body
	CreateOAuth2ApplicationOptions api.CreateOAuth2ApplicationOptions

	// in:body
	CreateAccessTokenOption api.CreateAccessTokenOption

	// in:body
	UserSettingsOptions api.UserSettingsOptions

	// in:body
	UserBadgeOption api.UserBadgeOption
}
