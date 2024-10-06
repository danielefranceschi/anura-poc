// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cron

import (
	"context"
	"time"

	"code.gitea.io/gitea/models/system"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/setting"
	user_service "code.gitea.io/gitea/services/user"
)

func registerDeleteInactiveUsers() {
	RegisterTaskFatal("delete_inactive_accounts", &OlderThanConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@annually",
		},
		OlderThan: time.Minute * time.Duration(setting.Service.ActiveCodeLives),
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		olderThanConfig := config.(*OlderThanConfig)
		return user_service.DeleteInactiveUsers(ctx, olderThanConfig.OlderThan)
	})
}

func registerDeleteOldSystemNotices() {
	RegisterTaskFatal("delete_old_system_notices", &OlderThanConfig{
		BaseConfig: BaseConfig{
			Enabled:    false,
			RunAtStart: false,
			Schedule:   "@every 168h",
		},
		OlderThan: 365 * 24 * time.Hour,
	}, func(ctx context.Context, _ *user_model.User, config Config) error {
		olderThanConfig := config.(*OlderThanConfig)
		return system.DeleteOldSystemNotices(ctx, olderThanConfig.OlderThan)
	})
}

func initExtendedTasks() {
	registerDeleteInactiveUsers()
	registerDeleteOldSystemNotices()
}
