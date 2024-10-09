// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth

import (
	"errors"
	"fmt"
	"net/http"

	"code.gitea.io/gitea/models/auth"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/auth/password"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/modules/web/middleware"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/forms"
	user_service "code.gitea.io/gitea/services/user"
)

// tplMustChangePassword template for updating a user's password
var tplMustChangePassword base.TplName = "user/auth/change_passwd"

func commonResetPassword(ctx *context.Context) (*user_model.User, *auth.TwoFactor) {
	code := ctx.FormString("code")

	ctx.Data["Title"] = ctx.Tr("auth.reset_password")
	ctx.Data["Code"] = code

	if nil != ctx.Doer {
		ctx.Data["user_signed_in"] = true
	}

	if len(code) == 0 {
		ctx.Flash.Error(ctx.Tr("auth.invalid_code_forgot_password", fmt.Sprintf("%s/user/forgot_password", setting.AppSubURL)), true)
		return nil, nil
	}

	// Fail early, don't frustrate the user
	u := user_model.VerifyUserActiveCode(ctx, code)
	if u == nil {
		ctx.Flash.Error(ctx.Tr("auth.invalid_code_forgot_password", fmt.Sprintf("%s/user/forgot_password", setting.AppSubURL)), true)
		return nil, nil
	}

	twofa, err := auth.GetTwoFactorByUID(ctx, u.ID)
	if err != nil {
		if !auth.IsErrTwoFactorNotEnrolled(err) {
			ctx.Error(http.StatusInternalServerError, "CommonResetPassword", err.Error())
			return nil, nil
		}
	} else {
		ctx.Data["has_two_factor"] = true
		ctx.Data["scratch_code"] = ctx.FormBool("scratch_code")
	}

	// Show the user that they are affecting the account that they intended to
	ctx.Data["user_email"] = u.Email

	if nil != ctx.Doer && u.ID != ctx.Doer.ID {
		ctx.Flash.Error(ctx.Tr("auth.reset_password_wrong_user", ctx.Doer.Email, u.Email), true)
		return nil, nil
	}

	return u, twofa
}

// MustChangePassword renders the page to change a user's password
func MustChangePassword(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("auth.must_change_password")
	ctx.Data["ChangePasscodeLink"] = setting.AppSubURL + "/user/settings/change_password"
	ctx.Data["MustChangePassword"] = true
	ctx.HTML(http.StatusOK, tplMustChangePassword)
}

// MustChangePasswordPost response for updating a user's password after their
// account was created by an admin
func MustChangePasswordPost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.MustChangePasswordForm)
	ctx.Data["Title"] = ctx.Tr("auth.must_change_password")
	ctx.Data["ChangePasscodeLink"] = setting.AppSubURL + "/user/settings/change_password"
	if ctx.HasError() {
		ctx.HTML(http.StatusOK, tplMustChangePassword)
		return
	}

	// Make sure only requests for users who are eligible to change their password via
	// this method passes through
	if !ctx.Doer.MustChangePassword {
		ctx.ServerError("MustUpdatePassword", errors.New("cannot update password. Please visit the settings page"))
		return
	}

	if form.Password != form.Retype {
		ctx.Data["Err_Password"] = true
		ctx.RenderWithErr(ctx.Tr("form.password_not_match"), tplMustChangePassword, &form)
		return
	}

	opts := &user_service.UpdateAuthOptions{
		Password:           optional.Some(form.Password),
		MustChangePassword: optional.Some(false),
	}
	if err := user_service.UpdateAuth(ctx, ctx.Doer, opts); err != nil {
		switch {
		case errors.Is(err, password.ErrMinLength):
			ctx.Data["Err_Password"] = true
			ctx.RenderWithErr(ctx.Tr("auth.password_too_short", setting.MinPasswordLength), tplMustChangePassword, &form)
		case errors.Is(err, password.ErrComplexity):
			ctx.Data["Err_Password"] = true
			ctx.RenderWithErr(password.BuildComplexityError(ctx.Locale), tplMustChangePassword, &form)
		case errors.Is(err, password.ErrIsPwned):
			ctx.Data["Err_Password"] = true
			ctx.RenderWithErr(ctx.Tr("auth.password_pwned", "https://haveibeenpwned.com/Passwords"), tplMustChangePassword, &form)
		case password.IsErrIsPwnedRequest(err):
			ctx.Data["Err_Password"] = true
			ctx.RenderWithErr(ctx.Tr("auth.password_pwned_err"), tplMustChangePassword, &form)
		default:
			ctx.ServerError("UpdateAuth", err)
		}
		return
	}

	ctx.Flash.Success(ctx.Tr("settings.change_password_success"))

	log.Trace("User updated password: %s", ctx.Doer.Name)

	if redirectTo := ctx.GetSiteCookie("redirect_to"); redirectTo != "" {
		middleware.DeleteRedirectToCookie(ctx.Resp)
		ctx.RedirectToCurrentSite(redirectTo)
		return
	}

	ctx.Redirect(setting.AppSubURL + "/")
}
