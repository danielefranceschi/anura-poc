// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package activities

import (
	"context"

	"code.gitea.io/gitea/models/db"
	user_model "code.gitea.io/gitea/models/user"
)

// Statistic contains the database statistics
type Statistic struct {
	Counter struct {
		User, Oauth, AuthSource int64
	}
}

// GetStatistic returns the database statistics
func GetStatistic(ctx context.Context) (stats Statistic) {
	_ = db.GetEngine(ctx)
	stats.Counter.User = user_model.CountUsers(ctx, nil)
	stats.Counter.Oauth = 0
	return stats
}
