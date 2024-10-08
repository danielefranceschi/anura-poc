// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import api "code.gitea.io/gitea/modules/structs"

// GeneralUISettings
// swagger:response GeneralUISettings
type swaggerResponseGeneralUISettings struct {
	// in:body
	Body api.GeneralUISettings `json:"body"`
}

// GeneralAPISettings
// swagger:response GeneralAPISettings
type swaggerResponseGeneralAPISettings struct {
	// in:body
	Body api.GeneralAPISettings `json:"body"`
}
