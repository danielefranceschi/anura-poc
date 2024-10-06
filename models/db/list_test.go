// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	"code.gitea.io/gitea/models/db"

	"xorm.io/builder"
)

type mockListOptions struct {
	db.ListOptions
}

func (opts mockListOptions) IsListAll() bool {
	return true
}

func (opts mockListOptions) ToConds() builder.Cond {
	return builder.NewCond()
}
