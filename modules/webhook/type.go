// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

// HookEventType is the type of a hook event
type HookEventType string

// Types of hook events
const (
	HookEventCreate  HookEventType = "create"
	HookEventDelete  HookEventType = "delete"
	HookEventPackage HookEventType = "package"
)

// Event returns the HookEventType as an event string
func (h HookEventType) Event() string {
	switch h {
	case HookEventCreate:
		return "create"
	case HookEventDelete:
		return "delete"
	case HookEventPackage:
		return "package"
	}
	return ""
}

// HookType is the type of a webhook
type HookType = string

// Types of webhooks
const (
	SLACK   HookType = "slack"
	DISCORD HookType = "discord"
	MSTEAMS HookType = "msteams"
)

// HookStatus is the status of a web hook
type HookStatus int

// Possible statuses of a web hook
const (
	HookStatusNone HookStatus = iota
	HookStatusSucceed
	HookStatusFail
)
