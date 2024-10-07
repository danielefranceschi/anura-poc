// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"code.gitea.io/gitea/models/avatars"
	"code.gitea.io/gitea/models/db"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/translation"
	"code.gitea.io/gitea/modules/typesniffer"
	"code.gitea.io/gitea/modules/util"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/modules/web/middleware"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/forms"
	user_service "code.gitea.io/gitea/services/user"
	"code.gitea.io/gitea/services/webtheme"
)

const (
	tplSettingsProfile      base.TplName = "user/settings/profile"
	tplSettingsAppearance   base.TplName = "user/settings/appearance"
	tplSettingsOrganization base.TplName = "user/settings/organization"
	tplSettingsRepositories base.TplName = "user/settings/repos"
)

// Profile render user's profile page
func Profile(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.profile")
	ctx.Data["PageIsSettingsProfile"] = true
	ctx.Data["DisableGravatar"] = setting.Config().Picture.DisableGravatar.Value(ctx)

	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsProfile)
}

// ProfilePost response for change user's profile
func ProfilePost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsProfile"] = true
	ctx.Data["DisableGravatar"] = setting.Config().Picture.DisableGravatar.Value(ctx)
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	if ctx.HasError() {
		ctx.HTML(http.StatusOK, tplSettingsProfile)
		return
	}

	form := web.GetForm(ctx).(*forms.UpdateProfileForm)

	if form.Name != "" {
		if err := user_service.RenameUser(ctx, ctx.Doer, form.Name); err != nil {
			switch {
			case user_model.IsErrUserIsNotLocal(err):
				ctx.Flash.Error(ctx.Tr("form.username_change_not_local_user"))
			case user_model.IsErrUserAlreadyExist(err):
				ctx.Flash.Error(ctx.Tr("form.username_been_taken"))
			case db.IsErrNameReserved(err):
				ctx.Flash.Error(ctx.Tr("user.form.name_reserved", form.Name))
			case db.IsErrNamePatternNotAllowed(err):
				ctx.Flash.Error(ctx.Tr("user.form.name_pattern_not_allowed", form.Name))
			case db.IsErrNameCharsNotAllowed(err):
				ctx.Flash.Error(ctx.Tr("user.form.name_chars_not_allowed", form.Name))
			default:
				ctx.ServerError("RenameUser", err)
				return
			}
			ctx.Redirect(setting.AppSubURL + "/user/settings")
			return
		}
	}

	opts := &user_service.UpdateOptions{
		FullName:         optional.Some(form.FullName),
		KeepEmailPrivate: optional.Some(form.KeepEmailPrivate),
		Description:      optional.Some(form.Description),
		Website:          optional.Some(form.Website),
		Location:         optional.Some(form.Location),
	}
	if err := user_service.UpdateUser(ctx, ctx.Doer, opts); err != nil {
		ctx.ServerError("UpdateUser", err)
		return
	}

	log.Trace("User settings updated: %s", ctx.Doer.Name)
	ctx.Flash.Success(ctx.Tr("settings.update_profile_success"))
	ctx.Redirect(setting.AppSubURL + "/user/settings")
}

// UpdateAvatarSetting update user's avatar
// FIXME: limit size.
func UpdateAvatarSetting(ctx *context.Context, form *forms.AvatarForm, ctxUser *user_model.User) error {
	ctxUser.UseCustomAvatar = form.Source == forms.AvatarLocal
	if len(form.Gravatar) > 0 {
		if form.Avatar != nil {
			ctxUser.Avatar = avatars.HashEmail(form.Gravatar)
		} else {
			ctxUser.Avatar = ""
		}
		ctxUser.AvatarEmail = form.Gravatar
	}

	if form.Avatar != nil && form.Avatar.Filename != "" {
		fr, err := form.Avatar.Open()
		if err != nil {
			return fmt.Errorf("Avatar.Open: %w", err)
		}
		defer fr.Close()

		if form.Avatar.Size > setting.Avatar.MaxFileSize {
			return errors.New(ctx.Locale.TrString("settings.uploaded_avatar_is_too_big", form.Avatar.Size/1024, setting.Avatar.MaxFileSize/1024))
		}

		data, err := io.ReadAll(fr)
		if err != nil {
			return fmt.Errorf("io.ReadAll: %w", err)
		}

		st := typesniffer.DetectContentType(data)
		if !(st.IsImage() && !st.IsSvgImage()) {
			return errors.New(ctx.Locale.TrString("settings.uploaded_avatar_not_a_image"))
		}
		if err = user_service.UploadAvatar(ctx, ctxUser, data); err != nil {
			return fmt.Errorf("UploadAvatar: %w", err)
		}
	} else if ctxUser.UseCustomAvatar && ctxUser.Avatar == "" {
		// No avatar is uploaded but setting has been changed to enable,
		// generate a random one when needed.
		if err := user_model.GenerateRandomAvatar(ctx, ctxUser); err != nil {
			log.Error("GenerateRandomAvatar[%d]: %v", ctxUser.ID, err)
		}
	}

	if err := user_model.UpdateUserCols(ctx, ctxUser, "avatar", "avatar_email", "use_custom_avatar"); err != nil {
		return fmt.Errorf("UpdateUserCols: %w", err)
	}

	return nil
}

