// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"context"
	"fmt"

	_ "image/jpeg" // Needed for jpeg support

	auth_model "code.gitea.io/gitea/models/auth"
	"code.gitea.io/gitea/models/db"
	user_model "code.gitea.io/gitea/models/user"
)

// deleteUser deletes models associated to an user.
func deleteUser(ctx context.Context, u *user_model.User, purge bool) (err error) {
	// e := db.GetEngine(ctx)

	// ***** START: Star *****
	// starredRepoIDs, err := db.FindIDs(ctx, "star", "star.repo_id",
	// 	builder.Eq{"star.uid": u.ID})
	// if err != nil {
	// 	return fmt.Errorf("get all stars: %w", err)
	// } else if err = db.DecrByIDs(ctx, starredRepoIDs, "num_stars", new(repo_model.Repository)); err != nil {
	// 	return fmt.Errorf("decrease repository num_stars: %w", err)
	// }
	// ***** END: Star *****

	if err = db.DeleteBeans(ctx,
		&auth_model.AccessToken{UID: u.ID},
		&user_model.EmailAddress{UID: u.ID},
		&user_model.UserOpenID{UID: u.ID},
		&user_model.Setting{UserID: u.ID},
		&user_model.UserBadge{UserID: u.ID},
	); err != nil {
		return fmt.Errorf("deleteBeans: %w", err)
	}

	if err := auth_model.DeleteOAuth2RelictsByUserID(ctx, u.ID); err != nil {
		return err
	}

	// ***** START: ExternalLoginUser *****
	if err = user_model.RemoveAllAccountLinks(ctx, u); err != nil {
		return fmt.Errorf("ExternalLoginUser: %w", err)
	}
	// ***** END: ExternalLoginUser *****

	if err := auth_model.DeleteAuthTokensByUserID(ctx, u.ID); err != nil {
		return fmt.Errorf("DeleteAuthTokensByUserID: %w", err)
	}

	if _, err = db.DeleteByID[user_model.User](ctx, u.ID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
