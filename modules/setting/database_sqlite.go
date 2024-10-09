// Copyright 2014 The Gogs Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	_ "modernc.org/sqlite" // Import for SQLite3 support
)

func init() {
	EnableSQLite3 = true
	SupportedDatabaseTypes = append(SupportedDatabaseTypes, "sqlite")
}
