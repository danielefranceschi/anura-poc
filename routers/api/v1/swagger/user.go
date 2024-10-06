// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "code.gitea.io/gitea/modules/structs"
)

// User
// swagger:response User
type swaggerResponseUser struct {
	// in:body
	Body api.User `json:"body"`
}

// UserList
// swagger:response UserList
type swaggerResponseUserList struct {
	// in:body
	Body []api.User `json:"body"`
}

// swagger:model EditUserOption
type swaggerModelEditUserOption struct {
	// in:body
	Options api.EditUserOption
}

// UserSettings
// swagger:response UserSettings
type swaggerResponseUserSettings struct {
	// in:body
	Body []api.UserSettings `json:"body"`
}
