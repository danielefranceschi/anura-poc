// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "code.gitea.io/gitea/modules/structs"
)

// // Repository
// // swagger:response Repository
// type swaggerResponseRepository struct {
// 	// in:body
// 	Body api.Repository `json:"body"`
// }

// // RepositoryList
// // swagger:response RepositoryList
// type swaggerResponseRepositoryList struct {
// 	// in:body
// 	Body []api.Repository `json:"body"`
// }

// Hook
// swagger:response Hook
type swaggerResponseHook struct {
	// in:body
	Body api.Hook `json:"body"`
}

// HookList
// swagger:response HookList
type swaggerResponseHookList struct {
	// in:body
	Body []api.Hook `json:"body"`
}

// SearchResults
// swagger:response SearchResults
type swaggerResponseSearchResults struct {
	// in:body
	Body api.SearchResults `json:"body"`
}
