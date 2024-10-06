// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"context"
	"fmt"
	"net/http"

	webhook_model "code.gitea.io/gitea/models/webhook"
	api "code.gitea.io/gitea/modules/structs"
)

type (
	// MSTeamsFact for Fact Structure
	MSTeamsFact struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// MSTeamsSection is a MessageCard section
	MSTeamsSection struct {
		ActivityTitle    string        `json:"activityTitle"`
		ActivitySubtitle string        `json:"activitySubtitle"`
		ActivityImage    string        `json:"activityImage"`
		Facts            []MSTeamsFact `json:"facts"`
		Text             string        `json:"text"`
	}

	// MSTeamsAction is an action (creates buttons, links etc)
	MSTeamsAction struct {
		Type    string                `json:"@type"`
		Name    string                `json:"name"`
		Targets []MSTeamsActionTarget `json:"targets,omitempty"`
	}

	// MSTeamsActionTarget is the actual link to follow, etc
	MSTeamsActionTarget struct {
		Os  string `json:"os"`
		URI string `json:"uri"`
	}

	// MSTeamsPayload is the parent object
	MSTeamsPayload struct {
		Type            string           `json:"@type"`
		Context         string           `json:"@context"`
		ThemeColor      string           `json:"themeColor"`
		Title           string           `json:"title"`
		Summary         string           `json:"summary"`
		Sections        []MSTeamsSection `json:"sections"`
		PotentialAction []MSTeamsAction  `json:"potentialAction"`
	}
)

type msteamsConvertor struct{}

// Repository implements PayloadConvertor Repository method
// func (m msteamsConvertor) Repository(p *api.RepositoryPayload) (MSTeamsPayload, error) {
// 	var title, url string
// 	var color int
// 	switch p.Action {
// 	case api.HookRepoCreated:
// 		title = fmt.Sprintf("[%s] Repository created", p.Repository.FullName)
// 		url = p.Repository.HTMLURL
// 		color = greenColor
// 	case api.HookRepoDeleted:
// 		title = fmt.Sprintf("[%s] Repository deleted", p.Repository.FullName)
// 		color = yellowColor
// 	}

// 	return createMSTeamsPayload(
// 		p.Repository,
// 		p.Sender,
// 		title,
// 		"",
// 		url,
// 		color,
// 		nil,
// 	), nil
// }

func (m msteamsConvertor) Package(p *api.PackagePayload) (MSTeamsPayload, error) {
	title, color := getPackagePayloadInfo(p, noneLinkFormatter, false)

	return createMSTeamsPayload(
		nil, // p.Repository,
		p.Sender,
		title,
		"",
		p.Package.HTMLURL,
		color,
		&MSTeamsFact{"Package:", p.Package.Name},
	), nil
}

// func createMSTeamsPayload(r *api.Repository, s *api.User, title, text, actionTarget string, color int, fact *MSTeamsFact) MSTeamsPayload {
func createMSTeamsPayload(r *any, s *api.User, title, text, actionTarget string, color int, fact *MSTeamsFact) MSTeamsPayload {
	facts := make([]MSTeamsFact, 0, 2)
	if r != nil {
		facts = append(facts, MSTeamsFact{
			Name:  "Repository:",
			Value: "r.FullName", // TODO
		})
	}
	if fact != nil {
		facts = append(facts, *fact)
	}

	return MSTeamsPayload{
		Type:       "MessageCard",
		Context:    "https://schema.org/extensions",
		ThemeColor: fmt.Sprintf("%x", color),
		Title:      title,
		Summary:    title,
		Sections: []MSTeamsSection{
			{
				ActivityTitle:    s.FullName,
				ActivitySubtitle: s.UserName,
				ActivityImage:    s.AvatarURL,
				Text:             text,
				Facts:            facts,
			},
		},
		PotentialAction: []MSTeamsAction{
			{
				Type: "OpenUri",
				Name: "View in Gitea",
				Targets: []MSTeamsActionTarget{
					{
						Os:  "default",
						URI: actionTarget,
					},
				},
			},
		},
	}
}

func newMSTeamsRequest(_ context.Context, w *webhook_model.Webhook, t *webhook_model.HookTask) (*http.Request, []byte, error) {
	var pc payloadConvertor[MSTeamsPayload] = msteamsConvertor{}
	return newJSONRequest(pc, w, t, true)
}
