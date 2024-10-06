// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"fmt"
	"html"
	"net/url"

	webhook_model "code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	webhook_module "code.gitea.io/gitea/modules/webhook"
)

type linkFormatter = func(string, string) string

// noneLinkFormatter does not create a link but just returns the text
func noneLinkFormatter(url, text string) string {
	return text
}

// htmlLinkFormatter creates a HTML link
func htmlLinkFormatter(url, text string) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(url), html.EscapeString(text))
}

func getPackagePayloadInfo(p *api.PackagePayload, linkFormatter linkFormatter, withSender bool) (text string, color int) {
	refLink := linkFormatter(p.Package.HTMLURL, p.Package.Name+":"+p.Package.Version)

	switch p.Action {
	case api.HookPackageCreated:
		text = fmt.Sprintf("Package created: %s", refLink)
		color = greenColor
	case api.HookPackageDeleted:
		text = fmt.Sprintf("Package deleted: %s", refLink)
		color = redColor
	}
	if withSender {
		text += fmt.Sprintf(" by %s", linkFormatter(setting.AppURL+url.PathEscape(p.Sender.UserName), p.Sender.UserName))
	}

	return text, color
}

// ToHook convert models.Webhook to api.Hook
// This function is not part of the convert package to prevent an import cycle
func ToHook(repoLink string, w *webhook_model.Webhook) (*api.Hook, error) {
	config := map[string]string{
		"url":          w.URL,
		"content_type": w.ContentType.Name(),
	}
	if w.Type == webhook_module.SLACK {
		s := GetSlackHook(w)
		config["channel"] = s.Channel
		config["username"] = s.Username
		config["icon_url"] = s.IconURL
		config["color"] = s.Color
	}

	authorizationHeader, err := w.HeaderAuthorization()
	if err != nil {
		return nil, err
	}

	return &api.Hook{
		ID:                  w.ID,
		Type:                w.Type,
		URL:                 fmt.Sprintf("%s/settings/hooks/%d", repoLink, w.ID),
		Active:              w.IsActive,
		Config:              config,
		Events:              w.EventsArray(),
		AuthorizationHeader: authorizationHeader,
		Updated:             w.UpdatedUnix.AsTime(),
		Created:             w.CreatedUnix.AsTime(),
	}, nil
}
