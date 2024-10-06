// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	api "code.gitea.io/gitea/modules/structs"
)

func packageTestPayload() *api.PackagePayload {
	return &api.PackagePayload{
		Action: api.HookPackageCreated,
		Sender: &api.User{
			UserName:  "user1",
			AvatarURL: "http://localhost:3000/user1/avatar",
		},
		// Repository: nil,
		Package: &api.Package{
			Owner: &api.User{
				UserName:  "user1",
				AvatarURL: "http://localhost:3000/user1/avatar",
			},
			// Repository: nil,
			Creator: &api.User{
				UserName:  "user1",
				AvatarURL: "http://localhost:3000/user1/avatar",
			},
			Type:    "container",
			Name:    "GiteaContainer",
			Version: "latest",
			HTMLURL: "http://localhost:3000/user1/-/packages/container/GiteaContainer/latest",
		},
	}
}
