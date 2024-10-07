// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"context"

	"code.gitea.io/gitea/models"
	auth_model "code.gitea.io/gitea/models/auth"
	user_model "code.gitea.io/gitea/models/user"
	password_module "code.gitea.io/gitea/modules/auth/password"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/modules/setting"
)

type UpdateOptions struct {
	KeepEmailPrivate optional.Option[bool]
	FullName         optional.Option[string]
	Website          optional.Option[string]
	Location         optional.Option[string]
	Description      optional.Option[string]
	Language         optional.Option[string]
	Theme            optional.Option[string]
	IsActive         optional.Option[bool]
	IsAdmin          optional.Option[bool]
	SetLastLogin     bool
}

func UpdateUser(ctx context.Context, u *user_model.User, opts *UpdateOptions) error {
	cols := make([]string, 0, 20)

	if opts.KeepEmailPrivate.Has() {
		u.KeepEmailPrivate = opts.KeepEmailPrivate.Value()

		cols = append(cols, "keep_email_private")
	}

	if opts.FullName.Has() {
		u.FullName = opts.FullName.Value()

		cols = append(cols, "full_name")
	}
	if opts.Website.Has() {
		u.Website = opts.Website.Value()

		cols = append(cols, "website")
	}
	if opts.Location.Has() {
		u.Location = opts.Location.Value()

		cols = append(cols, "location")
	}
	if opts.Description.Has() {
		u.Description = opts.Description.Value()

		cols = append(cols, "description")
	}
	if opts.Language.Has() {
		u.Language = opts.Language.Value()

		cols = append(cols, "language")
	}
	if opts.Theme.Has() {
		u.Theme = opts.Theme.Value()

		cols = append(cols, "theme")
	}

	if opts.IsActive.Has() {
		u.IsActive = opts.IsActive.Value()

		cols = append(cols, "is_active")
	}

	if opts.IsAdmin.Has() {
		if !opts.IsAdmin.Value() && user_model.IsLastAdminUser(ctx, u) {
			return models.ErrDeleteLastAdminUser{UID: u.ID}
		}

		u.IsAdmin = opts.IsAdmin.Value()

		cols = append(cols, "is_admin")
	}

	if opts.SetLastLogin {
		u.SetLastLogin()

		cols = append(cols, "last_login_unix")
	}

	return user_model.UpdateUserCols(ctx, u, cols...)
}

type UpdateAuthOptions struct {
	LoginSource        optional.Option[int64]
	LoginName          optional.Option[string]
	Password           optional.Option[string]
	MustChangePassword optional.Option[bool]
	ProhibitLogin      optional.Option[bool]
}

func UpdateAuth(ctx context.Context, u *user_model.User, opts *UpdateAuthOptions) error {
	if opts.LoginSource.Has() {
		source, err := auth_model.GetSourceByID(ctx, opts.LoginSource.Value())
		if err != nil {
			return err
		}

		u.LoginType = source.Type
		u.LoginSource = source.ID
	}
	if opts.LoginName.Has() {
		u.LoginName = opts.LoginName.Value()
	}

	deleteAuthTokens := false
	if opts.Password.Has() && (u.IsLocal() || u.IsOAuth2()) {
		password := opts.Password.Value()

		if len(password) < setting.MinPasswordLength {
			return password_module.ErrMinLength
		}
		if !password_module.IsComplexEnough(password) {
			return password_module.ErrComplexity
		}
		if err := password_module.IsPwned(ctx, password); err != nil {
			return err
		}

		if err := u.SetPassword(password); err != nil {
			return err
		}

		deleteAuthTokens = true
	}

	if opts.MustChangePassword.Has() {
		u.MustChangePassword = opts.MustChangePassword.Value()
	}
	if opts.ProhibitLogin.Has() {
		u.ProhibitLogin = opts.ProhibitLogin.Value()
	}

	if err := user_model.UpdateUserCols(ctx, u, "login_type", "login_source", "login_name", "passwd", "passwd_hash_algo", "salt", "must_change_password", "prohibit_login"); err != nil {
		return err
	}

	if deleteAuthTokens {
		return auth_model.DeleteAuthTokensByUserID(ctx, u.ID)
	}
	return nil
}
