// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package notify

import (
	"context"

	packages_model "code.gitea.io/gitea/models/packages"
	user_model "code.gitea.io/gitea/models/user"
)

// NullNotifier implements a blank notifier
type NullNotifier struct{}

var _ Notifier = &NullNotifier{}

// Run places a place holder function
func (*NullNotifier) Run() {
}

// PackageCreate places a place holder function
func (*NullNotifier) PackageCreate(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor) {
}

// PackageDelete places a place holder function
func (*NullNotifier) PackageDelete(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor) {
}
