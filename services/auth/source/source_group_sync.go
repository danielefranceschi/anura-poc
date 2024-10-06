// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package source

import (
	"context"

	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/container"
)

// SyncGroupsToTeams maps authentication source groups to organization and team memberships
func SyncGroupsToTeams(ctx context.Context, user *user_model.User, sourceUserGroups container.Set[string], sourceGroupTeamMapping map[string]map[string][]string, performRemoval bool) error {
	return SyncGroupsToTeamsCached(ctx, user, sourceUserGroups, sourceGroupTeamMapping, performRemoval)
}

// SyncGroupsToTeamsCached maps authentication source groups to organization and team memberships
func SyncGroupsToTeamsCached(ctx context.Context, user *user_model.User, sourceUserGroups container.Set[string], sourceGroupTeamMapping map[string]map[string][]string, performRemoval bool) error {
	return nil
}
