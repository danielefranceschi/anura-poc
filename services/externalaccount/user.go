// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package externalaccount

import (
	"context"

	"code.gitea.io/gitea/models/auth"
	user_model "code.gitea.io/gitea/models/user"

	"github.com/markbates/goth"
)

func toExternalLoginUser(ctx context.Context, user *user_model.User, gothUser goth.User) (*user_model.ExternalLoginUser, error) {
	authSource, err := auth.GetActiveOAuth2SourceByName(ctx, gothUser.Provider)
	if err != nil {
		return nil, err
	}
	return &user_model.ExternalLoginUser{
		ExternalID:        gothUser.UserID,
		UserID:            user.ID,
		LoginSourceID:     authSource.ID,
		RawData:           gothUser.RawData,
		Provider:          gothUser.Provider,
		Email:             gothUser.Email,
		Name:              gothUser.Name,
		FirstName:         gothUser.FirstName,
		LastName:          gothUser.LastName,
		NickName:          gothUser.NickName,
		Description:       gothUser.Description,
		AvatarURL:         gothUser.AvatarURL,
		Location:          gothUser.Location,
		AccessToken:       gothUser.AccessToken,
		AccessTokenSecret: gothUser.AccessTokenSecret,
		RefreshToken:      gothUser.RefreshToken,
		ExpiresAt:         gothUser.ExpiresAt,
	}, nil
}

// LinkAccountToUser link the gothUser to the user
func LinkAccountToUser(ctx context.Context, user *user_model.User, gothUser goth.User) error {
	externalLoginUser, err := toExternalLoginUser(ctx, user, gothUser)
	if err != nil {
		return err
	}

	if err := user_model.LinkExternalToUser(ctx, user, externalLoginUser); err != nil {
		return err
	}

	return nil
}

// EnsureLinkExternalToUser link the gothUser to the user
func EnsureLinkExternalToUser(ctx context.Context, user *user_model.User, gothUser goth.User) error {
	externalLoginUser, err := toExternalLoginUser(ctx, user, gothUser)
	if err != nil {
		return err
	}

	return user_model.EnsureLinkExternalToUser(ctx, externalLoginUser)
}
