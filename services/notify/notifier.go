// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package notify

import (
	"context"

	packages_model "code.gitea.io/gitea/models/packages"
	user_model "code.gitea.io/gitea/models/user"
)

// Notifier defines an interface to notify receiver
type Notifier interface {
	Run()

	PackageCreate(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor)
	PackageDelete(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor)
}