// AvatarPost response for change user's avatar request
func AvatarPost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.AvatarForm)
	if err := UpdateAvatarSetting(ctx, form, ctx.Doer); err != nil {
		ctx.Flash.Error(err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("settings.update_avatar_success"))
	}

	ctx.Redirect(setting.AppSubURL + "/user/settings")
}

// DeleteAvatar render delete avatar page
func DeleteAvatar(ctx *context.Context) {
	if err := user_service.DeleteAvatar(ctx, ctx.Doer); err != nil {
		ctx.Flash.Error(err.Error())
	}

	ctx.JSONRedirect(setting.AppSubURL + "/user/settings")
}

// Appearance render user's appearance settings
func Appearance(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.appearance")
	ctx.Data["PageIsSettingsAppearance"] = true

	allThemes := webtheme.GetAvailableThemes()
	if webtheme.IsThemeAvailable(setting.UI.DefaultTheme) {
		allThemes = util.SliceRemoveAll(allThemes, setting.UI.DefaultTheme)
		allThemes = append([]string{setting.UI.DefaultTheme}, allThemes...) // move the default theme to the top
	}
	ctx.Data["AllThemes"] = allThemes
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsAppearance)
}

// UpdateUIThemePost is used to update users' specific theme
func UpdateUIThemePost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.UpdateThemeForm)
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsAppearance"] = true

	if ctx.HasError() {
		ctx.Flash.Error(ctx.GetErrMsg())
		ctx.Redirect(setting.AppSubURL + "/user/settings/appearance")
		return
	}

	if !webtheme.IsThemeAvailable(form.Theme) {
		ctx.Flash.Error(ctx.Tr("settings.theme_update_error"))
		ctx.Redirect(setting.AppSubURL + "/user/settings/appearance")
		return
	}

	opts := &user_service.UpdateOptions{
		Theme: optional.Some(form.Theme),
	}
	if err := user_service.UpdateUser(ctx, ctx.Doer, opts); err != nil {
		ctx.Flash.Error(ctx.Tr("settings.theme_update_error"))
	} else {
		ctx.Flash.Success(ctx.Tr("settings.theme_update_success"))
	}

	ctx.Redirect(setting.AppSubURL + "/user/settings/appearance")
}

// UpdateUserLang update a user's language
func UpdateUserLang(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.UpdateLanguageForm)
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsAppearance"] = true

	if form.Language != "" {
		if !util.SliceContainsString(setting.Langs, form.Language) {
			ctx.Flash.Error(ctx.Tr("settings.update_language_not_found", form.Language))
			ctx.Redirect(setting.AppSubURL + "/user/settings/appearance")
			return
		}
	}

	opts := &user_service.UpdateOptions{
		Language: optional.Some(form.Language),
	}
	if err := user_service.UpdateUser(ctx, ctx.Doer, opts); err != nil {
		ctx.ServerError("UpdateUser", err)
		return
	}

	// Update the language to the one we just set
	middleware.SetLocaleCookie(ctx.Resp, ctx.Doer.Language, 0)

	log.Trace("User settings updated: %s", ctx.Doer.Name)
	ctx.Flash.Success(translation.NewLocale(ctx.Doer.Language).TrString("settings.update_language_success"))
	ctx.Redirect(setting.AppSubURL + "/user/settings/appearance")
}
