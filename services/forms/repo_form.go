// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package forms

import (
	"net/http"
	"strings"

	"code.gitea.io/gitea/modules/web/middleware"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/webhook"

	"gitea.com/go-chi/binding"
)

// CreateRepoForm form for creating repository
type CreateRepoForm struct {
	UID           int64  `binding:"Required"`
	RepoName      string `binding:"Required;AlphaDashDot;MaxSize(100)"`
	Private       bool
	Description   string `binding:"MaxSize(2048)"`
	DefaultBranch string `binding:"GitRefName;MaxSize(100)"`
	AutoInit      bool
	Gitignores    string
	IssueLabels   string
	License       string
	Readme        string
	Template      bool

	RepoTemplate    int64
	GitContent      bool
	Topics          bool
	GitHooks        bool
	Webhooks        bool
	Avatar          bool
	Labels          bool
	ProtectedBranch bool

	ForkSingleBranch string
	ObjectFormatName string
}

// Validate validates the fields
func (f *CreateRepoForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// RepoSettingForm form for changing repository settings
type RepoSettingForm struct {
	RepoName               string `binding:"Required;AlphaDashDot;MaxSize(100)"`
	Description            string `binding:"MaxSize(2048)"`
	Website                string `binding:"ValidUrl;MaxSize(1024)"`
	Interval               string
	MirrorAddress          string
	MirrorUsername         string
	MirrorPassword         string
	LFS                    bool   `form:"mirror_lfs"`
	LFSEndpoint            string `form:"mirror_lfs_endpoint"`
	PushMirrorID           string
	PushMirrorAddress      string
	PushMirrorUsername     string
	PushMirrorPassword     string
	PushMirrorSyncOnCommit bool
	PushMirrorInterval     string
	Private                bool
	Template               bool
	EnablePrune            bool

	// Advanced settings
	EnableCode                            bool
	EnableWiki                            bool
	EnableExternalWiki                    bool
	DefaultWikiBranch                     string
	DefaultWikiEveryoneAccess             string
	ExternalWikiURL                       string
	EnableIssues                          bool
	EnableExternalTracker                 bool
	ExternalTrackerURL                    string
	TrackerURLFormat                      string
	TrackerIssueStyle                     string
	ExternalTrackerRegexpPattern          string
	EnableCloseIssuesViaCommitInAnyBranch bool
	EnableProjects                        bool
	ProjectsMode                          string
	EnableReleases                        bool
	EnablePackages                        bool
	EnablePulls                           bool
	EnableActions                         bool
	PullsIgnoreWhitespace                 bool
	PullsAllowMerge                       bool
	PullsAllowRebase                      bool
	PullsAllowRebaseMerge                 bool
	PullsAllowSquash                      bool
	PullsAllowFastForwardOnly             bool
	PullsAllowManualMerge                 bool
	PullsDefaultMergeStyle                string
	EnableAutodetectManualMerge           bool
	PullsAllowRebaseUpdate                bool
	DefaultDeleteBranchAfterMerge         bool
	DefaultAllowMaintainerEdit            bool
	EnableTimetracker                     bool
	AllowOnlyContributorsToTrackTime      bool
	EnableIssueDependencies               bool
	IsArchived                            bool

	// Signing Settings
	TrustModel string

	// Admin settings
	EnableHealthCheck  bool
	RequestReindexType string
}

// Validate validates the fields
func (f *RepoSettingForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

//  __      __      ___.   .__                   __
// /  \    /  \ ____\_ |__ |  |__   ____   ____ |  | __
// \   \/\/   // __ \| __ \|  |  \ /  _ \ /  _ \|  |/ /
//  \        /\  ___/| \_\ \   Y  (  <_> |  <_> )    <
//   \__/\  /  \___  >___  /___|  /\____/ \____/|__|_ \
//        \/       \/    \/     \/                   \/

// WebhookForm form for changing web hook
type WebhookForm struct {
	Events              string
	Repository          bool
	Package             bool
	Active              bool
	AuthorizationHeader string
}

// SendEverything if the hook will be triggered any event
func (f WebhookForm) SendEverything() bool {
	return f.Events == "send_everything"
}

// ChooseEvents if the hook will be triggered choose events
func (f WebhookForm) ChooseEvents() bool {
	return f.Events == "choose_events"
}

// NewWebhookForm form for creating web hook
type NewWebhookForm struct {
	PayloadURL  string `binding:"Required;ValidUrl"`
	HTTPMethod  string `binding:"Required;In(POST,GET)"`
	ContentType int    `binding:"Required"`
	Secret      string
	WebhookForm
}

// Validate validates the fields
func (f *NewWebhookForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// NewSlackHookForm form for creating slack hook
type NewSlackHookForm struct {
	PayloadURL string `binding:"Required;ValidUrl"`
	Channel    string `binding:"Required"`
	Username   string
	IconURL    string
	Color      string
	WebhookForm
}

// Validate validates the fields
func (f *NewSlackHookForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	if !webhook.IsValidSlackChannel(strings.TrimSpace(f.Channel)) {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"Channel"},
			Classification: "",
			Message:        ctx.Locale.TrString("repo.settings.add_webhook.invalid_channel_name"),
		})
	}
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// NewDiscordHookForm form for creating discord hook
type NewDiscordHookForm struct {
	PayloadURL string `binding:"Required;ValidUrl"`
	Username   string
	IconURL    string
	WebhookForm
}

// Validate validates the fields
func (f *NewDiscordHookForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

// NewMSTeamsHookForm form for creating MS Teams hook
type NewMSTeamsHookForm struct {
	PayloadURL string `binding:"Required;ValidUrl"`
	WebhookForm
}

// Validate validates the fields
func (f *NewMSTeamsHookForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}
