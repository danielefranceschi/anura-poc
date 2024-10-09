// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package structs

import (
	"errors"
	"time"

	"code.gitea.io/gitea/modules/json"
)

// ErrInvalidReceiveHook FIXME
var ErrInvalidReceiveHook = errors.New("Invalid JSON payload received over webhook")

// Hook a hook is a web hook when one repository changed
type Hook struct {
	ID                  int64             `json:"id"`
	Type                string            `json:"type"`
	BranchFilter        string            `json:"branch_filter"`
	URL                 string            `json:"-"`
	Config              map[string]string `json:"config"`
	Events              []string          `json:"events"`
	AuthorizationHeader string            `json:"authorization_header"`
	Active              bool              `json:"active"`
	// swagger:strfmt date-time
	Updated time.Time `json:"updated_at"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
}

// HookList represents a list of API hook.
type HookList []*Hook

// CreateHookOptionConfig has all config options in it
// required are "content_type" and "url" Required
type CreateHookOptionConfig map[string]string

// CreateHookOption options when create a hook
type CreateHookOption struct {
	// required: true
	// enum: dingtalk,discord,gitea,gogs,msteams,slack,telegram,feishu,wechatwork,packagist
	Type string `json:"type" binding:"Required"`
	// required: true
	Config              CreateHookOptionConfig `json:"config" binding:"Required"`
	Events              []string               `json:"events"`
	BranchFilter        string                 `json:"branch_filter" binding:"GlobPattern"`
	AuthorizationHeader string                 `json:"authorization_header"`
	// default: false
	Active bool `json:"active"`
}

// EditHookOption options when modify one hook
type EditHookOption struct {
	Config              map[string]string `json:"config"`
	Events              []string          `json:"events"`
	BranchFilter        string            `json:"branch_filter" binding:"GlobPattern"`
	AuthorizationHeader string            `json:"authorization_header"`
	Active              *bool             `json:"active"`
}

// Payloader payload is some part of one hook
type Payloader interface {
	JSONPayload() ([]byte, error)
}

// PayloadUser represents the author or committer of a commit
type PayloadUser struct {
	// Full name of the commit author
	Name string `json:"name"`
	// swagger:strfmt email
	Email    string `json:"email"`
	UserName string `json:"username"`
}

// _ Payloader = &RepositoryPayload{}
var _ Payloader = &PackagePayload{}

// //__________                           .__  __
// //\______   \ ____ ______   ____  _____|__|/  |_  ___________ ___.__.
// // |       _// __ \\____ \ /  _ \/  ___/  \   __\/  _ \_  __ <   |  |
// // |    |   \  ___/|  |_> >  <_> )___ \|  ||  | (  <_> )  | \/\___  |
// // |____|_  /\___  >   __/ \____/____  >__||__|  \____/|__|   / ____|
// //        \/     \/|__|              \/                       \/

// // HookRepoAction an action that happens to a repo
// type HookRepoAction string

// const (
// 	// HookRepoCreated created
// 	HookRepoCreated HookRepoAction = "created"
// 	// HookRepoDeleted deleted
// 	HookRepoDeleted HookRepoAction = "deleted"
// )

// // RepositoryPayload payload for repository webhooks
// type RepositoryPayload struct {
// 	Action       HookRepoAction `json:"action"`
// 	Repository   *Repository    `json:"repository"`
// 	Organization *User          `json:"organization"`
// 	Sender       *User          `json:"sender"`
// }

// // JSONPayload JSON representation of the payload
// func (p *RepositoryPayload) JSONPayload() ([]byte, error) {
// 	return json.MarshalIndent(p, "", " ")
// }

// HookPackageAction an action that happens to a package
type HookPackageAction string

const (
	// HookPackageCreated created
	HookPackageCreated HookPackageAction = "created"
	// HookPackageDeleted deleted
	HookPackageDeleted HookPackageAction = "deleted"
)

// PackagePayload represents a package payload
type PackagePayload struct {
	Action HookPackageAction `json:"action"`
	// Repository   *Repository       `json:"repository"`
	Package *Package `json:"package"`
	Sender  *User    `json:"sender"`
}

// JSONPayload implements Payload
func (p *PackagePayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
