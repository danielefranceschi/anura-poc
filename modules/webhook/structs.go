// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

// HookEvents is a set of web hook events
type HookEvents struct {
	Create  bool `json:"create"`
	Delete  bool `json:"delete"`
	Package bool `json:"package"`
}

// HookEvent represents events that will delivery hook.
type HookEvent struct {
	SendEverything bool `json:"send_everything"`
	ChooseEvents   bool `json:"choose_events"`

	HookEvents `json:"events"`
}
