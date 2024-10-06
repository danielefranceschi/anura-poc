// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package notify

import (
	"context"

	packages_model "code.gitea.io/gitea/models/packages"
	user_model "code.gitea.io/gitea/models/user"
)

var notifiers []Notifier

// RegisterNotifier providers method to receive notify messages
func RegisterNotifier(notifier Notifier) {
	go notifier.Run()
	notifiers = append(notifiers, notifier)
}

// PackageCreate notifies creation of a package to notifiers
func PackageCreate(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor) {
	for _, notifier := range notifiers {
		notifier.PackageCreate(ctx, doer, pd)
	}
}

// PackageDelete notifies deletion of a package to notifiers
func PackageDelete(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor) {
	for _, notifier := range notifiers {
		notifier.PackageDelete(ctx, doer, pd)
	}
}
