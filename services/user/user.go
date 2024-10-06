// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/models/db"
	system_model "code.gitea.io/gitea/models/system"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/eventsource"
	"code.gitea.io/gitea/modules/storage"
)

// RenameUser renames a user
func RenameUser(ctx context.Context, u *user_model.User, newUserName string) error {
	// Non-local users are not allowed to change their username.
	if !u.IsLocal() {
		return user_model.ErrUserIsNotLocal{
			UID:  u.ID,
			Name: u.Name,
		}
	}

	if newUserName == u.Name {
		return nil
	}

	if err := user_model.IsUsableUsername(newUserName); err != nil {
		return err
	}

	onlyCapitalization := strings.EqualFold(newUserName, u.Name)
	oldUserName := u.Name

	if onlyCapitalization {
		u.Name = newUserName
		if err := user_model.UpdateUserCols(ctx, u, "name"); err != nil {
			u.Name = oldUserName
			return err
		}
		return nil
	}
	// TODO: update last downloader of packages

	u.Name = newUserName
	u.LowerName = strings.ToLower(newUserName)
	if err := user_model.UpdateUserCols(ctx, u, "name", "lower_name"); err != nil {
		u.Name = oldUserName
		u.LowerName = strings.ToLower(oldUserName)
		return err
	}

	return nil
}

// DeleteUser completely and permanently deletes everything of a user,
// but issues/comments/pulls will be kept and shown as someone has been deleted,
// unless the user is younger than USER_DELETE_WITH_COMMENTS_MAX_DAYS.
func DeleteUser(ctx context.Context, u *user_model.User, purge bool) error {
	if u.IsActive && user_model.IsLastAdminUser(ctx, u) {
		return models.ErrDeleteLastAdminUser{UID: u.ID}
	}

	if purge {
		// Disable the user first
		// NOTE: This is deliberately not within a transaction as it must disable the user immediately to prevent any further action by the user to be purged.
		if err := user_model.UpdateUserCols(ctx, &user_model.User{
			ID:              u.ID,
			IsActive:        false,
			IsRestricted:    true,
			IsAdmin:         false,
			ProhibitLogin:   true,
			Passwd:          "",
			Salt:            "",
			PasswdHashAlgo:  "",
			MaxRepoCreation: 0,
		}, "is_active", "is_restricted", "is_admin", "prohibit_login", "max_repo_creation", "passwd", "salt", "passwd_hash_algo"); err != nil {
			return fmt.Errorf("unable to disable user: %s[%d] prior to purge. UpdateUserCols: %w", u.Name, u.ID, err)
		}

		// Force any logged in sessions to log out
		// FIXME: We also need to tell the session manager to log them out too.
		eventsource.GetManager().SendMessage(u.ID, &eventsource.Event{
			Name: "logout",
		})
	}

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return err
	}
	defer committer.Close()

	if err := deleteUser(ctx, u, purge); err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	if err := committer.Commit(); err != nil {
		return err
	}
	_ = committer.Close()

	if u.Avatar != "" {
		avatarPath := u.CustomAvatarRelativePath()
		if err = storage.Avatars.Delete(avatarPath); err != nil {
			err = fmt.Errorf("failed to remove %s: %w", avatarPath, err)
			_ = system_model.CreateNotice(ctx, system_model.NoticeTask, fmt.Sprintf("delete user '%s': %v", u.Name, err))
		}
	}

	return nil
}

// DeleteInactiveUsers deletes all inactive users and their email addresses.
func DeleteInactiveUsers(ctx context.Context, olderThan time.Duration) error {
	inactiveUsers, err := user_model.GetInactiveUsers(ctx, olderThan)
	if err != nil {
		return err
	}

	// FIXME: should only update authorized_keys file once after all deletions.
	for _, u := range inactiveUsers {
		if err = DeleteUser(ctx, u, false); err != nil {
			// Ignore inactive users that were ever active but then were set inactive by admin
			if models.IsErrUserOwnPackages(err) {
				continue
			}
			select {
			case <-ctx.Done():
				return db.ErrCancelledf("when deleting inactive user %q", u.Name)
			default:
				return err
			}
		}
	}
	return nil // TODO: there could be still inactive users left, and the number would increase gradually
}
