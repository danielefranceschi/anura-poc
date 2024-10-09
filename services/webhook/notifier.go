// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"context"

	packages_model "code.gitea.io/gitea/models/packages"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/log"
	api "code.gitea.io/gitea/modules/structs"
	webhook_module "code.gitea.io/gitea/modules/webhook"
	"code.gitea.io/gitea/services/convert"
	notify_service "code.gitea.io/gitea/services/notify"
)

func init() {
	notify_service.RegisterNotifier(NewNotifier())
}

type webhookNotifier struct {
	notify_service.NullNotifier
}

var _ notify_service.Notifier = &webhookNotifier{}

// NewNotifier create a new webhookNotifier notifier
func NewNotifier() notify_service.Notifier {
	return &webhookNotifier{}
}

// func (m *webhookNotifier) CreateRepository(ctx context.Context, doer, u *user_model.User, repo *repo_model.Repository) {
// 	// Add to hook queue for created repo after session commit.
// 	if err := PrepareWebhooks(ctx, EventSource{Repository: repo}, webhook_module.HookEventRepository, &api.RepositoryPayload{
// 		Action:       api.HookRepoCreated,
// 		Repository:   convert.ToRepo(ctx, repo, access_model.Permission{AccessMode: perm.AccessModeOwner}),
// 		Organization: convert.ToUser(ctx, u, nil),
// 		Sender:       convert.ToUser(ctx, doer, nil),
// 	}); err != nil {
// 		log.Error("PrepareWebhooks [repo_id: %d]: %v", repo.ID, err)
// 	}
// }

func (m *webhookNotifier) PackageCreate(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor) {
	notifyPackage(ctx, doer, pd, api.HookPackageCreated)
}

func (m *webhookNotifier) PackageDelete(ctx context.Context, doer *user_model.User, pd *packages_model.PackageDescriptor) {
	notifyPackage(ctx, doer, pd, api.HookPackageDeleted)
}

func notifyPackage(ctx context.Context, sender *user_model.User, pd *packages_model.PackageDescriptor, action api.HookPackageAction) {
	source := EventSource{
		Repository: nil, // pd.Repository,
		Owner:      pd.Owner,
	}

	apiPackage, err := convert.ToPackage(ctx, pd, sender)
	if err != nil {
		log.Error("Error converting package: %v", err)
		return
	}

	if err := PrepareWebhooks(ctx, source, webhook_module.HookEventPackage, &api.PackagePayload{
		Action:  action,
		Package: apiPackage,
		Sender:  convert.ToUser(ctx, sender, nil),
	}); err != nil {
		log.Error("PrepareWebhooks: %v", err)
	}
}
